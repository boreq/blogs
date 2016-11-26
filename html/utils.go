package html

import (
	"errors"
	"golang.org/x/net/html"
)

func GetAttrVal(n *html.Node, key string) (string, error) {
	if n != nil {
		for _, attr := range n.Attr {
			if attr.Key == key {
				return attr.Val, nil
			}
		}
	}
	return "", errors.New("No such attribute")
}

func HasAttr(n *html.Node, key string) bool {
	_, err := GetAttrVal(n, key)
	return err == nil
}

func HasAttrVal(n *html.Node, key string, val string) bool {
	v, err := GetAttrVal(n, key)
	return err == nil && v == val
}

func IsHtmlNode(n *html.Node, name string) bool {
	if n == nil {
		return false
	}
	return n.Type == html.ElementNode && n.Data == name
}

func IsTextNode(n *html.Node) bool {
	if n == nil {
		return false
	}
	return n.Type == html.TextNode
}

type WalkFunc func(n *html.Node)

func WalkAllNodes(root *html.Node, walkF WalkFunc) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		walkF(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
}
