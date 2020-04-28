package tfgraph

import (
	"github.com/awalterschulze/gographviz"
	"github.com/rs/zerolog/log"
	"strings"
)

// Builds the Terraform configuration element hierarchy from the specified Graphviz graph.
func BuildTfGraphFromGraphviz(graph *gographviz.Graph) *Module {
	// Create the root module.
	tfGraphRoot := NewModule(nil, "")
	for _, node := range graph.Nodes.Nodes {
		// Check the node name.
		if !tfConfigElementRegexp.MatchString(node.Name) {
			log.Fatal().
				Str("name", node.Name).
				Msg("Invalid node name")
		}

		// Get the config element qualified name and find its parent module sub-graph(es).
		qualifiedName := strings.ReplaceAll(node.Name, "\"", "")
		module := tfGraphRoot
		for {
			// Check if the current qualified name starts with a module reference.
			matches := tfModuleRegexp.FindStringSubmatch(qualifiedName)
			if matches == nil {
				// Ok, all modules were stripped from the qualified name
				// and the `module` variable is the config element parent.
				break
			}

			moduleName := matches[1]
			childModule, ok := module.Children[moduleName].(*Module)
			if !ok {
				// First time we see this module name, create a new module element.
				childModule = NewModule(module, moduleName)
				module.Children[moduleName] = childModule
			}
			module = childModule
			qualifiedName = matches[2]
		}

		// Add a new config element node to the current module.
		tfType := TfResource
		for managedTfType := range ManagedTerraformTypes {
			if strings.HasPrefix(qualifiedName, managedTfType+".") {
				tfType = managedTfType
				break
			}
		}
		elt := NewBaseConfigElement(module, qualifiedName, tfType)
		module.AddChild(elt)

		// TODO : edges
	}

	// Return the "root" module.
	return tfGraphRoot.Children["root"].(*Module)
}
