package types


// Built in types.
var (
	String = &Type{
		Name: Name{Name: "string"},
		Kind: Builtin,
	}
	Int64 = &Type{
		Name: Name{Name: "int64"},
		Kind: Builtin,
	}
	Int32 = &Type{
		Name: Name{Name: "int32"},
		Kind: Builtin,
	}
	Int16 = &Type{
		Name: Name{Name: "int16"},
		Kind: Builtin,
	}
	Int = &Type{
		Name: Name{Name: "int"},
		Kind: Builtin,
	}
	Uint64 = &Type{
		Name: Name{Name: "uint64"},
		Kind: Builtin,
	}
	Uint32 = &Type{
		Name: Name{Name: "uint32"},
		Kind: Builtin,
	}
	Uint16 = &Type{
		Name: Name{Name: "uint16"},
		Kind: Builtin,
	}
	Uint = &Type{
		Name: Name{Name: "uint"},
		Kind: Builtin,
	}
	Uintptr = &Type{
		Name: Name{Name: "uintptr"},
		Kind: Builtin,
	}
	Float64 = &Type{
		Name: Name{Name: "float64"},
		Kind: Builtin,
	}
	Float32 = &Type{
		Name: Name{Name: "float32"},
		Kind: Builtin,
	}
	Float = &Type{
		Name: Name{Name: "float"},
		Kind: Builtin,
	}
	Bool = &Type{
		Name: Name{Name: "bool"},
		Kind: Builtin,
	}
	Byte = &Type{
		Name: Name{Name: "byte"},
		Kind: Builtin,
	}

	builtins = &Package{
		Types: map[string]*Type{
			"bool":    Bool,
			"string":  String,
			"int":     Int,
			"int64":   Int64,
			"int32":   Int32,
			"int16":   Int16,
			"int8":    Byte,
			"uint":    Uint,
			"uint64":  Uint64,
			"uint32":  Uint32,
			"uint16":  Uint16,
			"uint8":   Byte,
			"uintptr": Uintptr,
			"byte":    Byte,
			"float":   Float,
			"float64": Float64,
			"float32": Float32,
		},
		Imports: map[string]*Package{},
		Path:    "",
		Name:    "",
	}
)
