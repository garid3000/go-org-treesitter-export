package main

import (
	"context"
	"fmt"
	"log"
	"os"

	sitter "github.com/garid3000/go-tree-sitter" //sitter "github.com/smacker/go-tree-sitter"
	"github.com/garid3000/go-tree-sitter/org"
)

var fileString = []byte(``)
var outputFile *os.File

var treeNodeId = []int{}

//
//  <h2 id="orgfbd0474">
//  <span class="section-number-2">1.</span> lvl1 heading Lorem ipsum
//  </h2>
//
//  <div class="outline-text-2" id="text-1">
//  | <p>
//  |   Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam consequat
//  |   rhoncus ex at varius. Etiam lacinia ex nisl, ac vestibulum lacus
//  |   interdum pellentesque. Pellentesque rhoncus quis sem vitae convallis.
//  | </p>
//  </div>

func list_of_int_2_string_as_id(alist []int) string {
	result := "text"
	for _, element := range alist {
		result += fmt.Sprintf("-%d", element)
	}
	return result
}

func VisitNode(node *sitter.Node, depth int, outputfile *os.File) {
	treeNodeId = append(treeNodeId, 0)
	fmt.Println(treeNodeId, "\t", node.Type(), "\t|\t", "\t^\t", node.Content(fileString)) // Symbol ~~ type

	h_int := len(treeNodeId) - 1
	switch node.Type() {
	case "headline":
		fmt.Println("FUCKING HEALDINE")
		outputfile.WriteString(
			fmt.Sprintf(
				"<h%d id=\"%s\"> <span class=\"section-number-2\"> %s </span> %s </h%d>\n",
				h_int, "orgfbd0474", list_of_int_2_string_as_id(treeNodeId),  node.Content(fileString), h_int,
			),
		)
	case "body":
		outputfile.WriteString(
			fmt.Sprintf(
				"<div class=\"outline-text-%d\" id=\"%s\"> <p> %s </p> </div>\n",
				h_int, list_of_int_2_string_as_id(treeNodeId),  node.Content(fileString),
			),
		)

	default:
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		treeNodeId[len(treeNodeId)-1] = i
		child_node := node.Child(i)
		VisitNode(child_node, depth+1, outputfile)
	}
	treeNodeId = treeNodeId[0 : len(treeNodeId)-1]
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
	outputfile, err := os.OpenFile("/tmp/asdf.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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
