package symtab

import (
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
)

type SymbolKind int

const (
	SkNone SymbolKind = iota
	SkClass
	SkSubroutine
	SkStatic
	SkField
	SkArg
	SkVar
)

func (sk SymbolKind) String() string {
	switch sk {
	case SkClass:
		return "class"
	case SkSubroutine:
		return "subroutine"
	case SkStatic:
		return "static"
	case SkField:
		return "field"
	case SkArg:
		return "argument"
	case SkVar:
		return "var"
	default: //SkNone
		return "none"
	}
}

type symbolAttr struct {
	typ   string
	kind  SymbolKind
	index int
}

func (sa symbolAttr) String() string {
	return fmt.Sprintf("{%s %s %d}", sa.typ, sa.kind, sa.index)
}

type Symtab struct {
	ClassName       string
	SubroutineName  string
	classTable      map[string]symbolAttr
	subroutineTable map[string]symbolAttr
	indexTable      map[SymbolKind]int
}

func New() *Symtab {
	return &Symtab{
		ClassName:       "",
		classTable:      make(map[string]symbolAttr),
		subroutineTable: make(map[string]symbolAttr),
		indexTable: map[SymbolKind]int{
			SkStatic: 0,
			SkField:  0,
			SkArg:    0,
			SkVar:    0,
		},
	}
}

func (st *Symtab) StartSubroutine() {
	st.subroutineTable = make(map[string]symbolAttr)
	st.indexTable[SkArg] = 0
	st.indexTable[SkVar] = 0
}

func (st *Symtab) Define(name, typ string, kind SymbolKind) {
	if kind == SkStatic || kind == SkField {
		st.classTable[name] = symbolAttr{
			typ:   typ,
			kind:  kind,
			index: st.indexTable[kind],
		}
	} else {
		st.subroutineTable[name] = symbolAttr{
			typ:   typ,
			kind:  kind,
			index: st.indexTable[kind],
		}
	}
	st.indexTable[kind] += 1
}

func (st *Symtab) VarCount(kind SymbolKind) int {
	return st.indexTable[kind]
}

func (st *Symtab) KindOf(name string) SymbolKind {
	{
		v, ok := st.subroutineTable[name]
		if ok {
			return v.kind
		}
	}

	{
		v, ok := st.classTable[name]
		if ok {
			return v.kind
		}
	}

	return SkNone
}

func (st *Symtab) TypeOf(name string) string {
	{
		v, ok := st.subroutineTable[name]
		if ok {
			return v.typ
		}
	}

	{
		v, ok := st.classTable[name]
		if ok {
			return v.typ
		}
	}

	return ""
}

func (st *Symtab) IndexOf(name string) int {
	{
		v, ok := st.subroutineTable[name]
		if ok {
			return v.index
		}
	}

	{
		v, ok := st.classTable[name]
		if ok {
			return v.index
		}
	}

	return -1
}

func (st *Symtab) ClassTable() map[string]symbolAttr {
	return st.classTable
}

func (st *Symtab) SubroutineTable() map[string]symbolAttr {
	return st.subroutineTable
}

func ElmToTyp(elm tokenizer.Element) string {
	switch v := elm.(type) {
	case *tokenizer.Keyword:
		return v.Label
	case *tokenizer.Identifier:
		return v.Label
	default:
		panic(fmt.Errorf("unexpected element %v", v))
	}
}

func KwdToKind(kwd *tokenizer.Keyword) SymbolKind {
	switch kwd.Val() {
	case tokenizer.KwdStatic:
		return SkStatic
	case tokenizer.KwdField:
		return SkField
	case tokenizer.KwdVar:
		return SkVar
	default:
		return SkArg
	}
}
