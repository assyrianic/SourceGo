/*
 * srcgo.go
 * 
 * Copyright 2022 Nirari Technologies.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package SrcGo


import (
	"fmt"
	"bytes"
	"strings"
	//"unicode"
	"errors"
	"go/token"
	"go/ast"
	"go/types"
	"go/format"
	"go/constant"
)


const (
	ERR_STR  = "[ERROR]"
	WARN_STR = "[WARNING]"
	FMT_STR  = "%-100s %s\n"
)

type (
	StmtMutator  func(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator, ctxt any)
	BlockMutator func(b *ast.BlockStmt, mutator StmtMutator, ctxt any)
)

func MutateBlock(b *ast.BlockStmt, mutator StmtMutator, ctxt any) {
	for i, stmt := range b.List {
		mutator(&b.List, i, stmt, MutateBlock, ctxt)
	}
}


type AstTransmitter struct {
	Errors  []error
	FileSet  *token.FileSet
	TypeInfo *types.Info
	Userdata  any
	TmpVar    uint
	DebugMode bool
}

func MakeAstTransmitter(fs *token.FileSet, ti *types.Info, userdata any, debug_mode bool) AstTransmitter {
	return AstTransmitter{
		FileSet:   fs,
		TypeInfo:  ti,
		Userdata:  userdata,
		DebugMode: debug_mode,
	}
}

func (a *AstTransmitter) PrintErr(p token.Pos, msg string) {
	a.Errors = append(a.Errors, errors.New("SourceGo :: " + a.FileSet.PositionFor(p, false).String() + ": " + msg))
}

func (a *AstTransmitter) PrintErrs() {
	for _, msg := range a.Errors {
		fmt.Printf("%s\n", msg)
	}
}


// These below are for convenience to help devs make their passes.

func PtrizeExpr(x ast.Expr) *ast.StarExpr {
	ptr := new(ast.StarExpr)
	ptr.X = x
	return ptr
}

func MakeReference(x ast.Expr) *ast.UnaryExpr {
	ref := new(ast.UnaryExpr)
	ref.X = x
	ref.Op = token.AND
	return ref
}

func MakeBasicLit(tok token.Token, value string) *ast.BasicLit {
	bl := new(ast.BasicLit)
	bl.Kind = tok
	bl.Value = value
	return bl
}

func Arrayify(typ ast.Expr, len ast.Expr) *ast.ArrayType {
	a := new(ast.ArrayType)
	a.Len = len
	a.Elt = typ
	return a
}

func MakeIndex(index, x ast.Expr) *ast.IndexExpr {
	i := new(ast.IndexExpr)
	i.Index = index
	i.X = x
	return i
}


func MakeAssign(create bool) *ast.AssignStmt {
	assign := new(ast.AssignStmt)
	assign.TokPos = token.NoPos
	if create {
		assign.Tok = token.DEFINE
	} else {
		assign.Tok = token.ASSIGN
	}
	return assign
}


func MakeVarDecl(names []*ast.Ident, val ast.Expr, typ types.Type, ti *types.Info) *ast.DeclStmt {
	decl_stmt := new(ast.DeclStmt)
	gen_decl := new(ast.GenDecl)
	gen_decl.Tok = token.VAR
	gen_decl.Lparen = token.NoPos
	val_spec := new(ast.ValueSpec)
	for _, name := range names {
		val_spec.Names = append(val_spec.Names, name)
	}
	if val != nil {
		val_spec.Type = ValueToTypeExpr(val, ti)
	} else {
		val_spec.Type = TypeToASTExpr(typ)
	}
	gen_decl.Specs = append(gen_decl.Specs, val_spec)
	decl_stmt.Decl = gen_decl
	return decl_stmt
}


func MakeBitNotExpr(e ast.Expr) *ast.UnaryExpr {
	u := new(ast.UnaryExpr)
	u.Op, u.X = token.XOR, e
	return u
}

func MakeParenExpr(e ast.Expr) *ast.ParenExpr {
	p := new(ast.ParenExpr)
	p.X = e
	return p
}

func GetTypeBase(t types.Type) types.Type {
	switch t := t.(type) {
	case *types.Array:
		return t.Elem()
	case *types.Slice:
		return t.Elem()
	case *types.Pointer:
		return t.Elem()
	default:
		return nil
	}
}

func TypeToASTExpr(typ types.Type) ast.Expr {
	var type_stack []types.Type
	for typ != nil {
		type_stack = append(type_stack, typ)
		typ = GetTypeBase(typ)
	}
	// iterate backwards to build the AST
	x := (ast.Expr)(nil)
	for i := len(type_stack) - 1; i >= 0; i-- {
		switch t := type_stack[i].(type) {
		case *types.Array:
			len_str := fmt.Sprintf("%d", t.Len())
			x = Arrayify(x, MakeBasicLit(token.INT, len_str))
		case *types.Pointer:
			x = PtrizeExpr(x)
		case *types.Basic, *types.Named:
			x = ast.NewIdent(t.String())
		}
	}
	return x
}

// Turn a value expr into a type expr.
func ValueToTypeExpr(val ast.Expr, ti *types.Info) ast.Expr {
	return TypeToASTExpr(ti.TypeOf(val))
}

func InsertToIndex[T any](a []T, index int, value T) []T {
	if len(a) <= index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func IndexOfValue[T comparable](a []T, x T) int {
	for i, n := range a {
		if x==n {
			return i
		}
	}
	return -1
}

func FindStmt(a []ast.Stmt, x ast.Stmt) int {
	for i, n := range a {
		if x==n {
			return i
		}
	}
	return -1
}

func FindParam(fn *ast.FuncDecl, name string) (*ast.Field, int) {
	for _, field := range fn.Type.Params.List {
		for i, iden := range field.Names {
			if iden.Name==name {
				return field, i
			}
		}
	}
	return nil, 0
}


func IsMapType(expr ast.Expr, typ *types.Info) bool {
	if expr != nil {
		switch e := expr.(type) {
		case *ast.SelectorExpr:
			return IsMapType(e.Sel, typ)
		case *ast.IndexExpr:
			return IsMapType(e.X, typ)
		case *ast.Ident:
			for k, v := range typ.Defs {
				if k.Name==e.Name {
					typ_str := v.String()
					is_var := strings.Contains(typ_str, "var ") || strings.Contains(typ_str, "field ")
					if is_var && strings.Contains(typ_str, "map") {
						return true
					}
				}
			}
			return false
		default:
			break
		}
	}
	return false
}

func IsFuncPtr(expr ast.Expr, typ *types.Info) bool {
	if expr != nil {
		switch e := expr.(type) {
		case *ast.SelectorExpr:
			return IsFuncPtr(e.Sel, typ)
		case *ast.IndexExpr:
			return true
		case *ast.Ident:
			for k, v := range typ.Defs {
				if k.Name==e.Name {
					typ_str := v.String()
					is_var := strings.Contains(typ_str, "var ") || strings.Contains(typ_str, "field ")
					if is_var && strings.Contains(typ_str, "func(") {
						return true
					}
				}
			}
			return false
		case *ast.CallExpr:
			return IsFuncPtr(e.Fun, typ)
		default:
			break
		}
	}
	return false
}

func GetFuncName(expr ast.Expr) string {
	if expr != nil {
		switch e := expr.(type) {
		case *ast.SelectorExpr:
			return GetFuncName(e.Sel)
		case *ast.Ident:
			return e.Name
		case *ast.CallExpr:
			return GetFuncName(e.Fun)
		default:
			break
		}
	}
	return ""
}

func CollectFuncNames(file *ast.File) map[string]*ast.FuncDecl {
	funcs := make(map[string]*ast.FuncDecl)
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			funcs[d.Name.Name] = d
		}
	}
	return funcs
}


func ExpandFuncPtrCalls(x *ast.CallExpr, retvals []ast.Expr, retTypes []types.Type, ti *types.Info, currFunc *ast.FuncDecl) []ast.Stmt {
	var new_stmts []ast.Stmt
	Call_StartFunction := new(ast.CallExpr)
	Call_StartFunction.Fun = ast.NewIdent("Call_StartFunction")
	Call_StartFunction.Args = append(Call_StartFunction.Args, ast.NewIdent("nil"))
	Call_StartFunction.Args = append(Call_StartFunction.Args, x.Fun)
	
	Call_StartFunction_stmt := new(ast.ExprStmt)
	Call_StartFunction_stmt.X = Call_StartFunction
	new_stmts = append(new_stmts, Call_StartFunction_stmt)
	
	for _, arg := range x.Args {
		if func_call := MakeFuncPtrArgCall(arg, false, nil, ti); func_call != nil {
			new_stmts = append(new_stmts, func_call)
		}
	}
	
	if retvals != nil && len(retvals) > 0 {
		rets := len(retvals)
		for i:=1; i<rets; i++ {
			_, is_ptr := retvals[i].(*ast.StarExpr)
			func_call := MakeFuncPtrArgCall(retvals[i], !is_ptr, nil, ti)
			if func_call != nil {
				new_stmts = append(new_stmts, func_call)
			} else if retTypes != nil {
				func_call = MakeFuncPtrArgCall(retvals[i], !is_ptr, retTypes[i], ti)
				new_stmts = append(new_stmts, func_call)
			}
		}
		Call_Finish := new(ast.CallExpr)
		Call_Finish.Fun = ast.NewIdent("Call_Finish")
		
		var retval_type ast.Expr
		if field, _ := FindParam(currFunc, retvals[0].(*ast.Ident).Name); field != nil {
			retval_type = field.Type
		}
		
		if _, is_ptr := retval_type.(*ast.StarExpr); is_ptr {
			Call_Finish.Args = append(Call_Finish.Args, retvals[0])
		} else {
			Call_Finish.Args = append(Call_Finish.Args, MakeReference(retvals[0]))
		}
		
		Call_Finish_stmt := new(ast.ExprStmt)
		Call_Finish_stmt.X = Call_Finish
		new_stmts = append(new_stmts, Call_Finish_stmt)
	} else {
		Call_Finish := new(ast.CallExpr)
		Call_Finish.Fun = ast.NewIdent("Call_Finish")
		Call_Finish.Args = nil
		
		Call_Finish_stmt := new(ast.ExprStmt)
		Call_Finish_stmt.X = Call_Finish
		new_stmts = append(new_stmts, Call_Finish_stmt)
	}
	return new_stmts
}

func MakeFuncPtrArgCall(arg ast.Expr, by_ref bool, pretyp types.Type, ti *types.Info) *ast.ExprStmt {
	var typ types.Type
	if typ = ti.TypeOf(arg); typ == nil && pretyp != nil {
		typ = pretyp
	}
	switch t := typ.(type) {
	case *types.Array:
		switch t.Elem().String() {
		case "char":
			Call_PushStringEx := new(ast.CallExpr)
			Call_PushStringEx.Fun = ast.NewIdent("Call_PushStringEx")
			Call_PushStringEx.Args = append(Call_PushStringEx.Args, arg)
			Call_PushStringEx.Args = append(Call_PushStringEx.Args, MakeBasicLit(token.INT, fmt.Sprintf("%d", t.Len())))
			Call_PushStringEx.Args = append(Call_PushStringEx.Args, MakeBasicLit(token.INT, "1"))
			Call_PushStringEx.Args = append(Call_PushStringEx.Args, MakeBasicLit(token.INT, "2"))
			Call_PushString_stmt := new(ast.ExprStmt)
			Call_PushString_stmt.X = Call_PushStringEx
			return Call_PushString_stmt
		default:
			Call_PushArrayEx := new(ast.CallExpr)
			Call_PushArrayEx.Fun = ast.NewIdent("Call_PushArrayEx")
			Call_PushArrayEx.Args = append(Call_PushArrayEx.Args, arg)
			Call_PushArrayEx.Args = append(Call_PushArrayEx.Args, MakeBasicLit(token.INT, fmt.Sprintf("%d", t.Len())))
			Call_PushArrayEx.Args = append(Call_PushArrayEx.Args, MakeBasicLit(token.INT, "1"))
			Call_PushArrayEx_stmt := new(ast.ExprStmt)
			Call_PushArrayEx_stmt.X = Call_PushArrayEx
			return Call_PushArrayEx_stmt
		}
	case *types.Pointer:
		switch t.Elem().String() {
		case "float":
			Call_PushFloatRef := new(ast.CallExpr)
			Call_PushFloatRef.Fun = ast.NewIdent("Call_PushFloatRef")
			if by_ref {
				Call_PushFloatRef.Args = append(Call_PushFloatRef.Args, MakeReference(arg))
			} else {
				Call_PushFloatRef.Args = append(Call_PushFloatRef.Args, arg)
			}
			
			Call_PushFloatRef_stmt := new(ast.ExprStmt)
			Call_PushFloatRef_stmt.X = Call_PushFloatRef
			return Call_PushFloatRef_stmt
		default:
			Call_PushCellRef := new(ast.CallExpr)
			Call_PushCellRef.Fun = ast.NewIdent("Call_PushCellRef")
			if by_ref {
				Call_PushCellRef.Args = append(Call_PushCellRef.Args, MakeReference(arg))
			} else {
				Call_PushCellRef.Args = append(Call_PushCellRef.Args, arg)
			}
			
			Call_PushCellRef_stmt := new(ast.ExprStmt)
			Call_PushCellRef_stmt.X = Call_PushCellRef
			return Call_PushCellRef_stmt
		}
	case *types.Basic:
		switch t.Name() {
		case "string":
			Call_PushString := new(ast.CallExpr)
			Call_PushString.Fun = ast.NewIdent("Call_PushString")
			Call_PushString.Args = append(Call_PushString.Args, arg)
			
			Call_PushString_stmt := new(ast.ExprStmt)
			Call_PushString_stmt.X = Call_PushString
			return Call_PushString_stmt
		case "float", "float32":
			Call_PushFloat := new(ast.CallExpr)
			if by_ref {
				Call_PushFloat.Fun = ast.NewIdent("Call_PushFloatRef")
				Call_PushFloat.Args = append(Call_PushFloat.Args, MakeReference(arg))
			} else {
				Call_PushFloat.Fun = ast.NewIdent("Call_PushFloat")
				Call_PushFloat.Args = append(Call_PushFloat.Args, arg)
			}
			Call_PushFloat_stmt := new(ast.ExprStmt)
			Call_PushFloat_stmt.X = Call_PushFloat
			return Call_PushFloat_stmt
		default:
			Call_PushCell := new(ast.CallExpr)
			if by_ref {
				Call_PushCell.Fun = ast.NewIdent("Call_PushCellRef")
				Call_PushCell.Args = append(Call_PushCell.Args, MakeReference(arg))
			} else {
				Call_PushCell.Fun = ast.NewIdent("Call_PushCell")
				Call_PushCell.Args = append(Call_PushCell.Args, arg)
			}
			Call_PushCell_stmt := new(ast.ExprStmt)
			Call_PushCell_stmt.X = Call_PushCell
			return Call_PushCell_stmt
		}
	}
	return nil
}


func MakeTypeAlias(name string, typ types.Type, strong bool, builtins map[string]types.Object) {
	alias := types.NewTypeName(token.NoPos, nil, name, typ)
	if strong {
		builtins[name] = alias
	}
	types.Universe.Insert(alias)
}

func MakeNamedType(name string, typ types.Type, methods []*types.Func, builtins map[string]types.Object) {
	alias := types.NewTypeName(token.NoPos, nil, name, nil)
	types.NewNamed(alias, typ, methods)
	types.Universe.Insert(alias)
	builtins[name] = alias
}

func MakeIntConst(name string, num int64) {
	types.Universe.Insert(types.NewConst(token.NoPos, nil, name, types.Typ[types.Int], constant.MakeInt64(num)))
}

func MakeIntVar(name string) {
	types.Universe.Insert(types.NewVar(token.NoPos, nil, name, types.Typ[types.Int]))
}

func MakeEnumType(name string, names []string, values []int64) {
	alias := types.NewTypeName(token.NoPos, nil, name, nil)
	named := types.NewNamed(alias, types.Typ[types.Int], nil)
	types.Universe.Insert(alias)
	for i, n := range names {
		types.Universe.Insert(types.NewConst(token.NoPos, nil, n, named, constant.MakeInt64(values[i])))
	}
}

func MakeParams(param_names []string, param_types []types.Type) *types.Tuple {
	var params []*types.Var
	for i, param := range param_names {
		params = append(params, types.NewParam(token.NoPos, nil, param, param_types[i]))
	}
	return types.NewTuple(params...)
}

func MakeRet(param_types []types.Type) *types.Tuple {
	var rets []*types.Var
	for _, param := range param_types {
		rets = append(rets, types.NewParam(token.NoPos, nil, "", param))
	}
	return types.NewTuple(rets...)
}

func MakeFunc(name string, recv *types.Var, params, results *types.Tuple, variadic bool) {
	sig := types.NewSignature(recv, params, results, variadic)
	func_ := types.NewFunc(token.NoPos, nil, name, sig)
	types.Universe.Insert(func_)
}


func PrintAST(n ast.Node) string {
	var ast_str strings.Builder
	ast.Inspect(n, func(n ast.Node) bool {
		if n != nil {
			ast_str.WriteString(fmt.Sprintf("%p - %T:\t\t%+v", n, n, n) + "\n")
		}
		return true
	})
	return ast_str.String()
}

func PrettyPrintAST(n ast.Node) string {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	err := format.Node(&buf, fset, n)
	if err != nil {
		fmt.Println(err)
	}
	return buf.String()
}
