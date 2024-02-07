package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Response body from wibly.me/json.
type response []struct {
	URL         string `json:"URL"`
	Title       string `json:"Title"`
	Snippet     string `json:"Snippet"`
	Description string `json:"Description"`
}
type elState int

const (
	Head elState = iota
	Body
	Paragraph
	Anchor
	Strong
)

func main() {
	// Currently just for testing; will also make this use my own instance in future.
	req, err := http.NewRequest("GET", "https://wiby.me/json/?q=test", nil)
	if err != nil {
		println(err)
		return
	}
	htmlData, err := os.Open("Test.html")
	if err != nil {
		println(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err)
		return
	}
	defer res.Body.Close()
	var resp response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		println(err)
		return
	}
	json.Unmarshal(body, &resp)
	fmt.Print(resp)
	doc, err := html.Parse(htmlData)
	if err != nil {
		println(err)
	}
	traverseNode(doc, 0)
}

var (
	CurrentElState elState = 3
	bBody                  = false
)

// Just traverses the node tree and tells the renderer what it's rendering.
func traverseNode(n *html.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	switch n.Type {
	case html.ElementNode:
		switch strings.TrimSpace(n.Data) {
		case "body":
			if !bBody {
				bBody = true
			}
		default:
			if bBody {
				switch strings.TrimSpace(n.Data) {
				case "strong":
					CurrentElState = Strong
				default:
					CurrentElState = Paragraph
				}
			}
		}
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" && bBody {
			renderNode(n, indent)
			CurrentElState = Paragraph
		}
	case html.CommentNode:

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNode(c, depth+1)
	}
}

// Renders node from the node tree, text being surrounded by a box.
func renderNode(n *html.Node, ind string) {
	switch CurrentElState {
	case Strong:
		fmt.Printf("%s", "\033[31m"+strings.ToUpper(n.Data)+"\033[0m")
		CurrentElState = Paragraph
	default:
		fmt.Printf("%s", n.Data)
	}
}

// Initial search screen
func startMenu() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Query server")
	url, err := reader.ReadString('\n')
	if err != nil {
		println(err)
		return
	}
	fmt.Println(url)
}
