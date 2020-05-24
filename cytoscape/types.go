package cytoscape

// The complete documentation of the Cytoscape.js elements is available here:
// See https://js.cytoscape.org/#notation/elements-json

// NodeData represents the element data of a graph node.
type NodeData struct {
	// Standard properties
	ID     string  `json:"id"`
	Parent *string `json:"parent,omitempty"`

	// Custom properties.
	Label string `json:"label"`
}

// Node represents a graph node.
type Node struct {
	Data    NodeData `json:"data"`
	Classes []string `json:"classes,omitempty"`
}

// EdgeData represents the element data of a graph edge.
type EdgeData struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// Edge represents a graph edge.
type Edge struct {
	Data    EdgeData `json:"data"`
	Classes []string `json:"classes,omitempty"`
}

// Elements represents a graph elements (nodes and edges).
type Elements struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}
