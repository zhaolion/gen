/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

// Universe is a map of all packages. The key is the package name, but you
// should use Package(), Type(), Function(), or Variable() instead of direct
// access.
type Universe map[string]*Package

// Type returns the canonical type for the given fully-qualified name. Builtin
// types will always be found, even if they haven't been explicitly added to
// the map. If a non-existing type is requested, this will create (a marker for)
// it.
func (u Universe) Type(n Name) *Type {
	return u.Package(n.Package).Type(n.Name)
}

// Function returns the canonical function for the given fully-qualified name.
// If a non-existing function is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the function by setting Underlying to the correct type.
func (u Universe) Function(n Name) *Type {
	return u.Package(n.Package).Function(n.Name)
}

// Variable returns the canonical variable for the given fully-qualified name.
// If a non-existing variable is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the variable by setting Underlying to the correct type.
func (u Universe) Variable(n Name) *Type {
	return u.Package(n.Package).Variable(n.Name)
}

// Constant returns the canonical constant for the given fully-qualified name.
// If a non-existing constant is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the constant by setting Underlying to the correct type.
func (u Universe) Constant(n Name) *Type {
	return u.Package(n.Package).Constant(n.Name)
}

// AddImports registers import lines for packageName. May be called multiple times.
// You are responsible for canonicalizing all package paths.
func (u Universe) AddImports(packagePath string, importPaths ...string) {
	p := u.Package(packagePath)
	for _, i := range importPaths {
		p.Imports[i] = u.Package(i)
	}
}

// Package returns the Package for the given path.
// If a non-existing package is requested, this will create (a marker for) it.
// If a marker is created, it's the caller's responsibility to finish
// construction of the package.
func (u Universe) Package(packagePath string) *Package {
	if p, ok := u[packagePath]; ok {
		return p
	}
	p := &Package{
		Path:      packagePath,
		Types:     map[string]*Type{},
		Functions: map[string]*Type{},
		Variables: map[string]*Type{},
		Constants: map[string]*Type{},
		Imports:   map[string]*Package{},
	}
	u[packagePath] = p
	return p
}

// Package holds package-level information.
// Fields are public, as everything in this package, to enable consumption by
// templates (for example). But it is strongly encouraged for code to build by
// using the provided functions.
type Package struct {
	// Canonical name of this package-- its path.
	Path string

	// The location this package was loaded from
	SourcePath string

	// Short name of this package; the name that appears in the
	// 'package x' line.
	Name string

	// The comment right above the package declaration in doc.go, if any.
	DocComments []string
	// All comments from *.go, if any.
	Comments []string

	// Types within this package, indexed by their name (*not* including
	// package name).
	Types map[string]*Type

	// Functions within this package, indexed by their name (*not* including
	// package name).
	Functions map[string]*Type

	// Global variables within this package, indexed by their name (*not* including
	// package name).
	Variables map[string]*Type

	// Global constants within this package, indexed by their name (*not* including
	// package name).
	Constants map[string]*Type

	// Packages imported by this package, indexed by (canonicalized)
	// package path.
	Imports map[string]*Package
}

// Has returns true if the given name references a type known to this package.
func (p *Package) Has(name string) bool {
	_, has := p.Types[name]
	return has
}

// Type gets the given Type in this Package.  If the Type is not already
// defined, this will add it and return the new Type value.  The caller is
// expected to finish initialization.
func (p *Package) Type(typeName string) *Type {
	if t, ok := p.Types[typeName]; ok {
		return t
	}
	if p.Path == "" {
		// Import the standard builtin types!
		if t, ok := builtins.Types[typeName]; ok {
			p.Types[typeName] = t
			return t
		}
	}
	t := &Type{Name: Name{Package: p.Path, Name: typeName}}
	p.Types[typeName] = t
	return t
}

// Function gets the given function Type in this Package. If the function is
// not already defined, this will add it.  If a function is added, it's the
// caller's responsibility to finish construction of the function by setting
// Underlying to the correct type.
func (p *Package) Function(funcName string) *Type {
	if t, ok := p.Functions[funcName]; ok {
		return t
	}
	t := &Type{Name: Name{Package: p.Path, Name: funcName}}
	t.Kind = DeclarationOf
	p.Functions[funcName] = t
	return t
}

// Variable gets the given variable Type in this Package. If the variable is
// not already defined, this will add it. If a variable is added, it's the caller's
// responsibility to finish construction of the variable by setting Underlying
// to the correct type.
func (p *Package) Variable(varName string) *Type {
	if t, ok := p.Variables[varName]; ok {
		return t
	}
	t := &Type{Name: Name{Package: p.Path, Name: varName}}
	t.Kind = DeclarationOf
	p.Variables[varName] = t
	return t
}

// Constant gets the given constant Type in this Package. If the constant is
// not already defined, this will add it. If a constant is added, it's the caller's
// responsibility to finish construction of the constant by setting Underlying
// to the correct type.
func (p *Package) Constant(constName string) *Type {
	if t, ok := p.Constants[constName]; ok {
		return t
	}
	t := &Type{Name: Name{Package: p.Path, Name: constName}}
	t.Kind = DeclarationOf
	p.Constants[constName] = t
	return t
}

