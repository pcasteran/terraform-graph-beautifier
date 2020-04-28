package main

import (
	"bufio"
	"github.com/awalterschulze/gographviz"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
)

type replaceStruct struct {
	replacePattern string
	report         bool
}

var /* const */ tfOutputFixes = map[*regexp.Regexp]replaceStruct{
	// Replace : `"[root] rsc_type.rsc_name"` => `"module.root.rsc_type.rsc_name"`
	regexp.MustCompile(`"\[root] `): replaceStruct{`"module.root.`, false},
	// Replace : map["foo"] => map['foo']
	regexp.MustCompile(`\["(.*?)"]`): replaceStruct{"['${1}']", true},
}

var /* const */ tfJunkMatches = []*regexp.Regexp{
	regexp.MustCompile(`"module.root.root"`),
	regexp.MustCompile(`"module.root.meta.count-boundary \(EachMode fixup\)"`),
	regexp.MustCompile(`"module.root.provider\..* \(close\)"`),
}

func loadGraph(inputFilePath string, keepTfJunk bool, excludePatterns []string) *gographviz.Graph {
	// Build all the patterns to exclude.
	var exclusionPatterns []*regexp.Regexp
	if !keepTfJunk {
		exclusionPatterns = append(exclusionPatterns, tfJunkMatches...)
	}
	for _, pattern := range excludePatterns {
		exclusionPatterns = append(exclusionPatterns, regexp.MustCompile(pattern))
	}

	// Get the file to use.
	var file *os.File
	var err error
	if inputFilePath == "" {
		// Read from stdin.
		file = os.Stdin
		err = nil
	} else {
		// Read from input file.
		file, err = os.Open(inputFilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot open specified file")
		}
		defer file.Close()
	}

	// Read the input line by line.
	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Fix some issues with the Terraform output.
		for matchPattern, replaceStruct := range tfOutputFixes {
			previous := line
			line = matchPattern.ReplaceAllString(line, replaceStruct.replacePattern)
			if line != previous && replaceStruct.report {
				log.
					Debug().
					Str("before", previous).
					Str("after", line).
					Msg("Line fixed")
			}
		}

		// Check if the line has to be excluded.
		exclude := false
		for _, matchPattern := range exclusionPatterns {
			if matchPattern.MatchString(line) {
				log.
					Debug().
					Str("line", line).
					Str("pattern", matchPattern.String()).
					Msg("Line filtered out")
				exclude = true
				break
			}
		}
		if exclude {
			continue
		}

		sb.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("Error while reading input")
	}

	// Load the Graphviz graph from the input.
	graph, err := gographviz.Read([]byte(sb.String()))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse the input as a Graphviz graph")
	}

	return graph
}
