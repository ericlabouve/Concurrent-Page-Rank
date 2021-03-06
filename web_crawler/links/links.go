// Copyright © 2016 The Go Programming Language 
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

//!+Extract

// Package links provides a link-extraction function.
package links

import (
	"fmt"
	"net/http"
	"strings"
	"regexp"
	"golang.org/x/net/html"
)

// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url, username, passwd string) ([]string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				// only save url if it is in the calpoly.edu domain
				link_str := link.String()
				regex := regexp.MustCompile("http")
				num_instances := len(regex.FindAllStringIndex(link_str, -1))
				
				if strings.Contains(link_str, "calpoly.edu") && num_instances == 1 {
					if strings.Contains(link_str, "#") {
						link_str = strings.Split(link_str, "#")[0]
					}
					links = append(links, link_str)
				}
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//!-Extract

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