// HasImport returns true if p imports packageName. Package names include the
// package directory.
func (p *Package) HasImport(packageName string) bool {
	_, has := p.Imports[packageName]
	return has
}

// Interfaces return all interface types in this package.
func (p *Package) Interfaces() []*Type {
	interfaces := make([]*Type, 0, len(p.Types))
	for k, v := range p.Types {
		if v.Kind == Interface {
			interfaces = append(interfaces, p.Types[k])
		}
	}

	return interfaces
}

// Structs return all struct types in this package.
func (p *Package) Structs() []*Type {
	result := make([]*Type, 0, len(p.Types))
	for k, v := range p.Types {
		if v.Kind == Struct {
			result = append(result, p.Types[k])
		}
	}

	return result
}

// TODO: support more useful functions
type StructType Type

// TODO: support more useful functions
type InterfaceType Type

// Type represents a subset of possible go types.
type Type struct {
	// There are two general categories of types, those explicitly named
	// and those anonymous. Named ones will have a non-empty package in the
	// name field.
	//
	// An exception: If Kind == DeclarationOf, then this name is the name of a
	// top-level function, variable, or const, and the type can be found in Underlying.
	// We do this to allow the naming system to work against these objects, even
	// though they aren't strictly speaking types.
	Name Name

	// The general kind of this type.
	Kind Kind

	// If there are comment lines immediately before the type definition,
	// they will be recorded here.
	CommentLines []string

	// If there are comment lines preceding the `CommentLines`, they will be
	// recorded here. There are two cases:
	// ---
	// SecondClosestCommentLines
	// a blank line
	// CommentLines
	// type definition
	// ---
	//
	// or
	// ---
	// SecondClosestCommentLines
	// a blank line
	// type definition
	// ---
	SecondClosestCommentLines []string

	// If Kind == Struct
	Members []Member

	// If Kind == Map, Slice, Pointer, or Chan
	Elem *Type

	// If Kind == Map, this is the map's key type.
	Key *Type

	// If Kind == Alias, this is the underlying type.
	// If Kind == DeclarationOf, this is the type of the declaration.
	Underlying *Type

	// If Kind == Interface, this is the set of all required functions.
	// Otherwise, if this is a named type, this is the list of methods that
	// type has. (All elements will have Kind=="Func")
	Methods map[string]*Type

	// If Kind == func, this is the signature of the function.
	Signature *Signature

	// TODO: Add:
	// * channel direction
	// * array length
}

// String returns the name of the type.
func (t *Type) String() string {
	return t.Name.String()
}

// Signature is a function's signature.
type Signature struct {
	// TODO: store the parameter names, not just types.

	// If a method of some type, this is the type it's a member of.
	Receiver   *Type
	Parameters []*Type
	Results    []*Type

	// True if the last in parameter is of the form ...T.
	Variadic bool

	// If there are comment lines immediately before this
	// signature/method/function declaration, they will be recorded here.
	CommentLines []string
}

// A single struct member
type Member struct {
	// The name of the member.
	Name string

	// If the member is embedded (anonymous) this will be true, and the
	// Name will be the type name.
	Embedded bool

	// If there are comment lines immediately before the member in the type
	// definition, they will be recorded here.
	CommentLines []string

	// If there are tags along with this member, they will be saved here.
	Tags string

	// The type of this member.
	Type *Type
}

// String returns the name and type of the member.
func (m Member) String() string {
	return m.Name + " " + m.Type.String()
}

// A type name may have a package qualifier.
type Name struct {
	// Empty if embedded or builtin. This is the package path unless Path is specified.
	Package string
	// The type name.
	Name string
	// An optional location of the type definition for languages that can have disjoint
	// packages and paths.
	Path string
}

// String returns the name formatted as a string.
func (n Name) String() string {
	if n.Package == "" {
		return n.Name
	}
	return n.Package + "." + n.Name
}

// The possible classes of types.
type Kind string

const (
	// Builtin is a primitive, like bool, string, int.
	Builtin Kind = "Builtin"
	Struct  Kind = "Struct"
	Map     Kind = "Map"
	Slice   Kind = "Slice"
	Pointer Kind = "Pointer"

	// Alias is an alias of another type, e.g. in:
	//   type Foo string
	//   type Bar Foo
	// Bar is an alias of Foo.
	//
	// In the real go type system, Foo is a "Named" string; but to simplify
	// generation, this type system will just say that Foo *is* a builtin.
	// We then need "Alias" as a way for us to say that Bar *is* a Foo.
	Alias Kind = "Alias"

	// Interface is any type that could have differing types at run time.
	Interface Kind = "Interface"

	// The remaining types are included for completeness, but are not well
	// supported.
	Array Kind = "Array" // Array is just like slice, but has a fixed length.
	Chan  Kind = "Chan"
	Func  Kind = "Func"

	// DeclarationOf is different from other Kinds; it indicates that instead of
	// representing an actual Type, the type is a declaration of an instance of
	// a type. E.g., a top-level function, variable, or constant. See the
	// comment for Type.Name for more detail.
	DeclarationOf Kind = "DeclarationOf"
	Unknown       Kind = ""
	Unsupported   Kind = "Unsupported"

	// Protobuf is protobuf type.
	Protobuf Kind = "Protobuf"
)
