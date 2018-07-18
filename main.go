// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/kr/pretty"
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/webidl/parser"
)

type testNode struct {
	nodeType   parser.NodeType
	properties map[string]interface{}
	children   map[string]*list.List
}

type parserTest struct {
	name     string
	filename string
}

func (pt *parserTest) input() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s.webidl", pt.filename))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (pt *parserTest) tree() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s.tree", pt.filename))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (pt *parserTest) writeTree(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("%s.tree", pt.filename), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

func createAstNode(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
	return &testNode{
		nodeType:   kind,
		properties: make(map[string]interface{}),
		children:   make(map[string]*list.List),
	}
}

func (tn *testNode) GetType() parser.NodeType {
	return tn.nodeType
}

func (tn *testNode) Connect(predicate string, other parser.AstNode) parser.AstNode {
	if tn.children[predicate] == nil {
		tn.children[predicate] = list.New()
	}

	tn.children[predicate].PushBack(other)
	return tn
}

func (tn *testNode) Decorate(property string, value string) parser.AstNode {
	if _, ok := tn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, tn.properties))
	}

	tn.properties[property] = value
	return tn
}

func (tn *testNode) DecorateWithInt(property string, value int) parser.AstNode {
	if _, ok := tn.properties[property]; ok {
		panic(fmt.Sprintf("Existing key for property %s\n\tNode: %v", property, tn.properties))
	}

	tn.properties[property] = value
	return tn
}

var parserTests = []parserTest{
	parserTest{"html-dom", "html-dom"},
}

func main() {
	for _, test := range parserTests {

		moduleNode := createAstNode(compilercommon.InputSource("html-dom"), parser.NodeTypeGlobalModule)

		parser.Parse(moduleNode, createAstNode, compilercommon.InputSource("html-dom"), test.input())
		parseTree := getParseTree((moduleNode).(*testNode), 0)

		//expected := strings.TrimSpace(test.tree())
		found := strings.TrimSpace(parseTree)

		//if os.Getenv("REGEN") == "true" {
		test.writeTree(found)
		//	} else {
		//}
		tn, ok := (moduleNode).(*testNode)
		if !ok {
			panic("bad assert")
		}
		pretty.Println(tn)
	}
}

func getParseTree(currentNode *testNode, indentation int) string {
	parseTree := ""
	parseTree = parseTree + strings.Repeat(" ", indentation)
	parseTree = parseTree + fmt.Sprintf("%v", currentNode.nodeType)
	parseTree = parseTree + "\n"

	keys := make([]string, 0)

	for key, _ := range currentNode.properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		parseTree = parseTree + strings.Repeat(" ", indentation+2)
		parseTree = parseTree + fmt.Sprintf("%s = %v", key, currentNode.properties[key])
		parseTree = parseTree + "\n"
	}

	keys = make([]string, 0)

	for key, _ := range currentNode.children {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := currentNode.children[key]
		parseTree = parseTree + fmt.Sprintf("%s%v =>", strings.Repeat(" ", indentation+2), key)
		parseTree = parseTree + "\n"

		for e := value.Front(); e != nil; e = e.Next() {
			parseTree = parseTree + getParseTree(e.Value.(*testNode), indentation+4)
		}
	}

	return parseTree
}
