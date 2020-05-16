package cytoscape

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"github.com/rs/zerolog/log"
	"html/template"
	"os"
)

type FormattingOptions struct {
	GraphName    string
	EmbedModules bool // TODO
}

func WriteGraph(
	outputFilePath string,
	root *tfgraph.Module,
	dependencies []*tfgraph.Dependency,
	formattingOptions *FormattingOptions,
) {
	// Get the Cytoscape graph elements.
	graphElements := getGraphElements(root, dependencies, formattingOptions)

	// Encode result to  JSON.
	var buf bytes.Buffer
	graphW := bufio.NewWriter(&buf)
	enc := json.NewEncoder(graphW)
	if err := enc.Encode(&graphElements); err != nil {
		log.Fatal().Err(err).Msg("Cannot encode Cytoscape graph buffer")
	}
	if err := graphW.Flush(); err != nil {
		log.Fatal().Err(err).Msg("Cannot flush Cytoscape graph buffer")
	}

	tmpl := template.Must(template.ParseFiles("index.gohtml"))
	err := tmpl.Execute(os.Stdout, &map[string]interface{}{
		"PageTitle":         formattingOptions.GraphName,
		"GraphElementsJSON": template.JS(buf.String()),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot render HTML template")
	}
}

func getGraphElements(
	root *tfgraph.Module,
	dependencies []*tfgraph.Dependency,
	formattingOptions *FormattingOptions,
) *Elements {
	// Get the graph nodes.
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

	// Get the graph edges.
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

	// Encode result to Cytoscape JSON.
	return &Elements{
		Nodes: nodes,
		Edges: edges,
	}
}
