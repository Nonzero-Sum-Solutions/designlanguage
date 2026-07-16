package model

import (
	"errors"
	"fmt"
	"strings"
)

type Design interface {
	Commented

	Author() string
	Namespace() string
	AllComponents() []Component
	Enums() []Enum
	BaseComponents() []Component
	Entities() []Entity
	Objects() []Object
}

type Named interface {
	Name() string
}

type named struct {
	name string
}

func (n *named) Name() string {
	return n.name
}

func newNamed(name string) Named {
	return &named{name}
}

type Param interface {
	Named
	Type() Type
}

type param struct {
	Named
	typ Type
}

func (p *param) Type() Type {
	return p.typ
}

func NewParam(name string, type1 Type) Param {
	return &param{
		newNamed(name),
		type1,
	}
}

type Commented interface {
	Comment() string
}

type Type interface {
	Named
	IsArray() bool
}

// Component is a computational object that may be composed of other components.
// It exposes its functionality through a single interface.
type Component interface {
	Named
	Commented
	Supertype() Type
}

type component struct {
	Named
	Commented
	supertype Type
}

// TODO Validate all strings more stringently for validity
func NewComponent(name, comment string, supertype Type) (Component, error) {
	if strings.Contains(name, " ") {
		return nil, fmt.Errorf("Invalid name: [%s]", name)
	}

	if strings.HasSuffix(comment, " ") {
		return nil, fmt.Errorf("Invalid comment: [%s]", comment)
	}

	return &component{
		newNamed(name),
		newCommented(comment),
		supertype,
	}, nil
}

func (c *component) Supertype() Type {
	return c.supertype
}

// Attribute is the named part of an interface, designating a sub-component.
type Attribute interface {
	Commented
	Param
}

// Method is a means by which a component computes.
// It's a function from some entity to some entity.
type Method interface {
	Named
	Commented
	Params() []Param
	ReturnVals() []Param
}

// Entity is a special type of Component with no methods, and whose
// attributes are all entities.
// It embeds the Component interface to inherit its base functionality.
type Entity interface {
	Component
	Attributes() []Attribute
}

type Object interface {
	Entity
	Methods() []Method
}

// Representation is an encoding of an entity.
type Representation interface {
	// Encode converts the entity into an encoded string representation.
	Encode() string
	// Decode converts an encoded string back into an entity.
	Decode(encoded string) (Entity, error)
}

type Enum interface {
	Component
	Values() []string
}

func NewEnum(name string, values ...string) (Enum, error) {
	if len(values) < 1 {
		return nil, errors.New("Not enough values for enum")
	}
	c, err := NewComponent(name, "", nil)
	if err != nil {
		return nil, err
	}
	return &enumExpr{c, values}, nil
}

type enumExpr struct {
	Component
	values []string
}

func (e *enumExpr) Values() []string {
	return e.values
}

func NewObject(name, comment string, superType Type, attributes []Attribute, methods []Method) (Object, error) {
	e, err := NewEntity(name, comment, superType, attributes)
	if err != nil {
		return nil, err
	}

	return &object{e, methods}, nil
}

type object struct {
	Entity
	methods []Method
}

func (o *object) Methods() []Method {
	return o.methods
}

func NewEntity(name, comment string, supertype Type, attributes []Attribute) (Entity, error) {
	c, err := NewComponent(name, comment, supertype)
	if err != nil {
		return nil, err
	}
	return &entity{c, attributes}, nil
}

type entity struct {
	Component
	attributes []Attribute
}

func (e *entity) Attributes() []Attribute {
	return e.attributes
}

type commented struct {
	comment string
}

func (c *commented) Comment() string {
	return c.comment
}

func newCommented(commentText string) Commented {
	return &commented{commentText}
}

// design is a concrete implementation of the Design interface
type design struct {
	Commented
	author         string
	namespace      string
	allComponents  []Component
	enums          []Enum
	baseComponents []Component
	entities       []Entity
	objects        []Object
}

func (d *design) Author() string {
	return d.author
}

func (d *design) Namespace() string {
	return d.namespace
}

func (d *design) AllComponents() []Component {
	return d.allComponents
}

func (d *design) Enums() []Enum {
	return d.enums
}

func (d *design) BaseComponents() []Component {
	return d.baseComponents
}
func (d *design) Entities() []Entity {
	return d.entities
}

func (d *design) Objects() []Object {
	return d.objects
}

func NewDesign(author, comment, namespace string, allComponents []Component) Design {
	return &design{
		newCommented(comment),
		author,
		namespace,
		allComponents,
		GetEnums(allComponents),
		GetBaseComponents(allComponents),
		GetEntities(allComponents),
		GetObjects(allComponents),
	}
}

// attribute is a concrete implementation of the Attribute interface
type attribute struct {
	Named
	Commented
	typ Type
}

func (a *attribute) Type() Type {
	return a.typ
}

// typ is a concrete implementation of the Type interface
type typ struct {
	name  string
	array bool
}

func (t *typ) Name() string {
	return t.name
}

func (t *typ) IsArray() bool {
	return t.array
}

func NewAttribute(name, comment, typeName string, isArray bool) (Attribute, error) {
	t, err := NewType(typeName, isArray)
	return &attribute{
		newNamed(name),
		newCommented(comment),
		t,
	}, err
}

func NewType(name string, isArray bool) (Type, error) {
	if strings.Contains(name, " ") {
		return nil, fmt.Errorf("Invalid type: [%s]", name)
	}

	return &typ{
		name:  name,
		array: isArray,
	}, nil
}

// method is a concrete implementation of the Method interface
type method struct {
	Named
	Commented
	params     []Param
	returnVals []Param
}

func (m *method) Params() []Param {
	return m.params
}

func (m *method) ReturnVals() []Param {
	return m.returnVals
}

func NewMethod(name, comment string, params, returnVals []Param) (Method, error) {
	if strings.Contains(name, " ") {
		return nil, fmt.Errorf("Invalid method name: [%s]", name)
	}

	if strings.HasSuffix(comment, " ") {
		return nil, fmt.Errorf("Invalid comment: [%s]", comment)
	}

	return &method{
		newNamed(name),
		newCommented(comment),
		params,
		returnVals,
	}, nil
}

func GetEntities(components []Component) []Entity {
	entities := []Entity{}
	for _, component := range components {
		if _, ok := component.(Object); ok {
			continue
		}

		if e, ok := component.(Entity); ok {
			entities = append(entities, e)
		}
	}
	return entities
}

func GetEnums(components []Component) []Enum {
	enums := []Enum{}
	for _, component := range components {
		if e, ok := component.(Enum); ok {
			enums = append(enums, e)
		}
	}
	return enums
}

func GetObjects(components []Component) []Object {
	objects := []Object{}
	for _, component := range components {
		if e, ok := component.(Object); ok {
			objects = append(objects, e)
		}
	}
	return objects
}

func GetBaseComponents(components []Component) []Component {
	filteredComponents := []Component{}
	for _, component := range components {
		switch t := component.(type) {
		case Enum:
		case Object:
		case Entity:
			continue
		default:
			filteredComponents = append(filteredComponents, t)
		}
	}
	return filteredComponents
}
