package tfgraph

type Module struct {
	*BaseConfigElement

	Children map[string]ConfigElement
}

func NewModule(parent *Module, name string) *Module {
	return &Module{
		BaseConfigElement: &BaseConfigElement{
			parent: parent,
			name:   name,
			tfType: TfModule,
		},
		Children: make(map[string]ConfigElement),
	}
}

func (m *Module) AddChild(e ConfigElement) {
	m.Children[e.GetName()] = e
}
