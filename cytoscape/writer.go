package cytoscape

import (
	"encoding/json"
	"github.com/pcasteran/terraform-graph-beautifier/tfgraph"
	"log"
	"os"
)

type FormattingOptions struct {
	EmbedModules bool
}

func WriteGraph(
	outputFilePath string,
	root *tfgraph.Module,
	dependencies []*tfgraph.Dependency,
	graphName string,
	formattingOptions *FormattingOptions,
) {
	// TODO : use params

	// TODO : temp for test
	elts := Elements{}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&elts); err != nil {
		log.Println(err)
	}
}
