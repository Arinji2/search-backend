package scraper

import (
	"regexp"
	"slices"
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
				link := extractLink(n)
				extensionRegex := `\.(png|jpe?g|gif|bmp|svg|webp|tiff?|pdf|docx?|xlsx?|pptx?|mp4|avi|mov|mkv|mp3|wav|flac|zip|rar|tar|gz|7z)$`
				isInvalidLink := regexp.MustCompile(extensionRegex).MatchString(link)
				if isInvalidLink {
					return
				}
				if slices.Contains(content, link) {
					return
				}
				content = append(content, link)
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
