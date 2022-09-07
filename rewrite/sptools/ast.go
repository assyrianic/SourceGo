package SPTools

import (
	"io"
	"fmt"
	///"time"
	"strings"
	"strconv"
)


type Node interface {
	// line & col.
	Tok() Token
	Span() Span
	aNode()
}

type node struct {
	tok Token
}
func (n *node) Tok() Token { return n.tok }
func (n *node) Span() Span { return n.tok.Span }
func (*node) aNode() {}

func copyPosToNode(n *node, t Token) {
	n.tok = t
}


// top-level plugin.
type Plugin struct {
	Decls []Decl
	node
}

func IsPluginNode(n Node) bool {
	if _, is_pl := n.(*Plugin); is_pl {
		return true
	}
	return false
}


type StorageClassFlags uint16
const (
	IsPublic = StorageClassFlags(1 << iota)
	IsConst
	IsNative
	IsForward
	IsStatic
	IsStock
	IsPrivate
	IsProtected
	IsReadOnly
	IsSealed
	IsVirtual
	MaxStorageClasses
)

var StorageClassToString = [...]string{
	IsPublic: "public",
	IsConst: "const",
	IsNative: "native",
	IsForward: "forward",
	IsStatic: "static",
	IsStock: "stock",
	IsPrivate: "private",
	IsProtected: "protected",
	IsReadOnly: "readonly",
	IsSealed: "sealed",
	IsVirtual: "virtual",
}

func (sc StorageClassFlags) String() string {
	var sb strings.Builder
	for flag := IsPublic; sc != 0 && flag < MaxStorageClasses; flag <<= 1 {
		if sc & flag > 0 {
			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(StorageClassToString[flag])
			sc &^= flag
		}
	}
	return sb.String()
}

func storageClassFromToken(tok Token) StorageClassFlags {
	switch tok.Kind {
	case TKConst:
		return IsConst
	case TKStock:
		return IsStock
	case TKPublic:
		return IsPublic
	case TKPrivate:
		return IsPrivate
	case TKProtected:
		return IsProtected
	case TKStatic:
		return IsStatic
	case TKForward:
		return IsForward
	case TKNative:
		return IsNative
	case TKReadOnly:
		return IsReadOnly
	case TKSealed:
		return IsSealed
	case TKVirtual:
		return IsVirtual
	default:
		return StorageClassFlags(0)
	}
}


// Declarations here.
type (
	Decl interface {
		Node
		aDecl()
	}
	
	BadDecl struct {
		decl
	}
	
	TypeDecl struct {
		Type Spec // anything not Func or Var Decl.
		decl
	}
	
	// name1, name2[n], name3=expr;
	VarDecl struct {
		Type Spec // *TypeSpec
		Names, Inits []Expr
		Dims [][]Expr // a var can have multiple dims, account for em all.
		// valid dim index but empty means [] auto counting.
		// nil index if there was no initializer.
		ClassFlags StorageClassFlags
		decl
	}
	
	// class type name() {}
	// class type name();
	// class type name1() = name2;
	FuncDecl struct {
		RetType Spec // *TypeSpec
		Ident Expr
		Params []Decl // []*VarDecl, *BadDecl if error.
		Body Node // Expr if alias, Stmt if body, nil if ';'.
		ClassFlags StorageClassFlags
		decl
	}
	
	StaticAssert struct {
		A, B Expr
		decl
	}
)
type decl struct{ node }
func (*decl) aDecl() {}

func IsDeclNode(n Node) bool {
	switch n.(type) {
	case *FuncDecl, *VarDecl, *TypeDecl, *BadDecl, *StaticAssert:
		return true
	}
	return false
}


// Specifications here.
// Spec represents a constant or type definition.
type (
	Spec interface {
		Node
		aSpec()
	}
	
	BadSpec struct {
		spec
	}
	
	// enum Name { ... }
	// enum { ... }
	// enum ( op= i ) { ... }
	EnumSpec struct {
		Ident Expr // can be nil.
		Step Expr // can be nil.
		StepOp TokenKind
		Names []Expr
		Values []Expr
		spec
	}
	
	// struct Name { ... }
	// enum struct Name { ... }
	StructSpec struct {
		Ident Expr
		IsEnum bool
		Fields []Decl // []*VarDecl
		Methods []Decl // []*FuncDecl
		spec
	}
	
	// using name;
	UsingSpec struct {
		Namespace Expr
		spec
	}
	
	// type[]&
	TypeSpec struct {
		Type Expr
		Dims int
		IsRef bool
		spec
	}
	
	// methodmap Name [< type] { ... };
	MethodMapSpec struct {
		Ident Expr
		Parent Expr // TypeExpr
		Props []Spec // []*MethodMapPropSpec
		Methods []Spec // []*MethodMapMethodSpec
		Nullable bool
		spec
	}
	
	/* property Type name {
	 *    public get() {}
	 *    public set(Type param) {}
	 *    
	 *    public native get();
	 *    public native set(Type param);
	 * }
	 */
	MethodMapPropSpec struct {
		Type, Ident Expr
		SetterParams []Decl
		GetterBlock, SetterBlock Stmt
		GetterClass, SetterClass StorageClassFlags
		spec
	}
	// public Type name(params) {}
	// public native Type name(params);
	MethodMapMethodSpec struct {
		Impl Decl
		IsCtor bool
		spec
	}
	
	// function type params;
	SignatureSpec struct {
		Type Spec
		Params []Decl // array of []*VarDecl, *BadDecl if error.
		spec
	}
	
	// typedef name = function type params;
	TypeDefSpec struct {
		Ident Expr
		Sig Spec // *SignatureSpec
		spec
	}
	
	// typeset name {};
	TypeSetSpec struct {
		Ident Expr
		Signatures []Spec // []*SignatureSpec
		spec
	}
)
type spec struct{ node }
func (*spec) aSpec() {}

func IsSpecNode(n Node) bool {
	switch n.(type) {
	case *EnumSpec, *StructSpec, *UsingSpec, *TypeSpec:
		return true
	case *MethodMapSpec, *MethodMapPropSpec, *MethodMapMethodSpec:
		return true
	case *TypeDefSpec, *TypeSetSpec, *SignatureSpec:
		return true
	case *BadSpec:
		return true
	}
	return false
}


