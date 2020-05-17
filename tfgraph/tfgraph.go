package tfgraph

import (
	"regexp"
)

// List of the managed Terraform types.
const (
	TfResource = "resource"
	TfVar      = "var"
	TfLocal    = "local"
	TfOutput   = "output"
	TfModule   = "module"
	TfProvider = "provider"
)

// ManagedTerraformTypes is a set of all the managed Terraform types.
var /* const */ ManagedTerraformTypes = map[string]interface{}{
	TfResource: nil,
	TfVar:      nil,
	TfLocal:    nil,
	TfOutput:   nil,
	TfModule:   nil,
	TfProvider: nil,
}

// TfConfigElementRegexp allows to match valid Terraform configuration elements and extract their type+name.
var /* const */ TfConfigElementRegexp = regexp.MustCompile(`^"module.root.(.*)"$`)

// TfModuleRegexp allow to match Terraform modules and extract their name.
var /* const */ TfModuleRegexp = regexp.MustCompile(`(module\..*?)\.(.*)`)

// ConfigElement represents a generic Terraform configuration element.
type ConfigElement interface {
	GetParent() *Module
	SetParent(parent *Module)
	GetName() string
	GetQualifiedName() string
	GetTfType() string
}

// BaseConfigElement represents a concrete Terraform configuration element.
type BaseConfigElement struct {
	parent *Module
	name   string
	tfType string
}

// NewBaseConfigElement creates a new module using the specified parent module, name and type.
func NewBaseConfigElement(parent *Module, name string, tfType string) *BaseConfigElement {
	return &BaseConfigElement{
		parent: parent,
		name:   name,
		tfType: tfType,
	}
}

// GetParent returns the parent module of this element.
func (e *BaseConfigElement) GetParent() *Module {
	return e.parent
}

// SetParent sets the parent module of this element.
func (e *BaseConfigElement) SetParent(parent *Module) {
	e.parent = parent
}

// GetName returns the name of this element.
func (e *BaseConfigElement) GetName() string {
	return e.name
}

// GetQualifiedName returns the qualified name ([parentQualifiedName.]name) of this element.
func (e *BaseConfigElement) GetQualifiedName() string {
	if e.parent != nil {
		return e.parent.GetQualifiedName() + "." + e.name
	}
	return e.name
}

// GetTfType returns the Terraform type of this element.
func (e *BaseConfigElement) GetTfType() string {
	return e.tfType
}

// Module represents a Terraform configuration module.
type Module struct {
	*BaseConfigElement

	Children map[string]ConfigElement
}

// NewModule creates a new module using the specified parent module and name.
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

// AddChild adds the specified element to this module.
func (m *Module) AddChild(e ConfigElement) {
	m.Children[e.GetName()] = e
}

// Dependency represents a directed dependency between two Terraform configuration elements.
type Dependency struct {
	// This element depends on the `Destination` element.
	Source ConfigElement

	// This element is a dependency of the `Source` element.
	Destination ConfigElement
}

// NewDependency creates a new dependency using the specified source and destination elements.
func NewDependency(source, destination ConfigElement) *Dependency {
	return &Dependency{
		Source:      source,
		Destination: destination,
	}
}

// Graph represents a Terraform configuration (elements + dependencies).
type Graph struct {
	Root         *Module
	Dependencies []*Dependency
}

// NewGraph creates a new graph using the specified root module and dependencies.
func NewGraph(root *Module, dependencies []*Dependency) *Graph {
	return &Graph{
		Root:         root,
		Dependencies: dependencies,
	}
}
