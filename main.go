package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type resp []struct {
	URL         string `json:"URL"`
	Title       string `json:"Title"`
	Snippet     string `json:"Snippet"`
	Description string `json:"Description"`
}

func main() {
	req, _ := http.NewRequest("GET", "https://wiby.me/json/?q=test", nil)
	r, err := os.Open("Test.html")
	if err != nil {
		println(err)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	var rep resp
	bo, _ := io.ReadAll(res.Body)
	json.Unmarshal(bo, &rep)
	fmt.Print(rep)
	doc, err := html.Parse(r)
	if err != nil {
		println(err)
	}
	traverseNode(doc, 0)

}

type elState int

const (
	Head elState = iota
	Body
	Paragraph
	Anchor
	Strong
)

var CurrentElState elState = 3
var bBody = false

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
		//fmt.Printf("%s<!-- %s -->\n", indent, n.Data)
	}
	//renderNode(n.Data, indent)
	//fmt.Print(n.Data)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNode(c, depth+1)
	}
}

// TODO instead of strong and italics, use color instead; the standard way in which it is done.
func renderNode(n *html.Node, ind string) {
	if CurrentElState == Strong {
		fmt.Printf("%s", "\033[31m"+strings.ToUpper(n.Data)+"\033[0m")
		CurrentElState = Paragraph
	} else {
		fmt.Printf("%s", n.Data)
	}
}

// TODO add rendering function for specifically boxes.
// TODO OR render everything with boxes.
func mainMenu() {

}
