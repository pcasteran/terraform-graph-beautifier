package tfgraph

type Resource struct {
	// Embed AbstractConfigurationComponent.
	*AbstractConfigurationComponent
}

func NewResource(parent *Module, name string) *Resource {
	return &Resource{
		&AbstractConfigurationComponent{
			Parent: parent,
			Name:   name,
		},
	}
}
