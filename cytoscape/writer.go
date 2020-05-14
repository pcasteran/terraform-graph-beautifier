package cytoscape

import (
	"encoding/json"
	"fmt"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"os"
)

type FormattingOptions struct {
	GraphName    string
	EmbedModules bool
}

func WriteGraph(
	outputFilePath string,
	root *tfgraph.Module,
	dependencies []*tfgraph.Dependency,
	formattingOptions *FormattingOptions,
) {
	// TODO : use params

	// Add the nodes.
	var nodes []*Node
	var addElement func(parent *tfgraph.Module, element tfgraph.ConfigElement)
	addElement = func(parent *tfgraph.Module, element tfgraph.ConfigElement) {
		// Add a node for the element.
		node := &Node{
			Data: NodeData{
				Id:    element.GetQualifiedName(),
				Label: element.GetName(),
			},
			Classes: []string{element.GetTfType()},
		}
		if parent != nil {
			parentName := parent.GetQualifiedName()
			node.Data.Parent = &parentName
		}
		nodes = append(nodes, node)

		// If the element is a module, recursively add its children.
		module, ok := element.(*tfgraph.Module)
		if ok {
			for _, child := range module.Children {
				addElement(module, child)
			}
		}
	}
	addElement(nil, root)

	// Add the edges.
	var edges []*Edge
	for _, dep := range dependencies {
		src := dep.Src.GetQualifiedName()
		dst := dep.Dst.GetQualifiedName()
		edge := &Edge{
			Data: EdgeData{
				Id:     fmt.Sprintf("%s -> %s", src, dst),
				Source: src,
				Target: dst,
			},
		}
		edges = append(edges, edge)
	}

	// TODO : temp for test
	elts := Elements{
		Nodes: nodes,
		Edges: edges,
	}

	enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", "  ")
	if err := enc.Encode(&elts); err != nil {
		// TODO
	}
}
