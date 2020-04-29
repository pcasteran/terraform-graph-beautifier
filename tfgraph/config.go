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
// TODO : utile ?
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

var /* const */ tfConfigElementRegexp = regexp.MustCompile(`^"module.root.(.*)"$`)
var /* const */ tfModuleRegexp = regexp.MustCompile(`(module\..*?)\.(.*)`)

type ConfigElement interface {
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
