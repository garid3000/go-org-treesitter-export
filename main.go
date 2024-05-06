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

func slice_of_int_to_string(alist []int, depht int) string {
	result := ""
	for i, element := range alist {
		if i >= depht {
			if i != 0 {
				result = result[:len(result)-1]
				// just to make sure it won't end with -
			}
			break
		}
		result += fmt.Sprintf("%d-", element)
	}
	return result
}

func slice_of_int_to_hash_string(alist []int, depht int) string {
	hmap.SetSeed(hmap.Seed())
	hmap.WriteString(slice_of_int_to_string(alist, depht)) // should have error handler?
	return strconv.FormatUint(hmap.Sum64(), 10)
}

func VisitNode(node *sitter.Node, depth int, outputfile *os.File) {
	//treeNodeId = append(treeNodeId, 0)
	fmt.Println(slice_of_int_to_string(treeNodeId, depth), "\t", depth, "-depth", node.Type(), "\t|\t", "\t^\t", node.Content(fileString)) // Symbol ~~ type

	switch node.Type() {
	case "headline":
		outputfile.WriteString(
			fmt.Sprintf(
				"<h%d id=\"%s\">",
				depth-1, slice_of_int_to_hash_string(treeNodeId, depth-1),
				// depth should be depth-of-section which is parent of headline
			),
		)
	case "section":
		outputfile.WriteString(
			fmt.Sprintf(
				"<div id=\"outline-container-%s\" class=\"outline-%d\">",
				slice_of_int_to_hash_string(treeNodeId, depth), depth,
			),
		)

	case "stars":
		outputfile.WriteString(
			fmt.Sprintf(
				"<span class=\"section-number-%d\">%s ",
				depth-2,
				slice_of_int_to_string(treeNodeId, depth-2),
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
				"<div class=\"outline-text-%d\" id=\"text-%s\">",
				depth, slice_of_int_to_string(treeNodeId, depth),
			),
		)

	case "paragraph":
		outputfile.WriteString(
			fmt.Sprintf(
				"<p>%s",
				node.Content(fileString),
			),
		)

	default:
	}

	//visit the child-nodes
	for i := 0; i < int(node.ChildCount()); i++ {
		treeNodeId[depth] = i
		treeNodeId[depth+1] = 0
		child_node := node.Child(i)
		VisitNode(child_node, depth+1, outputfile)
	}

	switch node.Type() {
	case "headline":
		outputfile.WriteString(fmt.Sprintf("</h%d>", depth-1))
		// also depth of section;  section -> headline
	case "section":
		outputfile.WriteString("</div>")
	case "item":
		outputfile.WriteString("")
	case "stars":
		outputfile.WriteString("</span>")
	case "body":
		outputfile.WriteString("</div>")
	case "paragraph":
		outputfile.WriteString("</p>")
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
	outputfile, err := os.OpenFile("/tmp/asdf.html", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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
