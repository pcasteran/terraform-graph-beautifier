package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	tf_graph "github.com/pcasteran/terraform-graph-beautifier/tf-graph"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var /* const */ tfResourceRegex = regexp.MustCompile(`^"\[root] (.*)"$`)
var /* const */ tfModuleRegex = regexp.MustCompile(`module\.(.*?)\.(.*)`)

func main() {
	// Read the input file.
	path := "samples/tf_plan.gv"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Fix some issues with the Terraform output.
	// Replace entries : map["foo"] => map['foo']
	re := regexp.MustCompile(`\["(.*?)"]`)
	b = re.ReplaceAll(b, []byte("['${1}']"))

	// Read the input graph.
	graphIn, err := gographviz.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	// Reconstruct the Terraform resources hierarchy.
	root := tf_graph.NewModule(nil, "root")
	for _, node := range graphIn.Nodes.Nodes {
		// TODO : temp for test
		fmt.Println(node.Name)

		// Check node name.
		if !tfResourceRegex.MatchString(node.Name) {
			log.Fatal("Invalid node name : ", node.Name)
		}

		// Get the TF resource qualified name from the 'label' attribute
		// and find its parent module sub-graph(es).
		// TODO : use name instead of label
		resourceQualifiedName := strings.ReplaceAll(node.Attrs[gographviz.Label], "\"", "")
		module := root
		for {
			// Check if the current qualified name (still) starts with a module reference.
			matches := tfModuleRegex.FindStringSubmatch(resourceQualifiedName)
			if matches == nil {
				// Ok, all modules were stripped from the qualified name
				// and the `module` variable is the resource parent.
				break
			}

			moduleName := matches[1]
			childModule, ok := module.Children[moduleName].(*tf_graph.Module)
			if !ok {
				childModule = tf_graph.NewModule(module, moduleName)
				module.Children[moduleName] = childModule
			}
			module = childModule
			resourceQualifiedName = matches[2]
		}

		// Add a resource node to the current module.
		// TODO
		module.AddChild(&tf_graph.AbstractConfigurationComponent{
			Parent: nil,
			Name:   resourceQualifiedName,
		})
	}

	// Build the Graphviz graph.
	graphOut := gographviz.NewGraph()
	graphOut.Name = "" // TODO : name or current directory
	graphOut.Directed = true

	output := graphOut.String()
	fmt.Println(output)

	fo, err := os.Create("output.gv")
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

// TODO : tester modules imbriqués
// Attribute "tf-type"
// TODO : shape box(ressource), arrondi(var, local), ?? (output)
//"[root] module.cloudfunction.module.storage.google_storage_bucket.buckets["artefacts"] (orphan)" [label = "module.cloudfunction.module.storage.google_storage_bucket.buckets", shape = "box"]
//   Voir (orphan)
//   Corriger label
