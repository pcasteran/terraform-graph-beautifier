package graphviz

import (
	"bufio"
	"github.com/awalterschulze/gographviz"
	"github.com/rs/zerolog/log"
	"io"
	"regexp"
	"strings"
)

type replaceStruct struct {
	replacePattern string
	report         bool
}

var /* const */ tfOutputFixes = map[*regexp.Regexp]replaceStruct{
	// Replace : `"[root] rsc_type.rsc_name"` => `"module.root.rsc_type.rsc_name"`
	regexp.MustCompile(`"\[root] `): {`"module.root.`, false},
	// Replace : map["foo"] => map['foo']
	regexp.MustCompile(`\["(.*?)"]`): {"['${1}']", true},
	// Replace : x.y.z (expand) => x.y.z
	regexp.MustCompile(`"(.*?) \(expand\)"`): {`"${1}"`, true},
}

var /* const */ tfJunkMatches = []*regexp.Regexp{
	regexp.MustCompile(`"module\.root\.root"`),
	regexp.MustCompile(`"module\.root\.meta\.count-boundary \(EachMode fixup\)"`),
	regexp.MustCompile(` \(close\)"`),
	regexp.MustCompile(`"module\.root.*\.provider\[.*].*"`),
	// Exclude module nodes. They did not exist in TF 12 and there is some code in the following to handle them specifically.
	regexp.MustCompile(`"module\.root.*\.module\.[^\.]+"`),
}

func readGraph(reader io.Reader, keepTfJunk bool, excludePatterns []string) *gographviz.Graph {
	// Build all the patterns to exclude.
	var exclusionPatterns []*regexp.Regexp
	if !keepTfJunk {
		exclusionPatterns = append(exclusionPatterns, tfJunkMatches...)
	}
	for _, pattern := range excludePatterns {
		exclusionPatterns = append(exclusionPatterns, regexp.MustCompile(pattern))
	}

	// Read the input line by line.
	scanner := bufio.NewScanner(reader)
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
					Str("pattern", matchPattern.String()).
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
