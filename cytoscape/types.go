package cytoscape

type NodeData struct {
	// Standard properties
	Id     string  `json:"id"`
	Parent *string `json:"parent,omitempty"`

	// Custom properties.
	Label string `json:"label"`
}

type Node struct {
	Data       NodeData `json:"data"`
	Selected   *bool    `json:"selected,omitempty"`
	Selectable *bool    `json:"selectable,omitempty"`
	Locked     *bool    `json:"locked,omitempty"`
	Grabbable  *bool    `json:"grabbable,omitempty"`
	Pannable   *bool    `json:"pannable,omitempty"`
	Classes    []string `json:"classes,omitempty"`
}

type EdgeData struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type Edge struct {
	Data     EdgeData `json:"data"`
	Pannable *bool    `json:"pannable,omitempty"`
	Classes  []string `json:"classes,omitempty"`
}

type Elements struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}
