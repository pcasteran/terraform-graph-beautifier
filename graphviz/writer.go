package graphviz

import (
	"github.com/awalterschulze/gographviz"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"io"
)

const clusterRefNodeName = "clusterRef"

type FormattingOptions struct {
	GraphName    string
	EmbedModules bool
}

func WriteGraph(writer io.Writer, graph *tfgraph.Graph, formattingOptions *FormattingOptions) {
	// Build the output Graphviz outputGraph.
	outputGraph := gographviz.NewGraph()
	outputGraph.Name = escape(formattingOptions.GraphName)
	outputGraph.Directed = true
	_ = outputGraph.AddAttr(outputGraph.Name, string(gographviz.NewRank), "true")
	_ = outputGraph.AddAttr(outputGraph.Name, string(gographviz.Compound), "true")
	_ = outputGraph.AddAttr(outputGraph.Name, string(gographviz.RankDir), "TB")

	// Add all the modules as clusters.
	createCluster(outputGraph, "", graph.Root, formattingOptions)

	// Add the edges.
	for _, dep := range graph.Dependencies {
		shape, style := getEdgeShapeAndStyle(dep)
		_ = outputGraph.AddEdge(
			escape(dep.Src.GetQualifiedName()),
			escape(dep.Dst.GetQualifiedName()),
			true,
			map[string]string{
				string(gographviz.Shape): escape(shape),
				string(gographviz.Style): escape(style),
			},
		)
	}

	// Output the result.
	_, _ = io.WriteString(writer, outputGraph.String())
}

func createCluster(graph *gographviz.Graph, parentClusterName string, module *tfgraph.Module, formattingOptions *FormattingOptions) {
	// Create the module's cluster.
	parent := module.GetParent()
	parentName := graph.Name
	if formattingOptions.EmbedModules && parent != nil {
		parentName = escape(parentClusterName)
	}
	clusterName := "cluster_" + module.GetQualifiedName()
	_ = graph.AddSubGraph(
		parentName,
		escape(clusterName),
		map[string]string{
			string(gographviz.Label): escape(module.GetName()),
		},
	)
	if !formattingOptions.EmbedModules {
		// Add an invisible node to the cluster used for the module's dependency edges.
		clusterRef := clusterName + "." + clusterRefNodeName
		_ = graph.AddNode(
			escape(clusterName),
			escape(clusterRef),
			map[string]string{
				string(gographviz.Label):       escape(""),
				string(gographviz.Style):       escape("invis"),
				string(gographviz.Width):       "0",
				string(gographviz.Height):      "0",
				string(gographviz.Peripheries): "0",
			},
		)

		// Add an edge from the parent cluster.
		if parent != nil {
			dep := tfgraph.NewDependency(parent, module)
			shape, style := getEdgeShapeAndStyle(dep)
			_ = graph.AddEdge(
				escape(parentClusterName+"."+clusterRefNodeName),
				escape(clusterRef),
				true,
				map[string]string{
					string(gographviz.Shape):      escape(shape),
					string(gographviz.Style):      escape(style),
					string(gographviz.LTail):      escape(parentClusterName),
					string(gographviz.LHead):      escape(clusterName),
					string(gographviz.Constraint): "false",
					string(gographviz.Dir):        "both",
					string(gographviz.ArrowTail):  "diamond",
					string(gographviz.ArrowHead):  "vee",
				},
			)
		}
	}

	// Add the module's children to the graph.
	for _, child := range module.Children {
		// Check if the current child is itself a module.
		subModule, ok := child.(*tfgraph.Module)
		if ok {
			// Yes, recursively add the sub-module
			createCluster(graph, clusterName, subModule, formattingOptions)
		} else {
			// No, add the config element to the current cluster.
			shape, style := getNodeShapeAndStyle(child)
			_ = graph.AddNode(
				escape(clusterName),
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
	case tfgraph.TfModule:
		style = "solid"
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
