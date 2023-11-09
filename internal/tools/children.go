package tools

import "golang.org/x/net/html"

func GetHTMLChildren(node *html.Node, filter func(*html.Node) bool) (children []*html.Node) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ok := true
		if filter != nil {
			ok = filter(c)
		}
		if ok {
			children = append(children, c)
		}
	}
	return
}
