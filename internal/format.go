package internal

import (
	"log"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func toMarkdown(input string) string {
	md, err := htmltomarkdown.ConvertString(input)
	if err != nil {
		return ""
	}
	return md
}

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

func toHtml(input string) string {
	// we cannot parse the HTML document via correct methods,
	// so now lets focus on reducing the token count so the LLM
	// can parse the HTML document fast and cost-effectively.
	// Yes it's ugly, but it works.
	for _, rx := range rxHtml {
		found := rx.FindAllString(input, -1)
		for i := range found {
			log.Println(found[i])
			input = strings.Replace(input, found[i], "", 1)
		}
	}

	return input
}
