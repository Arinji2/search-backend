package scraper

import (
	"regexp"
	"slices"
	"strings"

	"github.com/Arinji2/search-backend/utils"
	"golang.org/x/net/html"
)

var blockedTerms = []string{
	"facebook", "twitter", "dashboard",
}

func extractLinks(n *html.Node) []string {
	var content []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				link := extractLink(n)
				extensionRegex := `\.(png|jpe?g|gif|bmp|svg|webp|tiff?|pdf|docx?|xlsx?|pptx?|mp4|avi|mov|mkv|mp3|wav|flac|zip|rar|tar|gz|7z)$`
				blockedPattern := regexp.MustCompile(`(?i)` + strings.Join(blockedTerms, "|"))
				isInvalidLink := regexp.MustCompile(extensionRegex).MatchString(link)
				if isInvalidLink {
					return
				}
				if blockedPattern.MatchString(link) {
					return
				}
				if slices.Contains(content, link) {
					return
				}
				if !utils.IsEnglishURL(link) {
					return
				}

				//make sure no blocked links exist
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
