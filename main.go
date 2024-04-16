package main

import (
	"context"
	"log"
	"os"

	sitter "github.com/garid3000/go-tree-sitter" //sitter "github.com/smacker/go-tree-sitter"
	"github.com/garid3000/go-tree-sitter/org"
)

var fileString = []byte(``)
var outputFile *os.File

func VisitNode(node *sitter.Node, depth int) {

}

func main(){
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
	VisitNode(rootNode, 0)


}
