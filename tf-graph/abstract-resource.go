package tf_graph

type ConfigurationComponent interface {
	GetQualifiedName() string
}

type AbstractConfigurationComponent struct {
	ConfigurationComponent
	Parent *Module
	Name   string
}

func (c *AbstractConfigurationComponent) GetQualifiedName() string {
	if c.Parent != nil {
		return c.Parent.GetQualifiedName() + "." + c.Name
	} else {
		return c.Name
	}
}
