package internal

import (
	"fmt"
	"github.com/we-are-discussing-rest/web-crawler/workers/utils"
	"golang.org/x/net/html"
	"strings"
)

func ExtractHtmlLinks(rawHtml string) ([]string, error) {
	var links []string
	parser, err := html.Parse(strings.NewReader(rawHtml))
	if err != nil {
		return nil, fmt.Errorf("error parsing html %v", err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if utils.IsValidURL(a.Val) {
						links = append(links, a.Val)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(parser)

	return links, nil
}
