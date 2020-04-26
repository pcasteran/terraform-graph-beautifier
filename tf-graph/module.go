package tfgraph

type Module struct {
	// Embed AbstractConfigurationComponent.
	*AbstractConfigurationComponent

	Children map[string]ConfigurationComponent
}

func NewModule(parent *Module, name string) *Module {
	return &Module{
		AbstractConfigurationComponent: &AbstractConfigurationComponent{
			Parent: parent,
			Name:   name,
		},
		Children: make(map[string]ConfigurationComponent),
	}
}

func (m *Module) AddChild(c *AbstractConfigurationComponent) {
	c.Parent = m
	m.Children[c.Name] = c
}
