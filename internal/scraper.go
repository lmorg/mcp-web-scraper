package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

var (
	// I know you shouldn't use regex to parse HTML.
	// This is only used in the extreme edge case that a markdown document
	// cannot be automatically generated from the HTML document. At that
	// point the HTML parser has already failed and we are now looking to
	// use an LLM for parsing. In that instance, our token count will be
	// massive so stripping the following HTML tags via regexp, while crude,
	// will reduce the token count.
	rxHtml = []*regexp.Regexp{
		regexp.MustCompile(`(?si)<head( |>).*?</head>`),
		regexp.MustCompile(`(?si)<svg( |>).*?</svg>`),
		regexp.MustCompile(`(?si)<script( |>).*?</script>`),
		regexp.MustCompile(`(?si)<!--.*?-->`),
	}
)

func Scrape(ctx context.Context, url string) (string, error) {
	response, chromeScraperErr := ChromeScraper(ctx, url)

	if chromeScraperErr != nil || response == "" {
		var goScraperErr error
		response, goScraperErr = GoScraper(url)

		if goScraperErr != nil {
			return "", fmt.Errorf("cannot scrape URL: Chrome[%v], Go[%v]", chromeScraperErr, goScraperErr)
		}

		if response == "" {
			return "", fmt.Errorf("empty page")
		}
	}

	/*md, mdErr := htmltomarkdown.ConvertString(response)
	if mdErr == nil {
		// markdown conversion successful so lets return early
		return md, nil
	}*/

	// we cannot parse the HTML document via correct methods,
	// so now lets focus on reducing the token count so the LLM
	// can parse the HTML document fast and cost-effectively.
	// Yes it's ugly, but it works.
	for _, rx := range rxHtml {
		found := rx.FindAllString(response, -1)
		for i := range found {
			log.Println(found[i])
			response = strings.Replace(response, found[i], "", 1)
		}
	}

	return response, nil
}

// ChromeScraper requires Google Chrome installed.
// The benefit of this is that single page applications and other documents
// which require Javascript to fetch their content, can still be scraped as if
// they were a traditional web page.
func ChromeScraper(ctx context.Context, url string) (string, error) {
	chromeCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var body, article string

	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second), // this is a kludge to allow dynamic sites which require JS to finish rendering
		chromedp.InnerHTML("body", &body, chromedp.ByQuery),
		chromedp.InnerHTML("article", &article, chromedp.ByQuery),
	)

	if err != nil {
		return "", err
	}

	// ideally we'd expect pages we want scraped to be articles but lets
	// fallback to returning the entire body if the page either isn't an
	// article and/or doesn't include the HTML5 <article> tag.
	if article != "" {
		return article, nil
	}

	return body, nil
}

// GoScraper returns a document using Go's HTTP user agent.
// No fancy rendering tricks happen here. So if the document is, for example,
// a single page application that requires Javascript to return any specific
// content, then you'll not get much useful with this function. This this is
// used as a fallback rather than a primary tool for web scraping.
func GoScraper(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating fallback request: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making fallback request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading fallback request: %v", err)
	}

	return string(body), err
}
