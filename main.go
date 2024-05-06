package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	sitter "github.com/garid3000/go-tree-sitter" //sitter "github.com/smacker/go-tree-sitter"
	"github.com/garid3000/go-tree-sitter/org"

	"hash/maphash"
)

var fileString []byte             // this contains the source code
var outputFile *os.File           // this for the ooutput file
var treeNodeId = make([]int, 100) // this contains the id
var hmap maphash.Hash

// var map_of_nodes map[string]string // 1-1-2 -> uint64

func slice_of_int_to_string(alist []int, depht int, sep string, endWithSep bool) string {
	result := ""
	for i, element := range alist {
		if i >= depht {
			if i != 0 && endWithSep == false {
				result = result[:len(result)-1]
				// just to make sure it won't end with -
			}
			break
		}
		result += fmt.Sprintf("%d%s", element, sep)
	}
	return result
}

func slice_of_int_to_hash_string(alist []int, depht int) string {
	hmap.SetSeed(hmap.Seed())
	hmap.WriteString(slice_of_int_to_string(alist, depht, "-", false)) // should have error handler?
	return "org" + strconv.FormatUint(hmap.Sum64(), 36)
}

func VisitNode(node *sitter.Node, depth int, outputfile *os.File) {
	fmt.Println(
		slice_of_int_to_string(treeNodeId, depth, "-", false),
		"\t", depth, "-depth", node.Type(), "\t|\t", "\t^\t", node.Content(fileString)) // Symbol ~~ type

	switch node.Type() {
	case "headline":
		outputfile.WriteString(
			fmt.Sprintf(
				"<h%d id=\"%s\">",
				depth-1+1, slice_of_int_to_hash_string(treeNodeId, depth-1+1),
				// depth should be depth-of-section which is parent of headline
				// but the html header start at 2 thus +1
			),
		)
	case "section":
		outputfile.WriteString(
			fmt.Sprintf(
				"<div id=\"outline-container-%s\" class=\"outline-%d\">\n",
				slice_of_int_to_hash_string(treeNodeId, depth), depth+1, // not sure why +1
			),
		)

	case "stars":
		outputfile.WriteString(
			fmt.Sprintf(
				"<span class=\"section-number-%d\">%s ",
				depth-2+1, // not sure why +1
				slice_of_int_to_string(treeNodeId, depth-2, ".", true),
				// depth is depth-of-section which is grandparent of stars
				// section -> headline -> stars
			),
		)

	case "item":
		outputfile.WriteString(
			fmt.Sprintf("%s", node.Content(fileString)),
		)

	case "body":
		outputfile.WriteString(
			fmt.Sprintf(
				"<div class=\"outline-text-%d\" id=\"text-%s\">\n",
				depth-1+1, slice_of_int_to_string(treeNodeId, depth-1, "-", false),
				//depth is the section-of-depth which is parent of body
				// not sure why +1
			),
		)

	case "paragraph":
		outputfile.WriteString(
			fmt.Sprintf(
				"<p>\n%s",
				node.Content(fileString),
			),
		)

	default:
	}

	//visit the child-nodes
	for i := 0; i < int(node.ChildCount()); i++ {
		child_node := node.Child(i)
		treeNodeId[depth] = i // should be the section
		treeNodeId[depth+1] = 0
		if node.Type() == "section" {
			treeNodeId[depth] = i - 1
			// headline, body, section
			// next_depth+ // depth changes only when readling child of section?
		}
		VisitNode(child_node, depth+1, outputfile)
	}

	switch node.Type() {
	case "headline":
		outputfile.WriteString(fmt.Sprintf("</h%d>\n", depth-1+1))
		// also depth of section;  section -> headline  thus -1
		// but the html header starts with 2 thus +1
	case "section":
		outputfile.WriteString("</div>\n")
	case "item":
		outputfile.WriteString("")
	case "stars":
		outputfile.WriteString("</span>")
	case "body":
		outputfile.WriteString("</div>\n")
	case "paragraph":
		outputfile.WriteString("</p>\n")
	default:
	}

}

func main() {
	// reading data
	fname := os.Args[1] // cmd1: filepath
	readData, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	fileString = readData

	lang := org.GetLanguage()

	//open file to output
	outputfile, err := os.OpenFile("/tmp/asdf.html", os.O_WRONLY|os.O_CREATE, 0600) //|os.O_APPEND
	if err != nil {
		log.Fatal(err)
	}
	defer outputfile.Close()

	//surfing the nodes
	rootNode, err := sitter.ParseCtx(context.Background(), fileString, lang)
	if err != nil {
		log.Fatal(err)
	}
	VisitNode(rootNode, 0, outputfile)
}
