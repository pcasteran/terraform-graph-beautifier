package tfgraph

// TODO : rename ConfigElement
type ConfigurationComponent interface {
	// TODO : à virer ?
	GetQualifiedName() string
}

type AbstractConfigurationComponent struct {
	// TODO : utile ?
	//ConfigurationComponent

	// TODO : à virer ?
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
