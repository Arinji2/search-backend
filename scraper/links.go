package scraper

import (
	"strings"

	"golang.org/x/net/html"
)

func extractLinks(n *html.Node) []string {
	var content []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				links := extractLink(n)
				content = append(content, links)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return content
}

func extractLink(n *html.Node) string {
	var href string

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	if href == "" {
		return ""
	}

	return strings.TrimSpace(href)

}
