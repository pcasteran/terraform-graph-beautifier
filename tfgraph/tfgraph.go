package tfgraph

import (
	"regexp"
)

const (
	TfResource = "resource"
	TfVar      = "var"
	TfLocal    = "local"
	TfOutput   = "output"
	TfModule   = "module"
	TfProvider = "provider"
)

// Set of all managed TF types.
type void struct{}

var empty void
var /* const */ ManagedTerraformTypes = map[string]void{
	TfResource: empty,
	TfVar:      empty,
	TfLocal:    empty,
	TfOutput:   empty,
	TfModule:   empty,
	TfProvider: empty,
}

var /* const */ TfConfigElementRegexp = regexp.MustCompile(`^"module.root.(.*)"$`)
var /* const */ TfModuleRegexp = regexp.MustCompile(`(module\..*?)\.(.*)`)

type ConfigElement interface {
	GetParent() *Module
	SetParent(parent *Module)
	GetName() string
	GetQualifiedName() string
	GetTfType() string
}

type BaseConfigElement struct {
	parent *Module
	name   string
	tfType string
}

func NewBaseConfigElement(parent *Module, name string, tfType string) *BaseConfigElement {
	return &BaseConfigElement{
		parent: parent,
		name:   name,
		tfType: tfType,
	}
}

func (e *BaseConfigElement) GetParent() *Module {
	return e.parent
}

func (e *BaseConfigElement) SetParent(parent *Module) {
	e.parent = parent
}

func (e *BaseConfigElement) GetName() string {
	return e.name
}

func (e *BaseConfigElement) GetQualifiedName() string {
	if e.parent != nil {
		return e.parent.GetQualifiedName() + "." + e.name
	} else {
		return e.name
	}
}

func (e *BaseConfigElement) GetTfType() string {
	return e.tfType
}

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

type Dependency struct {
	// This element depends on the `Dst` element.
	Src ConfigElement

	// This element is a dependency of the `Src` element.
	Dst ConfigElement
}

func NewDependency(src ConfigElement, dst ConfigElement) *Dependency {
	return &Dependency{
		Src: src,
		Dst: dst,
	}
}

type Graph struct {
	Root         *Module
	Dependencies []*Dependency
}

func NewGraph(root *Module, dependencies []*Dependency) *Graph {
	return &Graph{
		Root:         root,
		Dependencies: dependencies,
	}
}