// statement nodes here.
// Statement syntax write here.
type (
	Stmt interface {
		Node
		aStmt()
	}
	
	BadStmt struct {
		stmt
	}
	
	DeleteStmt struct {
		X Expr
		stmt
	}
	
	// { *Stmts }
	BlockStmt struct {
		Stmts []Stmt
		stmt
	}
	
	// if cond body [else/else if body]
	IfStmt struct {
		Cond Expr
		Then, Else Stmt
		stmt
	}
	
	// while cond body
	WhileStmt struct {
		Cond Expr
		Body Stmt
		Do bool // do-while version.
		stmt
	}
	
	// for [init] ; [cond] ; [post] body 
	ForStmt struct {
		Init Node
		Cond, Post Expr
		Body Stmt
		stmt
	}
	
	// expr ';'
	ExprStmt struct {
		X Expr
		stmt
	}
	
	// switch cond body
	SwitchStmt struct {
		Cases []Stmt
		Default Stmt
		Cond Expr
		stmt
	}
	
	// case a[, b, ...z]:
	CaseStmt struct {
		Case Expr // single or comma expression.
		Body Stmt
		stmt
	}
	
	// break, continue
	FlowStmt struct {
		Kind TokenKind
		stmt
	}
	
	// Type i;
	DeclStmt struct {
		D Decl // VarSpec, TypeDecl, or FuncDecl
		stmt
	}
	
	// return i;
	// return;
	RetStmt struct {
		X Expr
		stmt
	}
	
	// assert a;
	AssertStmt struct {
		X Expr
		stmt
	}
	
	// static_assert(a, b);
	StaticAssertStmt struct {
		A Decl
		stmt
	}
)
type stmt struct{ node }
func (*stmt) aStmt() {}

func IsStmtNode(n Node) bool {
	switch n.(type) {
	case *BlockStmt, *IfStmt, *WhileStmt, *ForStmt, *ExprStmt, *SwitchStmt, *CaseStmt:
		return true
	case *RetStmt, *DeclStmt, *DeleteStmt, *FlowStmt, *AssertStmt, *StaticAssertStmt:
		return true
	case *BadStmt:
		return true
	}
	return false
}


// expression nodes here.
// Expression syntax write here.
type LitKind uint8
const (
	IntLit LitKind = iota
	BoolLit
	FloatLit
	CharLit
	StringLit
)
var LitKindToStr = [...]string{
	IntLit: "int literal",
	FloatLit: "float literal",
	CharLit: "char literal",
	StringLit: "string literal",
	BoolLit: "bool literal",
}


type (
	Expr interface {
		Node
		aExpr()
		Tag() Type
	}
	
	BadExpr struct {
		expr
	}
	
	// function type (args) {}
	FuncLit struct {
		Sig  Spec // *SignatureSpec
		Body Stmt // *BlockStmt
		expr
	}
	
	// { a, b, c }
	BracketExpr struct {
		Exprs []Expr
		expr
	}
	
	// ...
	EllipsesExpr struct {
		expr
	}
	
	// a,b,c
	CommaExpr struct {
		Exprs []Expr
		expr
	}
	
	// a? b : c
	TernaryExpr struct {
		A, B, C Expr
		expr
	}
	
	// a # b [ # c ... z ]
	ChainExpr struct {
		A Expr
		Bs []Expr
		Kinds []TokenKind
		expr
	}
	
	// a # b
	BinExpr struct {
		L, R Expr
		Kind TokenKind
		expr
	}
	
	// <T>
	TypedExpr struct {
		TypeName Token
		expr
	}
	
	// view_as<T>(expr)
	ViewAsExpr struct {
		Type Expr // *TypedExpr
		X Expr
		expr
	}
	
	// ++i, i++ sizeof new
	UnaryExpr struct {
		X Expr
		Kind TokenKind
		Post bool
		expr
	}
	
	// id.name
	FieldExpr struct {
		X, Sel Expr
		expr
	}
	
	// name::id
	NameSpaceExpr struct {
		N, Id Expr
		expr
	}
	
	// a[i]
	IndexExpr struct {
		X, Index Expr
		expr
	}
	
	// .a = expr
	NamedArg struct {
		X Expr
		expr
	}
	
	// f(a,b,...z)
	CallExpr struct {
		ArgList []Expr // nil means no arguments
		Func      Expr
		expr
	}
	
	// this.a
	ThisExpr struct {
		expr
	}
	
	// null
	NullExpr struct {
		expr
	}
	
	// i
	Name struct {
		Value string
		expr
	}
	
	// 1, 1.0, '1', "1"
	BasicLit struct {
		Value string
		Kind  LitKind
		expr
	}
)
type expr struct{
	node
	tag Type
}
func (*expr) aExpr() {}
func (e *expr) Tag() Type {
	return e.tag
}

func IsExprNode(n Node) bool {
	switch n.(type) {
	case *BracketExpr, *EllipsesExpr, *CommaExpr, *TernaryExpr, *BinExpr, *TypedExpr, *ViewAsExpr, *ChainExpr:
		return true
	case *UnaryExpr, *FieldExpr, *NameSpaceExpr, *IndexExpr, *NamedArg, *CallExpr, *ThisExpr, *NullExpr:
		return true
	case *Name, *BasicLit:
		return true
	case *BadExpr:
		return true
	}
	return false
}


func printTabs(c rune, tabs int, w io.Writer) {
	for i := 0; i < tabs; i++ {
		fmt.Fprintf(w, "%c%c", c, c)
	}
}

