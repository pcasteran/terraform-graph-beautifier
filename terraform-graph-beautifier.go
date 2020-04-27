package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/awalterschulze/gographviz"
	tfgraph "github.com/pcasteran/terraform-graph-beautifier/tf-graph"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
)

var /* const */ tfResourceRegex = regexp.MustCompile(`^"\[root] (.*)"$`)
var /* const */ tfModuleRegex = regexp.MustCompile(`module\.(.*?)\.(.*)`)

func main() {
	// Parse command line arguments.
	inputFilePath := flag.String("input", "", "Path of the input Graphviz file to read, if not set reads stdin")
	debug := flag.Bool("debug", false, "Prints debugging information to stderr")
	flag.Parse()

	// Configure logging.
	// Default level for this example is info, unless debug flag is present.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Read the input line by line.
	var file *os.File
	var err error
	if *inputFilePath == "" {
		// Read from stdin.
		file = os.Stdin
		err = nil
	} else {
		// Read from input file.
		file, err = os.Open(*inputFilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot open specified file")
		}
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()

		// Fix some issues with the Terraform output.
		// Replace entries : map["foo"] => map['foo']
		// TODO : const
		re := regexp.MustCompile(`\["(.*?)"]`)
		line2 := re.ReplaceAllString(line, "['${1}']")
		if line != line2 {
			log.
				Debug().
				Str("before", line).
				Str("after", line2).
				Msg("Line fixed")
			line = line2
		}

		sb.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("Error while reading input")
	}

	// Load the Graphviz graph from the input.
	graphIn, err := gographviz.Read([]byte(sb.String()))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse the input as a Graphviz graph")
	}

	// Reconstruct the Terraform resources hierarchy.
	root := tfgraph.NewModule(nil, "root")
	for _, node := range graphIn.Nodes.Nodes {
		// Check node name.
		if !tfResourceRegex.MatchString(node.Name) {
			log.Fatal().
				Err(err).
				Str("name", node.Name).
				Msg("Invalid node name")
		}

		// Get the TF resource qualified name from the 'label' attribute
		// and find its parent module sub-graph(es).
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
			childModule, ok := module.Children[moduleName].(*tfgraph.Module)
			if !ok {
				childModule = tfgraph.NewModule(module, moduleName)
				module.Children[moduleName] = childModule
			}
			module = childModule
			resourceQualifiedName = matches[2]
		}

		// Add a resource node to the current module.
		// TODO
		module.AddChild(&tfgraph.AbstractConfigurationComponent{
			Parent: nil,
			Name:   resourceQualifiedName,
		})
	}

	// Build the Graphviz graph.
	graphOut := gographviz.NewGraph()
	graphOut.Name = "" // TODO : name or current directory
	graphOut.Directed = true

	// TODO : temp for tests
	output := graphIn.String()
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
