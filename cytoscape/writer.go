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
	"io/ioutil"
	"net/http"
)

// RenderingOptions contains all the options used during graph rendering.
type RenderingOptions struct {
	GraphName    string
	EmbedModules bool
	HTMLTemplate http.File
}

// WriteGraphJSON writes the specified Terraform graph in Cytoscape.js JSON format.
func WriteGraphJSON(writer io.Writer, graph *tfgraph.Graph, options *RenderingOptions) {
	// Get the graph elements.
	graphElements := getGraphElements(graph, options)

	// Encode the result to JSON.
	enc := json.NewEncoder(writer)
	if err := enc.Encode(&graphElements); err != nil {
		log.Fatal().Err(err).Msg("Cannot encode Cytoscape.js graph to JSON")
	}
}

// WriteGraphHTML writes the specified Terraform graph to an HTML file using the given template.
func WriteGraphHTML(writer io.Writer, graph *tfgraph.Graph, options *RenderingOptions) {
	// Get the graph elements JSON.
	var buf bytes.Buffer
	graphWriter := bufio.NewWriter(&buf)
	WriteGraphJSON(graphWriter, graph, options)
	graphWriter.Flush()

	// Load and execute the template.
	b, err := ioutil.ReadAll(options.HTMLTemplate)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot read HTML template")
	}
	tmpl := template.Must(template.New("cyto-html").Parse(string(b)))
	err = tmpl.Execute(writer, &map[string]interface{}{
		"PageTitle":         options.GraphName,
		"GraphElementsJSON": template.JS(buf.String()),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot render HTML template")
	}
}

func getGraphElements(graph *tfgraph.Graph, options *RenderingOptions) *Elements {
	// First, copy the graph dependencies as we may need to had some for the module -> module relations.
	deps := make([]*tfgraph.Dependency, len(graph.Dependencies))
	copy(deps, graph.Dependencies)

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
			if !options.EmbedModules && parent != nil {
				// Remove the parent and create a dependency for the module -> module relation.
				node.Data.Parent = nil
				deps = append(deps, &tfgraph.Dependency{
					Source:      parent,
					Destination: element,
				})
			}

			for _, child := range module.Children {
				addElement(module, child)
			}
		}
	}
	addElement(nil, graph.Root)

	// Get the graph edges.
	var edges []*Edge
	for _, dep := range deps {
		src := dep.Source.GetQualifiedName()
		dst := dep.Destination.GetQualifiedName()
		edge := &Edge{
			Data: EdgeData{
				ID:     fmt.Sprintf("%s-%s", src, dst),
				Source: src,
				Target: dst,
			},
			Classes: []string{fmt.Sprintf("%s-%s", dep.Source.GetTfType(), dep.Destination.GetTfType())},
		}
		edges = append(edges, edge)
	}

	return &Elements{
		Nodes: nodes,
		Edges: edges,
	}
}