func PrintNode(n Node, tabs int, w io.Writer) {
	const c = '-'
	printTabs(c, tabs, w)
	
	if IsExprNode(n) {
		fmt.Fprintf(w, "Type of Expr Node: '%s'\n", GetTypeName(n.(Expr).Tag()))
	}
	
	switch ast := n.(type) {
	case nil:
		fmt.Fprintf(w, "nil Node\n")
	case *BadStmt:
		fmt.Fprintf(w, "Bad/Errored Stmt Node:: %q\n", ast.node.tok.ToString())
	case *BadExpr:
		fmt.Fprintf(w, "Bad/Errored Expr Node:: %q\n", ast.node.tok.ToString())
	case *BadSpec:
		fmt.Fprintf(w, "Bad/Errored Spec Node:: %q\n", ast.node.tok.ToString())
	case *BadDecl:
		fmt.Fprintf(w, "Bad/Errored Decl Node:: %q\n", ast.node.tok.ToString())
	case *NullExpr:
		fmt.Fprintf(w, "'null' expr\n")
	case *BasicLit:
		fmt.Fprintf(w, "Basic Lit :: Value: %q - Kind: %q\n", ast.Value, LitKindToStr[ast.Kind])
	case *FuncLit:
		fmt.Fprintf(w, "Function Lit\n")
		PrintNode(ast.Sig, tabs + 1, w)
		PrintNode(ast.Body, tabs + 1, w)
	case *ThisExpr:
		fmt.Fprintf(w, "'this' expr\n")
	case *Name:
		fmt.Fprintf(w, "Ident: '%s'\n", ast.Value)
	case *UnaryExpr:
		fmt.Fprintf(w, "Unary Expr Kind: %q, Post: '%t'\n", TokenToStr[ast.Kind], ast.Post)
		PrintNode(ast.X, tabs + 1, w)
	case *CallExpr:
		fmt.Fprintf(w, "Call Expr\n")
		PrintNode(ast.Func, tabs + 1, w)
		if ast.ArgList != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Call Expr Arg List\n")
			for i := range ast.ArgList {
				PrintNode(ast.ArgList[i], tabs + 1, w)
			}
		}
	case *IndexExpr:
		fmt.Fprintf(w, "Index Expr Obj\n")
		PrintNode(ast.X, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Index Expr Index\n")
		PrintNode(ast.Index, tabs + 1, w)
	case *NameSpaceExpr:
		fmt.Fprintf(w, "Namespace Expr Name\n")
		PrintNode(ast.N, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Namespace Expr Id\n")
		PrintNode(ast.Id, tabs + 1, w)
	case *FieldExpr:
		fmt.Fprintf(w, "Field Expr\n")
		PrintNode(ast.X, tabs + 1, w)
		PrintNode(ast.Sel, tabs + 1, w)
	case *ViewAsExpr:
		fmt.Fprintf(w, "view_as Expr\n")
		PrintNode(ast.Type, tabs + 1, w)
		PrintNode(ast.X, tabs + 1, w)
	case *BinExpr:
		fmt.Fprintf(w, "Binary Expr - Kind: %q\n", TokenToStr[ast.Kind])
		PrintNode(ast.L, tabs + 1, w)
		PrintNode(ast.R, tabs + 1, w)
	case *ChainExpr:
		fmt.Fprintf(w, "Chain Expr\n")
		PrintNode(ast.A, tabs + 1, w)
		for i := range ast.Kinds {
			fmt.Fprintf(w, "Chain Expr Kind: %q\n", TokenToStr[ast.Kinds[i]])
			PrintNode(ast.Bs[i], tabs + 1, w)
		}
	case *TernaryExpr:
		fmt.Fprintf(w, "Ternary Expr\n")
		PrintNode(ast.A, tabs + 1, w)
		PrintNode(ast.B, tabs + 1, w)
		PrintNode(ast.C, tabs + 1, w)
	case *NamedArg:
		fmt.Fprintf(w, "Named Arg Expr\n")
		PrintNode(ast.X, tabs + 1, w)
	case *TypedExpr:
		fmt.Fprintf(w, "Typed Expr - Kind: %q\n", ast.TypeName.String())
	case *CommaExpr:
		fmt.Fprintf(w, "Comma Expr\n")
		for i := range ast.Exprs {
			PrintNode(ast.Exprs[i], tabs + 1, w)
		}
	case *BracketExpr:
		fmt.Fprintf(w, "Bracket Expr\n")
		for i := range ast.Exprs {
			PrintNode(ast.Exprs[i], tabs + 1, w)
		}
	case *EllipsesExpr:
		fmt.Fprintf(w, "Ellipses '...' Expr\n")
	case *RetStmt:
		fmt.Fprintf(w, "Return Statement\n")
		if ast.X != nil {
			PrintNode(ast.X, tabs + 1, w)
		}
	case *IfStmt:
		fmt.Fprintf(w, "If Statement\n")
		PrintNode(ast.Cond, tabs + 1, w)
		PrintNode(ast.Then, tabs + 1, w)
		if ast.Else != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "If Statement :: Else\n")
			PrintNode(ast.Else, tabs + 1, w)
		}
	case *WhileStmt:
		fmt.Fprintf(w, "While Statement: is Do-While? %t\n", ast.Do)
		PrintNode(ast.Cond, tabs + 1, w)
		PrintNode(ast.Body, tabs + 1, w)
	case *ForStmt:
		fmt.Fprintf(w, "For Statement\n")
		if ast.Init != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "For Statement Init\n")
			PrintNode(ast.Init, tabs + 1, w)
		}
		if ast.Cond != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "For Statement Cond\n")
			PrintNode(ast.Cond, tabs + 1, w)
		}
		if ast.Post != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "For Statement Post\n")
			PrintNode(ast.Post, tabs + 1, w)
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "For Statement Body\n")
		PrintNode(ast.Body, tabs + 1, w)
	case *ExprStmt:
		fmt.Fprintf(w, "Expr Statement\n")
		PrintNode(ast.X, tabs + 1, w)
	case *BlockStmt:
		fmt.Fprintf(w, "Block Statement\n")
		for i := range ast.Stmts {
			PrintNode(ast.Stmts[i], tabs + 1, w)
		}
	case *DeleteStmt:
		fmt.Fprintf(w, "Delete Statement\n")
		PrintNode(ast.X, tabs + 1, w)
	case *SwitchStmt:
		fmt.Fprintf(w, "Switch Statement Condition\n")
		PrintNode(ast.Cond, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Switch Statement Cases\n")
		for i := range ast.Cases {
			PrintNode(ast.Cases[i], tabs + 1, w)
		}
		if ast.Default != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Switch Statement Default Case\n")
			PrintNode(ast.Default, tabs + 1, w)
		}
	case *CaseStmt:
		fmt.Fprintf(w, "Case Statement Exprs\n")
		PrintNode(ast.Case, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Case Statement Body\n")
		PrintNode(ast.Body, tabs + 1, w)
	case *FlowStmt:
		fmt.Fprintf(w, "Flow Statement: %q\n", TokenToStr[ast.Kind])
	case *AssertStmt:
		fmt.Fprintf(w, "Assert Statement\n")
		PrintNode(ast.X, tabs + 1, w)
	case *StaticAssertStmt:
		fmt.Fprintf(w, "Static Assert Statement\n")
		PrintNode(ast.A, tabs + 1, w)
	case *DeclStmt:
		fmt.Fprintf(w, "Declaration Statement\n")
		PrintNode(ast.D, tabs + 1, w)
	case *TypeSpec:
		fmt.Fprintf(w, "Type Specification\n")
		PrintNode(ast.Type, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Type Specification Dims:: %d\n", ast.Dims)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Type Specification Is Reference:: %t\n", ast.IsRef)
	case *EnumSpec:
		fmt.Fprintf(w, "Enum Specification\n")
		PrintNode(ast.Ident, tabs + 1, w)
		if ast.Step != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Enum Specification Step:: Op: %s\n", TokenToStr[ast.StepOp])
			PrintNode(ast.Step, tabs + 1, w)
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Enum Specification Names\n")
		for i := range ast.Names {
			PrintNode(ast.Names[i], tabs + 1, w)
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Enum Specification Values\n")
		for i := range ast.Values {
			PrintNode(ast.Values[i], tabs + 1, w)
		}
	case *StructSpec:
		if ast.IsEnum {
			fmt.Fprintf(w, "Enum Struct Specification\n")
		} else {
			fmt.Fprintf(w, "Struct Specification\n")
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Struct Ident\n")
		PrintNode(ast.Ident, tabs + 1, w)
		if len(ast.Fields) > 0 {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Struct Fields\n")
			for i := range ast.Fields {
				PrintNode(ast.Fields[i], tabs + 1, w)
			}
		}
		if len(ast.Methods) > 0 {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Struct Methods\n")
			for i := range ast.Methods {
				PrintNode(ast.Methods[i], tabs + 1, w)
			}
		}
	case *UsingSpec:
		fmt.Fprintf(w, "Using Specification\n")
		PrintNode(ast.Namespace, tabs + 1, w)
	case *SignatureSpec:
		fmt.Fprintf(w, "Function Signature Specification\n")
		PrintNode(ast.Type, tabs + 1, w)
		if len(ast.Params) > 0 {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Function Signature Params\n")
			for i := range ast.Params {
				PrintNode(ast.Params[i], tabs + 1, w)
			}
		}
	case *TypeDefSpec:
		fmt.Fprintf(w, "TypeDef Specification Ident\n")
		PrintNode(ast.Ident, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "TypeDef Specification Signature\n")
		PrintNode(ast.Sig, tabs + 1, w)
	case *TypeSetSpec:
		fmt.Fprintf(w, "Typeset Specification Ident\n")
		PrintNode(ast.Ident, tabs + 1, w)
		if len(ast.Signatures) > 0 {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Typeset Signatures\n")
			for i := range ast.Signatures {
				PrintNode(ast.Signatures[i], tabs + 1, w)
			}
		}
	case *MethodMapSpec:
		fmt.Fprintf(w, "Methodmap Specification Ident\n")
		PrintNode(ast.Ident, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Specification:: Is Nullable? %t\n", ast.Nullable)
		if ast.Parent != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Methodmap Specification Derived-From\n")
			PrintNode(ast.Parent, tabs + 1, w)
		}
		if len(ast.Props) > 0 {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Methodmap Specification Props\n")
			for i := range ast.Props {
				PrintNode(ast.Props[i], tabs + 1, w)
			}
		}
		if len(ast.Methods) > 0 {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Methodmap Specification Methods\n")
			for i := range ast.Methods {
				PrintNode(ast.Methods[i], tabs + 1, w)
			}
		}
	case *MethodMapPropSpec:
		fmt.Fprintf(w, "Methodmap Property Specification Type\n")
		PrintNode(ast.Type, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Property Ident\n")
		PrintNode(ast.Ident, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Property Get Storage Class: '%s'\n", ast.GetterClass.String())
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Property Get Block\n")
		PrintNode(ast.GetterBlock, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Property Set Storage Class: '%s'\n", ast.SetterClass.String())
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Property Set Params\n")
		for i := range ast.SetterParams {
			PrintNode(ast.SetterParams[i], tabs + 1, w)
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Property Set Block\n")
		PrintNode(ast.SetterBlock, tabs + 1, w)
	case *MethodMapMethodSpec:
		fmt.Fprintf(w, "Methodmap Method Specification\n")
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Methodmap Method Specification:: Is Constructor? %t\n", ast.IsCtor)
		PrintNode(ast.Impl, tabs + 1, w)
	case *VarDecl:
		fmt.Fprintf(w, "Var Declaration Type\n")
		PrintNode(ast.Type, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Var Declaration Names\n")
		for i := range ast.Names {
			PrintNode(ast.Names[i], tabs + 1, w)
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Var Declaration Array Dims\n")
		for i := range ast.Dims {
			if ast.Dims[i] != nil {
				for _, dim := range ast.Dims[i] {
					PrintNode(dim, tabs + 1, w)
				}
			}
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Var Declaration Inits\n")
		for i := range ast.Inits {
			PrintNode(ast.Inits[i], tabs + 1, w)
		}
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Var Declaration Flags:: %s\n", ast.ClassFlags.String())
	case *FuncDecl:
		fmt.Fprintf(w, "Func Declaration Type\n")
		PrintNode(ast.RetType, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Func Declaration Name\n")
		PrintNode(ast.Ident, tabs + 1, w)
		printTabs(c, tabs, w)
		fmt.Fprintf(w, "Func Declaration Flags:: %s\n", ast.ClassFlags.String())
		if len(ast.Params) > 0 {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Func Declaration Params\n")
			for i := range ast.Params {
				PrintNode(ast.Params[i], tabs + 1, w)
			}
		}
		if ast.Body != nil {
			printTabs(c, tabs, w)
			fmt.Fprintf(w, "Func Declaration Body\n")
			PrintNode(ast.Body, tabs + 1, w)
		}
	case *TypeDecl:
		fmt.Fprintf(w, "Type Declaration\n")
		PrintNode(ast.Type, tabs + 1, w)
	case *StaticAssert:
		fmt.Fprintf(w, "Static Assert\n")
		PrintNode(ast.A, tabs + 1, w)
		PrintNode(ast.B, tabs + 1, w)
	case *Plugin:
		fmt.Fprintf(w, "Plugin File\n")
		for i := range ast.Decls {
			PrintNode(ast.Decls[i], tabs + 1, w)
		}
	default:
		fmt.Fprintf(w, "default :: %T\n", ast)
	}
}

func Walk(n, parent Node, visitor func(n, parent Node) bool) {
	if !visitor(n, parent) {
		return
	}
	
	switch ast := n.(type) {
	case *FuncLit:
		Walk(ast.Sig, n, visitor)
		Walk(ast.Body, n, visitor)
	case *CallExpr:
		Walk(ast.Func, n, visitor)
		if ast.ArgList != nil {
			for i := range ast.ArgList {
				Walk(ast.ArgList[i], n, visitor)
			}
		}
	case *IndexExpr:
		Walk(ast.X, n, visitor)
		Walk(ast.Index, n, visitor)
	case *NameSpaceExpr:
		Walk(ast.N, n, visitor)
		Walk(ast.Id, n, visitor)
	case *FieldExpr:
		Walk(ast.X, n, visitor)
		Walk(ast.Sel, n, visitor)
	case *UnaryExpr:
		Walk(ast.X, n, visitor)
	case *ViewAsExpr:
		Walk(ast.Type, n, visitor)
		Walk(ast.X, n, visitor)
	case *BinExpr:
		Walk(ast.L, n, visitor)
		Walk(ast.R, n, visitor)
	case *ChainExpr:
		Walk(ast.A, n, visitor)
		for i := range ast.Bs {
			Walk(ast.Bs[i], n, visitor)
		}
	case *TernaryExpr:
		Walk(ast.A, n, visitor)
		Walk(ast.B, n, visitor)
		Walk(ast.C, n, visitor)
	case *NamedArg:
		Walk(ast.X, n, visitor)
	case *CommaExpr:
		for i := range ast.Exprs {
			Walk(ast.Exprs[i], n, visitor)
		}
	case *BracketExpr:
		for i := range ast.Exprs {
			Walk(ast.Exprs[i], n, visitor)
		}
	case *RetStmt:
		if ast.X != nil {
			Walk(ast.X, n, visitor)
		}
	case *IfStmt:
		Walk(ast.Cond, n, visitor)
		Walk(ast.Then, n, visitor)
		if ast.Else != nil {
			Walk(ast.Else, n, visitor)
		}
	case *WhileStmt:
		Walk(ast.Cond, n, visitor)
		Walk(ast.Body, n, visitor)
	case *ForStmt:
		if ast.Init != nil {
			Walk(ast.Init, n, visitor)
		}
		if ast.Cond != nil {
			Walk(ast.Cond, n, visitor)
		}
		if ast.Post != nil {
			Walk(ast.Post, n, visitor)
		}
		Walk(ast.Body, n, visitor)
	case *ExprStmt:
		Walk(ast.X, n, visitor)
	case *BlockStmt:
		for i := range ast.Stmts {
			Walk(ast.Stmts[i], n, visitor)
		}
	case *DeleteStmt:
		Walk(ast.X, n, visitor)
	case *SwitchStmt:
		Walk(ast.Cond, n, visitor)
		for i := range ast.Cases {
			Walk(ast.Cases[i], n, visitor)
		}
		if ast.Default != nil {
			Walk(ast.Default, n, visitor)
		}
	case *CaseStmt:
		Walk(ast.Case, n, visitor)
		Walk(ast.Body, n, visitor)
	case *AssertStmt:
		Walk(ast.X, n, visitor)
	case *StaticAssertStmt:
		Walk(ast.A, n, visitor)
	case *DeclStmt:
		Walk(ast.D, n, visitor)
	case *TypeSpec:
		Walk(ast.Type, n, visitor)
	case *EnumSpec:
		Walk(ast.Ident, n, visitor)
		Walk(ast.Step, n, visitor)
		for i := range ast.Names {
			Walk(ast.Names[i], n, visitor)
		}
		for i := range ast.Values {
			Walk(ast.Values[i], n, visitor)
		}
	case *StructSpec:
		Walk(ast.Ident, n, visitor)
		for i := range ast.Fields {
			Walk(ast.Fields[i], n, visitor)
		}
		for i := range ast.Methods {
			Walk(ast.Methods[i], n, visitor)
		}
	case *UsingSpec:
		Walk(ast.Namespace, n, visitor)
	case *SignatureSpec:
		Walk(ast.Type, n, visitor)
		for i := range ast.Params {
			Walk(ast.Params[i], n, visitor)
		}
	case *TypeDefSpec:
		Walk(ast.Ident, n, visitor)
		Walk(ast.Sig, n, visitor)
	case *TypeSetSpec:
		Walk(ast.Ident, n, visitor)
		for i := range ast.Signatures {
			Walk(ast.Signatures[i], n, visitor)
		}
	case *MethodMapSpec:
		Walk(ast.Ident, n, visitor)
		Walk(ast.Parent, n, visitor)
		for i := range ast.Props {
			Walk(ast.Props[i], n, visitor)
		}
		for i := range ast.Methods {
			Walk(ast.Methods[i], n, visitor)
		}
	case *MethodMapPropSpec:
		Walk(ast.Type, n, visitor)
		Walk(ast.Ident, n, visitor)
		Walk(ast.GetterBlock, n, visitor)
		for i := range ast.SetterParams {
			Walk(ast.SetterParams[i], n, visitor)
		}
		Walk(ast.SetterBlock, n, visitor)
	case *MethodMapMethodSpec:
		Walk(ast.Impl, n, visitor)
	case *VarDecl:
		Walk(ast.Type, n, visitor)
		for i := range ast.Names {
			if ast.Names[i] != nil {
				Walk(ast.Names[i], n, visitor)
			}
		}
		for i := range ast.Dims {
			if ast.Dims[i] != nil {
				for _, dim := range ast.Dims[i] {
					Walk(dim, n, visitor)
				}
			}
		}
		for i := range ast.Inits {
			if ast.Inits[i] != nil {
				Walk(ast.Inits[i], n, visitor)
			}
		}
	case *FuncDecl:
		Walk(ast.RetType, n, visitor)
		Walk(ast.Ident, n, visitor)
		for i := range ast.Params {
			Walk(ast.Params[i], n, visitor)
		}
		if ast.Body != nil {
			Walk(ast.Body, n, visitor)
		}
	case *TypeDecl:
		Walk(ast.Type, n, visitor)
	case *StaticAssert:
		Walk(ast.A, n, visitor)
		Walk(ast.B, n, visitor)
	case *Plugin:
		for i := range ast.Decls {
			Walk(ast.Decls[i], n, visitor)
		}
	}
}



const (
	SP_GENFLAG_NEWLINE   = 1 << iota
	SP_GENFLAG_SEMICOLON = 1 << iota
	SP_GENFLAG_ALL       = -1
)

func ExprToString(e Expr) string {
	var sb strings.Builder
	exprToString(e, &sb)
	return sb.String()
}

func exprToString(e Expr, sb *strings.Builder) {
	switch ast := e.(type) {
	case *BadExpr:
		sb.WriteString("<bad Expr>")
	case *NullExpr:
		sb.WriteString("null")
	case *BasicLit:
		switch ast.Kind {
		case StringLit, CharLit:
			q := Ternary[rune](ast.Kind==StringLit, '"', '\'')
			sb.WriteRune(q)
			for _, c := range []rune(ast.Value) {
				if strconv.IsGraphic(c) {
					sb.WriteRune(c)
				} else {
					s := strconv.QuoteRuneToASCII(c)
					sb.WriteString(s[1 : len(s)-1])
				}
			}
			sb.WriteRune(q)
		default:
			sb.WriteString(ast.Value)
		}
	case *FuncLit:
		specToString(ast.Sig, sb, 0)
		stmtToString(ast.Body, sb, 1)
	case *ThisExpr:
		sb.WriteString("this")
	case *Name:
		sb.WriteString(ast.Value)
	case *UnaryExpr:
		if ast.Post {
			exprToString(ast.X, sb)
			sb.WriteString(TokenToStr[ast.Kind])
		} else {
			switch ast.Kind {
			case TKSizeof:
				sb.WriteString(TokenToStr[ast.Kind])
				sb.WriteRune('(')
				exprToString(ast.X, sb)
				sb.WriteRune(')')
			case TKDefined, TKNew:
				sb.WriteString(TokenToStr[ast.Kind])
				sb.WriteRune(' ')
				exprToString(ast.X, sb)
			default:
				sb.WriteString(TokenToStr[ast.Kind])
				exprToString(ast.X, sb)
			}
		}
	case *CallExpr:
		exprToString(ast.Func, sb)
		sb.WriteRune('(')
		if ast.ArgList != nil {
			for i := range ast.ArgList {
				exprToString(ast.ArgList[i], sb)
				if i+1 != len(ast.ArgList) {
					sb.WriteString(", ")
				}
			}
		}
		sb.WriteRune(')')
	case *IndexExpr:
		exprToString(ast.X, sb)
		sb.WriteRune('[')
		exprToString(ast.Index, sb)
		sb.WriteRune(']')
	case *NameSpaceExpr:
		exprToString(ast.N, sb)
		sb.WriteString("::")
		exprToString(ast.Id, sb)
	case *FieldExpr:
		exprToString(ast.X, sb)
		sb.WriteRune('.')
		exprToString(ast.Sel, sb)
	case *ViewAsExpr:
		sb.WriteString("view_as< ")
		exprToString(ast.Type, sb)
		sb.WriteString(" >(")
		exprToString(ast.X, sb)
		sb.WriteRune(')')
	case *BinExpr:
		// bitwise ops need parentheses.
		switch ast.Kind {
		case TKAnd, TKOr, TKXor, TKShAL, TKShAR, TKShLR:
			sb.WriteRune('(')
			exprToString(ast.L, sb)
			sb.WriteRune(')')
			sb.WriteString(" " + TokenToStr[ast.Kind] + " ")
			sb.WriteRune('(')
			exprToString(ast.R, sb)
			sb.WriteRune(')')
		default:
			exprToString(ast.L, sb)
			sb.WriteString(" " + TokenToStr[ast.Kind] + " ")
			exprToString(ast.R, sb)
		}
	case *ChainExpr:
		exprToString(ast.A, sb)
		for i := range ast.Kinds {
			sb.WriteString(" " + TokenToStr[ast.Kinds[i]] + " ")
			exprToString(ast.Bs[i], sb)
		}
	case *TernaryExpr:
		sb.WriteRune('(')
		exprToString(ast.A, sb)
		sb.WriteString(")? ")
		exprToString(ast.B, sb)
		sb.WriteString(" : ")
		exprToString(ast.C, sb)
	case *NamedArg:
		sb.WriteRune('.')
		exprToString(ast.X, sb)
	case *TypedExpr:
		sb.WriteString(ast.TypeName.Lexeme)
	case *CommaExpr:
		for i := range ast.Exprs {
			exprToString(ast.Exprs[i], sb)
			if i+1 != len(ast.Exprs) {
				sb.WriteString(", ")
			}
		}
	case *BracketExpr:
		sb.WriteString("{ ")
		for i := range ast.Exprs {
			exprToString(ast.Exprs[i], sb)
			if i+1 != len(ast.Exprs) {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(" }")
	case *EllipsesExpr:
		sb.WriteString("...")
	default:
		sb.WriteString("")
	}
}


func writeTabs(sb *strings.Builder, tabs int, tab_rune rune) {
	for i := 0; i < tabs; i++ {
		sb.WriteRune(tab_rune)
	}
}

func StmtToString(e Stmt) string {
	var sb strings.Builder
	stmtToString(e, &sb, 0)
	return sb.String()
}

// TODO: if over 50 statements for one construct, add an ending construct comment.
func stmtToString(e Stmt, sb *strings.Builder, tabs int) {
	const tab_rune = '\t'
	switch ast := e.(type) {
	case *BadStmt:
		sb.WriteString("<bad Stmt>")
	case *RetStmt:
		sb.WriteString("return")
		if ast.X != nil {
			sb.WriteRune(' ')
			exprToString(ast.X, sb)
		}
		sb.WriteRune(';')
	case *IfStmt:
		sb.WriteString("if( ")
		exprToString(ast.Cond, sb)
		sb.WriteString(" ) ")
		if _, isblock := ast.Then.(*BlockStmt); isblock {
			stmtToString(ast.Then, sb, tabs)
		} else {
			sb.WriteString("{\n")
			writeTabs(sb, tabs+1, tab_rune)
			stmtToString(ast.Then, sb, tabs)
			sb.WriteRune('\n')
			writeTabs(sb, tabs, tab_rune)
			sb.WriteRune('}')
		}
		
		if ast.Else != nil {
			sb.WriteString(" else ")
			stmtToString(ast.Else, sb, tabs)
		}
	case *WhileStmt:
		if ast.Do {
			sb.WriteString("do ")
			stmtToString(ast.Body, sb, tabs)
			sb.WriteString(" while( ")
			exprToString(ast.Cond, sb)
			sb.WriteString(" );")
		} else {
			sb.WriteString("while( ")
			exprToString(ast.Cond, sb)
			sb.WriteString(" ) ")
			stmtToString(ast.Body, sb, tabs)
		}
	case *ForStmt:
		sb.WriteString("for( ")
		if ast.Init != nil {
			// could be Decl or Expr
			if _, is_decl := ast.Init.(*VarDecl); is_decl {
				declToString(ast.Init.(Decl), sb, 0, 0)
			} else {
				exprToString(ast.Init.(Expr), sb)
			}
		}
		sb.WriteRune(';')
		if ast.Cond != nil {
			sb.WriteRune(' ')
			exprToString(ast.Cond, sb)
		}
		sb.WriteRune(';')
		if ast.Post != nil {
			sb.WriteRune(' ')
			exprToString(ast.Post, sb)
		}
		sb.WriteString(" ) ")
		stmtToString(ast.Body, sb, tabs)
	case *ExprStmt:
		exprToString(ast.X, sb)
		sb.WriteRune(';')
	case *BlockStmt:
		sb.WriteString("{\n")
		for i := range ast.Stmts {
			writeTabs(sb, tabs + 1, tab_rune)
			stmtToString(ast.Stmts[i], sb, tabs + 1)
			if i + 1 != len(ast.Stmts) {
				sb.WriteRune('\n')
			}
		}
		sb.WriteRune('\n')
		writeTabs(sb, tabs, tab_rune)
		sb.WriteRune('}')
	case *DeleteStmt:
		sb.WriteString("delete ")
		exprToString(ast.X, sb)
		sb.WriteRune(';')
	case *SwitchStmt:
		sb.WriteString("switch( ")
		exprToString(ast.Cond, sb)
		sb.WriteString(" ) {\n")
		for i := range ast.Cases {
			writeTabs(sb, tabs + 1, tab_rune)
			stmtToString(ast.Cases[i], sb, tabs + 1)
			if i+1 != len(ast.Cases) {
				sb.WriteRune('\n')
			}
		}
		
		if ast.Default != nil {
			sb.WriteRune('\n')
			writeTabs(sb, tabs + 1, tab_rune)
			sb.WriteString("default: ")
			stmtToString(ast.Default, sb, tabs+1)
		}
		sb.WriteRune('\n')
		writeTabs(sb, tabs, tab_rune)
		sb.WriteRune('}')
	case *CaseStmt:
		sb.WriteString("case ")
		exprToString(ast.Case, sb)
		sb.WriteString(": ")
		if _, isblock := ast.Body.(*BlockStmt); isblock {
			stmtToString(ast.Body, sb, tabs)
		} else {
			sb.WriteString("{\n")
			writeTabs(sb, tabs + 1, tab_rune)
			stmtToString(ast.Body, sb, tabs)
			sb.WriteRune('\n')
			writeTabs(sb, tabs, tab_rune)
			sb.WriteRune('}')
		}
	case *FlowStmt:
		sb.WriteString(TokenToStr[ast.Kind] + ";")
	case *AssertStmt:
		sb.WriteString("assert( ")
		exprToString(ast.X, sb)
		sb.WriteString(" );")
	case *StaticAssertStmt:
		declToString(ast.A, sb, 0, SP_GENFLAG_SEMICOLON)
	case *DeclStmt:
		declToString(ast.D, sb, tabs, SP_GENFLAG_SEMICOLON)
	default:
		sb.WriteString("")
	}
}


func SpecToString(e Spec) string {
	var sb strings.Builder
	specToString(e, &sb, 0)
	return sb.String()
}

func specToString(e Spec, sb *strings.Builder, tabs int) {
	const tab_rune = '\t'
	switch ast := e.(type) {
	case *BadSpec:
		sb.WriteString("<bad Spec>")
	case *TypeSpec:
		exprToString(ast.Type, sb)
		if ast.Dims > 0 {
			for i := 0; i < ast.Dims; i++ {
				sb.WriteString("[]")
			}
		} else if ast.IsRef {
			sb.WriteString("&")
		}
	case *EnumSpec:
		sb.WriteString("enum ")
		if ast.Ident != nil {
			exprToString(ast.Ident, sb)
			sb.WriteRune(' ')
		}
		if ast.Step != nil {
			sb.WriteString("( ")
			sb.WriteString(TokenToStr[ast.StepOp])
			sb.WriteRune(' ')
			exprToString(ast.Step, sb)
			sb.WriteString(" ) ")
		}
		sb.WriteString("{\n")
		for i := range ast.Names {
			writeTabs(sb, tabs + 1, tab_rune)
			exprToString(ast.Names[i], sb)
			if ast.Values[i] != nil {
				sb.WriteString(" = ")
				exprToString(ast.Values[i], sb)
			}
			if i+1 != len(ast.Names) {
				sb.WriteString(",")
			}
			sb.WriteRune('\n')
		}
		writeTabs(sb, tabs, tab_rune)
		sb.WriteString("};\n")
	case *StructSpec:
		if ast.IsEnum {
			sb.WriteString("enum ")
		}
		sb.WriteString("struct ")
		exprToString(ast.Ident, sb)
		sb.WriteString(" {\n")
		if len(ast.Fields) > 0 {
			for i := range ast.Fields {
				writeTabs(sb, tabs + 1, tab_rune)
				declToString(ast.Fields[i], sb, tabs + 1, SP_GENFLAG_ALL)
			}
		}
		if len(ast.Methods) > 0 {
			sb.WriteRune('\n')
			for i := range ast.Methods {
				writeTabs(sb, tabs + 1, tab_rune)
				declToString(ast.Methods[i], sb, tabs + 1, SP_GENFLAG_NEWLINE)
			}
		}
		writeTabs(sb, tabs, tab_rune)
		sb.WriteRune('}')
		if !ast.IsEnum {
			sb.WriteRune(';')
		}
		sb.WriteRune('\n')
	case *UsingSpec:
		sb.WriteString("using ")
		exprToString(ast.Namespace, sb)
		sb.WriteString(";\n")
	case *SignatureSpec:
		sb.WriteString("function ")
		specToString(ast.Type, sb, 0)
		sb.WriteString(" (")
		if len(ast.Params) > 0 {
			for i := range ast.Params {
				declToString(ast.Params[i], sb, tabs + 1, 0)
				if i+1 != len(ast.Params) {
					sb.WriteString(", ")
				}
			}
		}
		sb.WriteString(")")
	case *TypeDefSpec:
		sb.WriteString("typedef ")
		exprToString(ast.Ident, sb)
		sb.WriteString(" = ")
		specToString(ast.Sig, sb, 0)
		sb.WriteRune(';')
	case *TypeSetSpec:
		sb.WriteString("typeset ")
		exprToString(ast.Ident, sb)
		sb.WriteString(" {\n")
		if len(ast.Signatures) > 0 {
			for i := range ast.Signatures {
				writeTabs(sb, tabs + 1, tab_rune)
				specToString(ast.Signatures[i], sb, tabs + 1)
				sb.WriteString(";\n")
			}
		}
		writeTabs(sb, tabs, tab_rune)
		sb.WriteString("};\n")
	case *MethodMapSpec:
		sb.WriteString("methodmap ")
		exprToString(ast.Ident, sb)
		sb.WriteRune(' ')
		if ast.Nullable {
			sb.WriteString("__nullable__ ")
		}
		if ast.Parent != nil {
			sb.WriteString("< ")
			exprToString(ast.Parent, sb)
			sb.WriteRune(' ')
		}
		sb.WriteString("{\n")
		if len(ast.Props) > 0 {
			for i := range ast.Props {
				writeTabs(sb, tabs + 1, tab_rune)
				specToString(ast.Props[i], sb, tabs + 1)
				sb.WriteRune('\n')
			}
		}
		sb.WriteRune('\n')
		if len(ast.Methods) > 0 {
			for i := range ast.Methods {
				writeTabs(sb, tabs + 1, tab_rune)
				specToString(ast.Methods[i], sb, tabs + 1)
				sb.WriteRune('\n')
			}
		}
		writeTabs(sb, tabs, tab_rune)
		sb.WriteString("};\n")
	case *MethodMapPropSpec:
		sb.WriteString("property ")
		exprToString(ast.Type, sb)
		sb.WriteRune(' ')
		exprToString(ast.Ident, sb)
		sb.WriteString(" {\n")
		wrote_getter := false
		if ast.GetterBlock != nil || ast.GetterClass > 0 {
			writeTabs(sb, tabs + 1, tab_rune)
			sb.WriteString(ast.GetterClass.String())
			sb.WriteString(" get()")
			if ast.SetterBlock==nil {
				sb.WriteRune(';')
			} else {
				sb.WriteRune(' ')
				stmtToString(ast.GetterBlock, sb, tabs + 1)
			}
			wrote_getter = true
		}
		if ast.SetterBlock != nil || ast.SetterClass > 0 || len(ast.SetterParams) > 0 {
			if wrote_getter {
				sb.WriteRune('\n')
			}
			writeTabs(sb, tabs + 1, tab_rune)
			sb.WriteString(ast.GetterClass.String())
			sb.WriteString(" set(")
			for i := range ast.SetterParams {
				declToString(ast.SetterParams[i], sb, 0, 0)
				if i+1 != len(ast.SetterParams) {
					sb.WriteString(", ")
				}
			}
			sb.WriteRune(')')
			if ast.SetterBlock==nil {
				sb.WriteRune(';')
			} else {
				sb.WriteRune(' ')
				stmtToString(ast.SetterBlock, sb, tabs + 1)
			}
		}
		sb.WriteRune('\n')
		writeTabs(sb, tabs, tab_rune)
		sb.WriteRune('}')
	case *MethodMapMethodSpec:
		declToString(ast.Impl, sb, tabs, 0)
	default:
		sb.WriteString("")
	}
}


func DeclToString(e Decl) string {
	var sb strings.Builder
	declToString(e, &sb, 0, SP_GENFLAG_ALL)
	return sb.String()
}

func declToString(e Decl, sb *strings.Builder, tabs, flags int) {
	switch ast := e.(type) {
	case *BadDecl:
		sb.WriteString("<bad Decl>")
	case *VarDecl:
		if ast.ClassFlags > 0 {
			sb.WriteString(ast.ClassFlags.String())
			sb.WriteRune(' ')
		}
		specToString(ast.Type, sb, 0)
		sb.WriteRune(' ')
		for i := range ast.Names {
			exprToString(ast.Names[i], sb)
			if ast.Dims[i] != nil {
				for _, dim := range ast.Dims[i] {
					sb.WriteRune('[')
					exprToString(dim, sb)
					sb.WriteRune(']')
				}
			}
			if ast.Inits[i] != nil {
				sb.WriteString(" = ")
				exprToString(ast.Inits[i], sb)
			}
			if i+1 != len(ast.Names) {
				sb.WriteString(", ")
			}
		}
		if flags & SP_GENFLAG_SEMICOLON > 0 {
			sb.WriteRune(';')
		}
		if flags & SP_GENFLAG_NEWLINE > 0 {
			sb.WriteRune('\n')
		}
	case *FuncDecl:
		if ast.ClassFlags > 0 {
			sb.WriteString(ast.ClassFlags.String())
			sb.WriteRune(' ')
		}
		specToString(ast.RetType, sb, 0)
		sb.WriteRune(' ')
		exprToString(ast.Ident, sb)
		sb.WriteRune('(')
		if len(ast.Params) > 0 {
			for i := range ast.Params {
				declToString(ast.Params[i], sb, 0, 0)
				if i+1 != len(ast.Params) {
					sb.WriteString(", ")
				}
			}
		}
		sb.WriteRune(')')
		if ast.Body != nil {
			if IsExprNode(ast.Body) {
				sb.WriteString(" = ")
				exprToString(ast.Body.(Expr), sb)
			} else if IsStmtNode(ast.Body) {
				sb.WriteRune(' ')
				stmtToString(ast.Body.(Stmt), sb, tabs)
			}
		} else {
			sb.WriteRune(';')
		}
		if flags & SP_GENFLAG_NEWLINE > 0 {
			sb.WriteRune('\n')
		}
	case *TypeDecl:
		specToString(ast.Type, sb, tabs)
	case *StaticAssert:
		sb.WriteString("static_assert( ")
		exprToString(ast.A, sb)
		if ast.B != nil {
			sb.WriteString(", ")
			exprToString(ast.B, sb)
		}
		sb.WriteString(" )")
		if flags & SP_GENFLAG_SEMICOLON > 0 {
			sb.WriteRune(';')
		}
		if flags & SP_GENFLAG_NEWLINE > 0 {
			sb.WriteRune('\n')
		}
	default:
		sb.WriteString("")
	}
}


func PluginToString(e Node) string {
	if plugin, is_plugin := e.(*Plugin); !is_plugin {
		return ""
	} else {
		var sb strings.Builder
		for i := range plugin.Decls {
			declToString(plugin.Decls[i], &sb, 0, SP_GENFLAG_ALL)
		}
		return sb.String()
	}
}

func AstToString(n Node) string {
	switch {
	case IsPluginNode(n):
		return PluginToString(n)
	case IsDeclNode(n):
		return DeclToString(n.(Decl))
	case IsSpecNode(n):
		return SpecToString(n.(Spec))
	case IsStmtNode(n):
		return StmtToString(n.(Stmt))
	case IsExprNode(n):
		return ExprToString(n.(Expr))
	default:
		return ""
	}
}