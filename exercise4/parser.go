package parser

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	filename := flag.String("filename", "test_data/test.html", "html file to extract all links from it")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	links, err := ParseLinks(file)
	fmt.Printf("%+v\n", links)
}

// Link represents an <a> tag
type Link struct {
	Href string
	Text string
}

// ParseLinks extracts all the links from an html page
func ParseLinks(r io.Reader) ([]Link, error) {
	links := make([]Link, 0)
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	dfs(doc, &links)

	return links, nil
}

func dfs(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.Data == "a" {
		var link Link
		link.Text = getText(n)
		for _, a := range n.Attr {
			if a.Key == "href" {
				link.Href = a.Val
				break
			}
		}
		*links = append(*links, link)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dfs(c, links)
	}
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	ans := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ans += getText(c)
	}
	return strings.Join(strings.Fields(ans), " ")
}
