package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"os"
)

func WriteGraph(root *tfgraph.Module, dependencies []*tfgraph.Dependency) {
	// Build the output Graphviz graph.
	graph := gographviz.NewGraph()
	graph.Name = "" // TODO : name or current directory
	graph.Directed = true

	// Add all the modules as clusters.
	createCluster(graph, "", root)

	// Add the edges.
	for _, dep := range dependencies {
		graph.AddEdge(
			escape(dep.Src.GetQualifiedName()),
			escape(dep.Dst.GetQualifiedName()),
			true,
			map[string]string{
				// TODO
			},
		)
	}

	// Output the result.
	output := graph.String()
	fmt.Println(output)

	// TODO : temp for tests
	fo, err := os.Create("samples/output.gv")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	_, _ = fo.WriteString(output)
}

func createCluster(graph *gographviz.Graph, parentName string, module *tfgraph.Module) {
	clusterName := escape("cluster_" + module.GetQualifiedName())
	graph.AddSubGraph(
		parentName,
		clusterName,
		map[string]string{
			string(gographviz.Label): escape(module.GetName()),
			// TODO
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
			graph.AddNode(
				clusterName,
				escape(child.GetQualifiedName()),
				map[string]string{
					string(gographviz.Label): escape(child.GetName()),
					// TODO
					// TODO custom TF type ? NewAttr()
				},
			)
		}
	}
}

func escape(s string) string {
	return "\"" + s + "\""
}
