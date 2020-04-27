package main

import (
	"bufio"
	"github.com/awalterschulze/gographviz"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
)

var /* const */ tfOutputFixes = map[*regexp.Regexp]string{
	// Replace : map["foo"] => map['foo']
	regexp.MustCompile(`\["(.*?)"]`): "['${1}']",
}

var /* const */ tfJunkMatches = []*regexp.Regexp{
	regexp.MustCompile(`"\[root] root"`),
	regexp.MustCompile(`"\[root] meta.count-boundary \(EachMode fixup\)"`),
	regexp.MustCompile(`"\[root] provider\..* \(close\)"`),
}

func loadGraph(inputFilePath string, keepTfJunk bool, excludePatterns []string) *gographviz.Graph {
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

	// Build all the patterns to exclude.
	var exclusionPatterns []*regexp.Regexp
	if !keepTfJunk {
		exclusionPatterns = append(exclusionPatterns, tfJunkMatches...)
	}
	for _, pattern := range excludePatterns {
		exclusionPatterns = append(exclusionPatterns, regexp.MustCompile(pattern))
	}

	// Read the input line by line.
	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line has to be excluded.
		exclude := false
		for _, match := range exclusionPatterns {
			if match.MatchString(line) {
				log.
					Debug().
					Str("line", line).
					Str("pattern", match.String()).
					Msg("Line filtered out")
				exclude = true
				break
			}
		}
		if exclude {
			continue
		}

		// Fix some issues with the Terraform output.
		for match, replace := range tfOutputFixes {
			lineFixed := match.ReplaceAllString(line, replace)
			if line != lineFixed {
				log.
					Debug().
					Str("before", line).
					Str("after", lineFixed).
					Msg("Line fixed")
				line = lineFixed
			}
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
