package cytoscape

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"github.com/rs/zerolog/log"
	"html/template"
	"io"
)

type FormattingOptions struct {
	GraphName    string
	EmbedModules bool // TODO
}

// WriteGraphJSON writes the specified Terraform graph in Cytoscape.js JSON format.
func WriteGraphJSON(writer io.Writer, graph *tfgraph.Graph) {
	// Get the graph elements.
	graphElements := getGraphElements(graph)

	// Encode the result to JSON.
	enc := json.NewEncoder(writer)
	if err := enc.Encode(&graphElements); err != nil {
		log.Fatal().Err(err).Msg("Cannot encode Cytoscape.js graph to JSON")
	}
}

// WriteGraphHTML writes the specified Terraform graph to an HTML file using the given template.
func WriteGraphHTML(writer io.Writer, graph *tfgraph.Graph, formattingOptions *FormattingOptions) {
	// Get the graph elements JSON.
	var buf bytes.Buffer
	graphWriter := bufio.NewWriter(&buf)
	WriteGraphJSON(graphWriter, graph)

	// TODO : give template as parameter
	tmpl := template.Must(template.ParseFiles("index.gohtml"))
	err := tmpl.Execute(writer, &map[string]interface{}{
		"PageTitle":         formattingOptions.GraphName,
		"GraphElementsJSON": template.JS(buf.String()),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot render HTML template")
	}
}

func getGraphElements(graph *tfgraph.Graph) *Elements {
	// Get the graph nodes.
	var nodes []*Node
	var addElement func(parent *tfgraph.Module, element tfgraph.ConfigElement)
	addElement = func(parent *tfgraph.Module, element tfgraph.ConfigElement) {
		// Add a node for the element.
		node := &Node{
			Data: NodeData{
				ID:    element.GetQualifiedName(),
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
	addElement(nil, graph.Root)

	// Get the graph edges.
	var edges []*Edge
	for _, dep := range graph.Dependencies {
		src := dep.Source.GetQualifiedName()
		dst := dep.Destination.GetQualifiedName()
		edge := &Edge{
			Data: EdgeData{
				ID:     fmt.Sprintf("%s -> %s", src, dst),
				Source: src,
				Target: dst,
			},
		}
		edges = append(edges, edge)
	}

	return &Elements{
		Nodes: nodes,
		Edges: edges,
	}
}
