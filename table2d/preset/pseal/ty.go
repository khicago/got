package pseal

import (
	"github.com/khicago/got/util/strs"
	"github.com/khicago/got/util/typer"
)

type (
	Type int

	SealFn func(val any) Seal
)

func (ty Type) Name() string {
	return tyNames[ty]
}

func (ty Type) SymAll() []string {
	return tySymbol[ty]
}

func (ty Type) Assert(val any) bool {
	return tyAsserter[ty](val)
}

func (ty Type) SymMatch(sym string) string {
	sym = strs.TrimLower(sym)
	for _, v := range ty.SymAll() {
		if sym == v {
			return sym
		}
	}
	return ""
}

func (ty Type) Default() any {
	return tyDefault[ty]
}

const (
	TyNil Type = iota
	TyAny

	TyPID    // @
	TyID     // ID
	TyBool   // BOOl, BOOLEAN, Y/N, N/Y
	TyInt    // INT, INTEGER
	TyFloat  // FLOAT, NUM
	TyString // STR, TEXT, STRING
	TyMemo   // #, MEMO, MEM

	// TyMark
	// - TyObjectStart  // {
	// - TyObjectEnd    // }
	// - TyListStart    // [
	// - TyListEnd      // ]
	TyMark
)

const (
	DefaultPID    = int64(-1)
	DefaultID     = int64(-1)
	DefaultBool   = false
	DefaultInt    = 0
	DefaultFloat  = 0.0
	DefaultString = ""
	DefaultMemo   = ""
	DefaultMark   = ""
)

var (
	DefaultAny any = nil
)

var tyNames = map[Type]string{
	TyNil:    "nil",
	TyAny:    "any",
	TyPID:    "pid",
	TyID:     "id",
	TyBool:   "bool",
	TyInt:    "int",
	TyFloat:  "float",
	TyString: "string",
	TyMemo:   "memo",
	TyMark:   "mark",
}

var tySymbol = map[Type][]string{
	TyPID:    {"@", "pid"},
	TyID:     {"id"},
	TyBool:   {"bool", "boolean", "n/y", "y/n"},
	TyInt:    {"int", "integer"},
	TyFloat:  {"float", "num"},
	TyString: {"string", "str", "text"},
	TyMemo:   {"#", "memo", "mem"},
	TyMark:   {"[", "]", "{", "}"},

	TyNil: {},
	TyAny: {},
}

var tyDefault = map[Type]any{
	TyNil:    DefaultAny,
	TyAny:    DefaultAny,
	TyPID:    DefaultPID,
	TyID:     DefaultID,
	TyBool:   DefaultBool,
	TyInt:    DefaultInt,
	TyFloat:  DefaultFloat,
	TyString: DefaultString,
	TyMemo:   DefaultMemo,
	TyMark:   DefaultMark,
}

var tyAsserter = map[Type]typer.Predicate[any]{
	TyNil:    func(any) bool { return false },
	TyAny:    func(any) bool { return true },
	TyPID:    typer.AssertType[int64, any],
	TyID:     typer.AssertType[int64, any],
	TyBool:   typer.AssertType[bool, any],
	TyInt:    typer.AssertType[int, any],
	TyFloat:  typer.AssertType[float64, any],
	TyString: typer.AssertType[string, any],
	TyMemo:   typer.AssertType[string, any],
	TyMark:   typer.AssertType[string, any],
}

func SymToType(sym string) Type {
	sym = strs.TrimLower(sym)
	for ty, syms := range tySymbol {
		for _, s := range syms {
			if s == sym {
				return ty
			}
		}
	}
	return TyNil
}
