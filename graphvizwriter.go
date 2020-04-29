package main

import (
	"github.com/awalterschulze/gographviz"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"github.com/rs/zerolog/log"
	"os"
)

func WriteGraph(outputFilePath string, root *tfgraph.Module, dependencies []*tfgraph.Dependency, graphName string) {
	// Build the output Graphviz graph.
	graph := gographviz.NewGraph()
	graph.Name = escape(graphName)
	graph.Directed = true
	graph.AddAttr(graph.Name, string(gographviz.NewRank), "true")

	// Add all the modules as clusters.
	createCluster(graph, graph.Name, root)

	// Add the edges.
	for _, dep := range dependencies {
		shape, style := getEdgeShapeAndStyle(dep)
		_ = graph.AddEdge(
			escape(dep.Src.GetQualifiedName()),
			escape(dep.Dst.GetQualifiedName()),
			true,
			map[string]string{
				string(gographviz.Shape): escape(shape),
				string(gographviz.Style): escape(style),
			},
		)
	}

	// Get the file to use.
	file := os.Stdout
	var err error
	if outputFilePath != "" {
		// Write to the specified file.
		file, err = os.Create(outputFilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot open the specified file for writing")
		}
		defer file.Close()
	}

	// Output the result.
	output := graph.String()
	_, _ = file.WriteString(output)
}

func createCluster(graph *gographviz.Graph, parentName string, module *tfgraph.Module) {
	clusterName := escape("cluster_" + module.GetQualifiedName())
	_ = graph.AddSubGraph(
		parentName,
		clusterName,
		map[string]string{
			string(gographviz.Label): escape(module.GetName()),
			string(gographviz.Style): escape("dashed"),
		},
	)

	for _, child := range module.Children {
		// Check if the current child is itself a module.
		subModule, ok := child.(*tfgraph.Module)
		if ok {
			// Yes, recursively add the sub-module
			createCluster(graph, clusterName, subModule)
		} else {
			// No, add the config element to the current cluster.
			shape, style := getNodeShapeAndStyle(child)
			_ = graph.AddNode(
				clusterName,
				escape(child.GetQualifiedName()),
				map[string]string{
					string(gographviz.Label): escape(child.GetName()),
					string(gographviz.Shape): escape(shape),
					string(gographviz.Style): escape(style),
				},
			)
		}
	}
}

func escape(s string) string {
	return "\"" + s + "\""
}

func getNodeShapeAndStyle(elt tfgraph.ConfigElement) (string, string) {
	shape := ""
	style := ""
	switch elt.GetTfType() {
	case tfgraph.TfResource:
		shape = "box"
		style = "rounded"
	case tfgraph.TfVar:
		shape = "ellipse"
	case tfgraph.TfLocal:
		shape = "ellipse"
	case tfgraph.TfOutput:
		shape = "note"
	case tfgraph.TfProvider:
		shape = "diamond"
	}

	return shape, style
}

func getEdgeShapeAndStyle(dep *tfgraph.Dependency) (string, string) {
	shape := ""
	style := ""
	switch dep.Dst.GetTfType() {
	case tfgraph.TfResource:
		style = "solid"
	case tfgraph.TfVar:
		style = "dotted"
	case tfgraph.TfLocal:
		style = "dotted"
	case tfgraph.TfOutput:
		style = "dashed"
	case tfgraph.TfProvider:
		style = "solid"
	}

	return shape, style
}
