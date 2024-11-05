package scraper

import (
	"golang.org/x/net/html"
)

func extractMetaInfo(n *html.Node) (string, string, string, string) {
	var title, ogImage, description, favicon string
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					title = n.FirstChild.Data
				}
			case "link":

				var href string
				isIcon := false
				for _, attr := range n.Attr {
					if attr.Key == "rel" && attr.Val == "icon" {
						isIcon = true
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if isIcon && href != "" {
					favicon = href
				}
			case "meta":
				var name, property, content string
				for _, attr := range n.Attr {
					switch attr.Key {
					case "name":
						name = attr.Val
					case "property":
						property = attr.Val
					case "content":
						content = attr.Val
					}
				}

				if name == "description" {
					description = content
				}
				if property == "og:image" {
					ogImage = content
				}
				if property == "og:title" && content != "" {
					title = content
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	return title, description, ogImage, favicon
}
