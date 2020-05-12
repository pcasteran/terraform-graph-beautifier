package cytoscape

type NodeData struct {
	Id     string `json:"id"`
	Parent string `json:"parent"`
}

type Node struct {
	Data       NodeData `json:"data"`
	Selected   bool     `json:"selected"`
	Selectable bool     `json:"selectable"`
	Locked     bool     `json:"locked"`
	Grabbable  bool     `json:"grabbable"`
	Pannable   bool     `json:"pannable"`
	Classes    []string `json:"classes"`
}

func NewNode(parentId, id string) *Node {
	return &Node{
		Data:       NodeData{Id: id, Parent: parentId},
		Selected:   false,
		Selectable: true,
		Locked:     false,
		Grabbable:  true,
		Pannable:   false,
		Classes:    nil,
	}
}

type EdgeData struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type Edge struct {
	Data     EdgeData `json:"data"`
	Pannable bool     `json:"pannable"`
	Classes  []string `json:"classes"`
}

func NewEdge(id, source, target string) *Edge {
	return &Edge{
		Data:     EdgeData{Id: id, Source: source, Target: target},
		Pannable: false,
		Classes:  nil,
	}
}

type Elements struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}
