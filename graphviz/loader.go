package graphviz

import (
	"github.com/awalterschulze/gographviz"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

// LoadGraph reads the Graphviz graph from the specified reader and returns the corresponding Terraform graph.
func LoadGraph(reader io.Reader, keepTfJunk bool, excludePatterns []string) *tfgraph.Graph {
	// Load the graph from the specified input.
	graphIn := readGraph(reader, keepTfJunk, excludePatterns)

	// Build the Terraform resource graph.
	return buildTfGraph(graphIn)
}

// Builds the Terraform configuration element hierarchy from the specified Graphviz graph.
func buildTfGraph(graph *gographviz.Graph) *tfgraph.Graph {
	// Create the graph root and build the graph from here.
	tfGraphRoot := tfgraph.NewModule(nil, "")
	nodeNameToConfigElement := make(map[string]tfgraph.ConfigElement, len(graph.Nodes.Nodes))
	for _, node := range graph.Nodes.Nodes {
		// Check the node name.
		if !tfgraph.TfConfigElementRegexp.MatchString(node.Name) {
			log.Fatal().
				Str("name", node.Name).
				Msg("Invalid node name")
		}

		// Get the config element qualified name and find its parent module sub-graph(es).
		qualifiedName := strings.ReplaceAll(node.Name, "\"", "")
		module := tfGraphRoot
		for {
			// Check if the current qualified name starts with a module reference.
			matches := tfgraph.TfModuleRegexp.FindStringSubmatch(qualifiedName)
			if matches == nil {
				// Ok, all modules were stripped from the qualified name
				// and the `module` variable is the config element parent.
				break
			}

			moduleName := matches[1]
			childModule, ok := module.Children[moduleName].(*tfgraph.Module)
			if !ok {
				// First time we see this module name, create a new module element.
				childModule = tfgraph.NewModule(module, moduleName)
				module.Children[moduleName] = childModule
			}
			module = childModule
			qualifiedName = matches[2]
		}

		// Add a new config element node to the current module.
		tfType := tfgraph.TfResource
		for managedTfType := range tfgraph.ManagedTerraformTypes {
			if strings.HasPrefix(qualifiedName, managedTfType+".") {
				tfType = managedTfType
				break
			}
		}
		elt := tfgraph.NewBaseConfigElement(module, qualifiedName, tfType)
		module.AddChild(elt)
		nodeNameToConfigElement[node.Name] = elt
	}

	// Build the edges of the graph.
	var edges []*tfgraph.Dependency
	for _, edge := range graph.Edges.Edges {
		src, ok := nodeNameToConfigElement[edge.Src]
		if !ok {
			log.Fatal().
				Str("source", edge.Src).
				Msg("Edge source is referencing an invalid node")
		}

		dst, ok := nodeNameToConfigElement[edge.Dst]
		if !ok {
			log.Fatal().
				Str("destination", edge.Dst).
				Msg("Edge destination is referencing an invalid node")
		}

		edges = append(edges, tfgraph.NewDependency(src, dst))
	}

	// Return the "root" module and the edges.
	root := tfGraphRoot.Children["module.root"].(*tfgraph.Module)
	root.SetParent(nil)
	return tfgraph.NewGraph(root, edges)
}
