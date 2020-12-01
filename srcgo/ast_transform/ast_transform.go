/**
 * ast_transform.go
 * 
 * Copyright 2020 Nirari Technologies.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package ASTMod


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


var ASTCtxt struct {
	TypeInfo      *types.Info
	CurrFile      **ast.File
	CurrFunc      *ast.FuncDecl
	FSet          *token.FileSet
	BuiltInTypes  map[string]types.Object
	Err           func(err error)
	RangeIter,TmpVar,TmpFunc uint
}

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


func MakeVarDecl(names []*ast.Ident, val ast.Expr, typ types.Type) *ast.DeclStmt {
	decl_stmt := new(ast.DeclStmt)
	gen_decl := new(ast.GenDecl)
	gen_decl.Tok = token.VAR
	gen_decl.Lparen = token.NoPos
	
	val_spec := new(ast.ValueSpec)
	for _, name := range names {
		val_spec.Names = append(val_spec.Names, name)
	}
	if val != nil {
		val_spec.Type = ValueToTypeExpr(val)
	} else {
		val_spec.Type = TypeToASTExpr(typ)
	}
	gen_decl.Specs = append(gen_decl.Specs, val_spec)
	
	decl_stmt.Decl = gen_decl
	return decl_stmt
}


func MakeBitNotExpr(e ast.Expr) *ast.UnaryExpr {
	u := new(ast.UnaryExpr)
	u.Op = token.XOR
	u.X = e
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
	/// iterate backwards to build the AST
	x := (ast.Expr)(nil)
	for i := len(type_stack) - 1; i >= 0; i-- {
		switch t := type_stack[i].(type) {
			case *types.Array:
				x = Arrayify( x, MakeBasicLit(token.INT, fmt.Sprintf("%d", t.Len())) )
			case *types.Pointer:
				x = PtrizeExpr(x)
			case *types.Basic, *types.Named:
				x = ast.NewIdent(t.String())
		}
	}
	return x
}

/// Turn a value expr into a type expr.
func ValueToTypeExpr(val ast.Expr) ast.Expr {
	typ := ASTCtxt.TypeInfo.TypeOf(val)
	return TypeToASTExpr(typ)
}

func InsertExpr(a []ast.Expr, index int, value ast.Expr) []ast.Expr {
	if len(a) <= index { /// nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) /// index < len(a)
	a[index] = value
	return a
}

func InsertStmt(a []ast.Stmt, index int, value ast.Stmt) []ast.Stmt {
	if len(a) <= index { /// nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) /// index < len(a)
	a[index] = value
	return a
}

func InsertDecl(a []ast.Decl, index int, value ast.Decl) []ast.Decl {
	if len(a) <= index { /// nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) /// index < len(a)
	a[index] = value
	return a
}

func FindExpr(a []ast.Expr, x ast.Expr) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func FindStmt(a []ast.Stmt, x ast.Stmt) int {
	for i, n := range a {
		if x == n {
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


/**
 * Modifies the return values of a function by mutating them into references and moving them to the parameters.
 * Example Go code: func f() (int, float) {}
 * Result  Go code: func f(f_param1 *float) int {}
 */
func MutateRetTypes(retvals **ast.FieldList, curr_params *ast.FieldList, obj_name string) []*ast.Field {
	if *retvals==nil || (*retvals).List==nil {
		return curr_params.List
	}
	
	new_params := make([]*ast.Field, 0)
	for _, param := range curr_params.List {
		new_params = append(new_params, param)
	}
	
	results := len((*retvals).List)
	
	/// multiple different return values.
	if results > 1 {
		for i := 1; i<results; i++ {
			ret := (*retvals).List[i]
			/// if they're named, treat as reference types.
			if ret.Names != nil && len(ret.Names) > 1 {
				ret.Type = PtrizeExpr(ret.Type)
				new_params = append(new_params, ret)
			} else {
				//param_num := len(new_params)
				ret.Names = append(ret.Names, ast.NewIdent(fmt.Sprintf("%s_param%d", obj_name, i)))
				ret.Type = PtrizeExpr(ret.Type)
				new_params = append(new_params, ret)
			}
		}
		(*retvals).List = (*retvals).List[:1]
	} else if results==1 && (*retvals).List[0].Names != nil && len((*retvals).List[0].Names) > 1 {
		/// This condition can happen if there's multiple return values of the same type but they're named!
		(*retvals).List[0].Type = PtrizeExpr((*retvals).List[0].Type)
		new_params = append(new_params, (*retvals).List[0])
		*retvals = nil
	}
	return new_params
}

func IsFuncPtr(expr ast.Expr) bool {
	if expr != nil {
		switch e := expr.(type) {
			case *ast.SelectorExpr:
				return IsFuncPtr(e.Sel)
			case *ast.IndexExpr:
				return true
			case *ast.Ident:
				for k, v := range ASTCtxt.TypeInfo.Defs {
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
				return IsFuncPtr(e.Fun)
			default:
				break
		}
	}
	return false
}


func PrintSrcGoErr(p token.Pos, msg string) {
	ASTCtxt.Err(errors.New("SourceGo :: " + ASTCtxt.FSet.PositionFor(p, false).String() + ": " + msg))
}


func MakeTypeAlias(name string, typ types.Type, strong bool) {
	alias := types.NewTypeName(token.NoPos, nil, name, typ)
	if strong {
		ASTCtxt.BuiltInTypes[name] = alias
	}
	types.Universe.Insert(alias)
}

func MakeNamedType(name string, typ types.Type, methods []*types.Func) {
	alias := types.NewTypeName(token.NoPos, nil, name, nil)
	types.NewNamed(alias, typ, methods)
	types.Universe.Insert(alias)
	ASTCtxt.BuiltInTypes[name] = alias
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


func AddSrcGoTypes() {
	/**
	 * func NewTypeName(pos token.Pos, pkg *Package, name string, typ Type) *TypeName
	 * 
	 * NewTypeName returns a new type name denoting the given typ. The remaining arguments set the attributes found with all Objects.
	 * The typ argument may be a defined (Named) type or an alias type. It may also be nil such that the returned TypeName can be used as argument for NewNamed, which will set the TypeName's type as a side-effect.
	 * 
	 * 
	 * func NewNamed(obj *TypeName, underlying Type, methods []*Func) *Named
	 * 
	 * NewNamed returns a new named type for the given type name, underlying type, and associated methods. If the given type name obj doesn't have a type yet, its type is set to the returned named type. The underlying type must not be a *Named
	 */
	ASTCtxt.BuiltInTypes = make(map[string]types.Object)
	//MakeTypeAlias("float", types.Typ[types.Float64], false)
	MakeNamedType("Handle", types.Typ[types.UnsafePointer], nil)
	MakeNamedType("__function__", types.Typ[types.UnsafePointer], nil)
}

func SetUpSrcGo(fset *token.FileSet, info *types.Info, err_fn func(err error)) {
	ASTCtxt.TypeInfo = info
	ASTCtxt.Err = err_fn
	ASTCtxt.FSet = fset
}


func AnalyzeIllegalCode(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch x := n.(type) {
				case *ast.FuncDecl:
					if x.Recv != nil && len(x.Recv.List) > 1 {
						PrintSrcGoErr(x.Pos(), "Multiple Receivers are not allowed.")
					}
					if x.Type.Results != nil {
						for _, ret := range x.Type.Results.List {
							if ptr, is_ptr := ret.Type.(*ast.StarExpr); is_ptr {
								PrintSrcGoErr(ptr.Pos(), "Returning Pointers isn't Allowed." + fmt.Sprintf(" Param %v is of pointer type", ret.Names))
							}
						}
					}
					
				case *ast.FuncType:
					if x.Results != nil {
						for _, ret := range x.Results.List {
							if ptr, is_ptr := ret.Type.(*ast.StarExpr); is_ptr {
								PrintSrcGoErr(ptr.Pos(), "Returning Pointers isn't Allowed." + fmt.Sprintf(" Param %v is of pointer type", ret.Names))
							}
						}
					}
				
				case *ast.StructType:
					for _, f := range x.Fields.List {
						switch t := f.Type.(type) {
							case *ast.StarExpr:
								PrintSrcGoErr(t.Pos(), "Pointers are not allowed in Structs.")
							case *ast.ArrayType:
								if t.Len==nil {
									PrintSrcGoErr(t.Pos(), "Arrays of unknown size are not allowed in Structs.")
								}
						}
					}
				case *ast.BranchStmt:
					if x.Tok==token.GOTO || x.Tok==token.FALLTHROUGH {
						PrintSrcGoErr(x.Pos(), fmt.Sprintf(" %s is Illegal.", x.Tok.String()))
					} else if x.Label != nil {
						PrintSrcGoErr(x.Pos(), "Branched Labels are Illegal.")
					}
				
				case *ast.CommClause:
					PrintSrcGoErr(x.Pos(), "Comm Select Cases are Illegal.")
				case *ast.DeferStmt:
					PrintSrcGoErr(x.Pos(), "Defer Statements are Illegal.")
				case *ast.TypeSwitchStmt:
					PrintSrcGoErr(x.Pos(), "Type-Switches are Illegal.")
				case *ast.LabeledStmt:
					PrintSrcGoErr(x.Pos(), "Labels are Illegal.")
				case *ast.GoStmt:
					PrintSrcGoErr(x.Pos(), "Goroutines are Illegal.")
				case *ast.SelectStmt:
					PrintSrcGoErr(x.Pos(), "Select Statements are Illegal.")
				case *ast.SendStmt:
					PrintSrcGoErr(x.Pos(), "Send Statements are Illegal.")
				
				case *ast.BasicLit:
					if x.Kind==token.IMAG {
						PrintSrcGoErr(x.Pos(), "Imaginary Numbers are Illegal.")
					}
				case *ast.TypeAssertExpr:
					PrintSrcGoErr(x.Pos(), "Type Assertions are Illegal.")
				case *ast.SliceExpr:
					PrintSrcGoErr(x.Pos(), "Slice Expressions are Illegal.")
			}
		}
		return true
	})
}


func MergeRetVals(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch f := n.(type) {
				case *ast.FuncDecl:
					f.Type.Params.List = MutateRetTypes(&f.Type.Results, f.Type.Params, f.Name.Name)
				case *ast.TypeSpec:
					if t, is_func_type := f.Type.(*ast.FuncType); is_func_type {
						t.Params.List = MutateRetTypes(&t.Results, t.Params, f.Name.Name)
					}
			}
		}
		return true
	})
}

func ChangeRecvrNames(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch f := n.(type) {
				case *ast.FuncDecl:
					if f.Recv != nil && f.Recv.List[0].Names != nil && len(f.Recv.List[0].Names) > 0 {
						recvr := f.Recv.List[0].Names[0].Name
						ast.Inspect(f.Body, func(n ast.Node) bool {
							if n != nil {
								switch i := n.(type) {
									case *ast.Ident:
										if recvr==i.Name {
											i.Name = "this"
										}
								}
							}
							return true
						})
					}
			}
		}
		return true
	})
}

func MakeFuncPtrArgCall(arg ast.Expr, by_ref bool, pretyp types.Type) *ast.ExprStmt {
	var typ types.Type
	if typ = ASTCtxt.TypeInfo.TypeOf(arg); typ == nil && pretyp != nil {
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

func ExpandFuncPtrCalls(x *ast.CallExpr, retvals []ast.Expr, retTypes []types.Type) []ast.Stmt {
	new_stmts := make([]ast.Stmt, 0)
	Call_StartFunction := new(ast.CallExpr)
	Call_StartFunction.Fun = ast.NewIdent("Call_StartFunction")
	Call_StartFunction.Args = append(Call_StartFunction.Args, ast.NewIdent("nil"))
	Call_StartFunction.Args = append(Call_StartFunction.Args, x.Fun)
	
	Call_StartFunction_stmt := new(ast.ExprStmt)
	Call_StartFunction_stmt.X = Call_StartFunction
	new_stmts = append(new_stmts, Call_StartFunction_stmt)
	
	for _, arg := range x.Args {
		if func_call := MakeFuncPtrArgCall(arg, false, nil); func_call != nil {
			new_stmts = append(new_stmts, func_call)
		}
	}
	
	if retvals != nil && len(retvals) > 0 {
		rets := len(retvals)
		for i:=1; i<rets; i++ {
			_, is_ptr := retvals[i].(*ast.StarExpr)
			func_call := MakeFuncPtrArgCall(retvals[i], !is_ptr, nil)
			if func_call != nil {
				new_stmts = append(new_stmts, func_call)
			} else if retTypes != nil {
				func_call = MakeFuncPtrArgCall(retvals[i], !is_ptr, retTypes[i])
				new_stmts = append(new_stmts, func_call)
			}
		}
		Call_Finish := new(ast.CallExpr)
		Call_Finish.Fun = ast.NewIdent("Call_Finish")
		
		var retval_type ast.Expr
		if field, _ := FindParam(ASTCtxt.CurrFunc, retvals[0].(*ast.Ident).Name); field != nil {
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


type (
	StmtMutator  func(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator)
	BlockMutator func(b *ast.BlockStmt, mutator StmtMutator)
)

func MutateRets(file *ast.File) {
	/** Case Studies of Returning statements to transform:
	 * 
	 * return m3() /// func m3() (type, type, type)
	 * 
	 * return int, float, m1() /// func m1() type
	 */
	for _, decl := range file.Decls {
		switch d := decl.(type) {
			case *ast.FuncDecl:
				ASTCtxt.CurrFunc = d
				if d.Body != nil {
					MutateBlock(d.Body, MutateRetStmts)
				}
				ASTCtxt.CurrFunc = nil
		}
	}
}

func MutateAssigns(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
			case *ast.FuncDecl:
				ASTCtxt.CurrFunc = d
				if d.Body != nil {
					fmt.Printf("MutateAssigns :: Modifying function %s\n", d.Name.Name)
					MutateBlock(d.Body, MutateAssignStmts)
				}
				ASTCtxt.CurrFunc = nil
		}
	}
}

func MutateRanges(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
			case *ast.FuncDecl:
				ASTCtxt.CurrFunc = d
				if d.Body != nil {
					MutateBlock(d.Body, MutateRangeStmts)
				}
				ASTCtxt.CurrFunc = nil
				ASTCtxt.RangeIter = 0
		}
	}
}

func MutateNoRetCalls(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
			case *ast.FuncDecl:
				ASTCtxt.CurrFunc = d
				if d.Body != nil {
					MutateBlock(d.Body, MutateNoRetCallStmts)
				}
				ASTCtxt.CurrFunc = nil
		}
	}
}

func NameAnonFuncs(file **ast.File) {
	/**
	 * Function Literals can be represented in different ways:
	 * 
	 * f := func(params){
	 *     code
	 * }
	 * f(args)
	 * 
	 * OR
	 * 
	 * func(params){
	 *     code
	 * }(args)
	 */
	ASTCtxt.CurrFile = file
	for _, decl := range (*file).Decls {
		switch d := decl.(type) {
			case *ast.FuncDecl:
				ASTCtxt.CurrFunc = d
				if d.Body != nil {
					MutateBlock(d.Body, MutateFuncLit)
				}
				ASTCtxt.CurrFunc = nil
		}
	}
	ASTCtxt.CurrFile = nil
}

func MutateBlock(b *ast.BlockStmt, mutator StmtMutator) {
	for i, stmt := range b.List {
		mutator(&b.List, i, stmt, MutateBlock)
	}
}

func MutateRetStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator) {
	switch n := s.(type) {
		case *ast.BlockStmt:
			bm(n, MutateRetStmts)
		
		case *ast.ForStmt:
			if n.Init != nil {
				MutateRetStmts(owner_list, index, n.Init, bm)
			}
			if n.Post != nil {
				MutateRetStmts(owner_list, index, n.Post, bm)
			}
			bm(n.Body, MutateRetStmts)
		
		case *ast.IfStmt:
			if n.Init != nil {
				MutateRetStmts(owner_list, index, n.Init, bm)
			}
			bm(n.Body, MutateRetStmts)
			if n.Else != nil {
				MutateRetStmts(owner_list, index, n.Else, bm)
			}
		
		case *ast.ReturnStmt:
			index := FindStmt(*owner_list, s)
			res_len := len(n.Results)
			func_calls := make([]*ast.CallExpr, 0)
			for i := range n.Results {
				switch call := n.Results[i].(type) {
					case *ast.CallExpr:
						func_calls = append(func_calls, call)
					default:
						func_calls = append(func_calls, nil)
				}
			}
			
			res := len(n.Results)
			for i:=1; i<res; i++ {
				call := func_calls[i]
				if call != nil && IsFuncPtr(call.Fun) {
					calls := ExpandFuncPtrCalls(call, []ast.Expr{ast.NewIdent(fmt.Sprintf("%s_param%d", ASTCtxt.CurrFunc.Name.Name, i))}, nil)
					for i:=len(calls)-1; i>=0; i-- {
						*owner_list = InsertStmt(*owner_list, index, calls[i])
					}
				} else {
					ptr_deref := PtrizeExpr(ast.NewIdent(fmt.Sprintf("%s_param%d", ASTCtxt.CurrFunc.Name.Name, i)))
					assign := MakeAssign(false)
					assign.Lhs = append(assign.Lhs, ptr_deref)
					assign.Rhs = append(assign.Rhs, n.Results[i])
					*owner_list = InsertStmt(*owner_list, index, assign)
					index++
				}
			}
			
			if res_len > 1 {
				n.Results = n.Results[:1]
				MutateRetStmts(owner_list, index, s, bm)
			} else if res_len==1 {
				/// check for function call.
				ast.Inspect(n.Results[0], func(node ast.Node) bool {
					if node != nil {
						switch call := node.(type) {
							case *ast.CallExpr:
								if IsFuncPtr(call.Fun) {
									ret_tmp := ast.NewIdent(fmt.Sprintf("fptr_temp%d", ASTCtxt.TmpVar))
									ASTCtxt.TmpVar++
									declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, n.Results[0], nil)
									calls := ExpandFuncPtrCalls(call, []ast.Expr{ret_tmp}, nil)
									calls = InsertStmt(calls, 0, declstmt)
									for i:=len(calls)-1; i>=0; i-- {
										*owner_list = InsertStmt(*owner_list, index, calls[i])
									}
									n.Results[0] = ret_tmp
								} else {
									for _, param := range ASTCtxt.CurrFunc.Type.Params.List {
										for _, name := range param.Names {
											if strings.HasPrefix(name.Name, ASTCtxt.CurrFunc.Name.Name + "_param") {
												call.Args = append(call.Args, name)
											}
										}
									}
								}
						}
					}
					return true
				})
			}
		
		case *ast.SwitchStmt:
			MutateRetStmts(owner_list, index, n.Init, bm)
			bm(n.Body, MutateRetStmts)
		
		case *ast.CaseClause:
			for i, stmt := range n.Body {
				MutateRetStmts(&n.Body, i, stmt, bm)
			}
		
		case *ast.RangeStmt:
			bm(n.Body, MutateRetStmts)
	}
}

func MutateAssignStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator) {
	fmt.Printf("MutateAssignStmts :: %T\n", s)
	switch n := s.(type) {
		case *ast.BlockStmt:
			bm(n, MutateAssignStmts)
		
		case *ast.ForStmt:
			if n.Init != nil {
				MutateAssignStmts(owner_list, index, n.Init, bm)
			}
			if n.Post != nil {
				MutateAssignStmts(owner_list, index, n.Post, bm)
			}
			bm(n.Body, MutateAssignStmts)
		
		case *ast.IfStmt:
			if n.Init != nil {
				MutateAssignStmts(owner_list, index, n.Init, bm)
			}
			bm(n.Body, MutateAssignStmts)
			if n.Else != nil {
				MutateAssignStmts(owner_list, index, n.Else, bm)
			}
		
		case *ast.SwitchStmt:
			MutateAssignStmts(owner_list, index, n.Init, bm)
			bm(n.Body, MutateAssignStmts)
		
		case *ast.CaseClause:
			for i, stmt := range n.Body {
				MutateAssignStmts(&n.Body, i, stmt, bm)
			}
		
		case *ast.RangeStmt:
			bm(n.Body, MutateAssignStmts)
		
		case *ast.AssignStmt:
			/**
			 * a,b,c := f()
			 * 
			 * if not func ptr:
			 *     var a,b,c type
			 *     a = f(&b, &c)
			 * 
			 * if func ptr:
			 *     var a,b,c type
			 *     Call_StartFunction(nil, f)
			 *     Call_PushCellRef(&b)
			 *     Call_PushCellRef(&c)
			 *     Call_Finish(&a)
			 */
			left_len, rite_len := len(n.Lhs), len(n.Rhs)
			fn, is_func_call := n.Rhs[0].(*ast.CallExpr)
			if rite_len==1 && left_len >= rite_len && is_func_call {
				fmt.Printf("ast.AssignStmt Mutator Func Call :: %+v\n", fn.Fun)
				/// a func call returning multiple items as a decl + init.
				switch n.Tok {
					case token.DEFINE:
						decl_stmt := new(ast.DeclStmt)
						gen_decl := new(ast.GenDecl)
						gen_decl.Tok = token.VAR
						
						/// first we get each name of a var and then map them to a type.
						var_map := make(map[types.Type][]ast.Expr)
						for _, e := range n.Lhs {
							if type_expr := ASTCtxt.TypeInfo.TypeOf(e); type_expr != nil {
								var_map[type_expr] = append(var_map[type_expr], e)
							} else {
								PrintSrcGoErr(n.TokPos, "Failed to expand assignment statement.")
							}
						}
						
						for key, val := range var_map {
							val_spec := new(ast.ValueSpec)
							for _, name := range val {
								val_spec.Names = append(val_spec.Names, name.(*ast.Ident))
							}
							val_spec.Type = TypeToASTExpr(key)
							gen_decl.Specs = append(gen_decl.Specs, val_spec)
						}
						
						decl_stmt.Decl = gen_decl
						*owner_list = InsertStmt(*owner_list, 0, decl_stmt)
						n.Tok = token.ASSIGN
						MutateAssignStmts(owner_list, index+1, n, bm)
					
					case token.ASSIGN:
						if IsFuncPtr(fn) {
							ret_tmp := ast.NewIdent(fmt.Sprintf("fptr_temp%d", ASTCtxt.TmpVar))
							ASTCtxt.TmpVar++
							declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, n.Lhs[0], nil)
							
							retvals := make([]ast.Expr, 0)
							retvals = append(retvals, ret_tmp)
							for i:=1; i<left_len; i++ {
								retvals = append(retvals, n.Lhs[i])
							}
							calls := ExpandFuncPtrCalls(fn, retvals, nil)
							calls = InsertStmt(calls, 0, declstmt)
							for i := len(calls)-1; i>=0; i-- {
								*owner_list = InsertStmt(*owner_list, index, calls[i])
							}
							n.Lhs = n.Lhs[:1]
							n.Rhs[0] = ret_tmp
						} else {
							/// transform the tuple return into a single return + pass by ref.
							for i:=1; i<left_len; i++ {
								switch e := n.Lhs[i].(type) {
									case *ast.Ident:
										fn.Args = append(fn.Args, MakeReference(e))
								}
							}
							if left_len > 1 {
								n.Lhs = n.Lhs[:1]
							}
						}
				}
			}
	}
}

func MutateRangeStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator) {
	switch n := s.(type) {
		case *ast.BlockStmt:
			bm(n, MutateRangeStmts)
		
		case *ast.EmptyStmt:
		
		case *ast.ForStmt:
			if n.Init != nil {
				MutateRangeStmts(owner_list, index, n.Init, bm)
			}
			if n.Post != nil {
				MutateRangeStmts(owner_list, index, n.Post, bm)
			}
			bm(n.Body, MutateRangeStmts)
		
		case *ast.IfStmt:
			if n.Init != nil {
				MutateRangeStmts(owner_list, index, n.Init, bm)
			}
			bm(n.Body, MutateRangeStmts)
			if n.Else != nil {
				MutateRangeStmts(owner_list, index, n.Else, bm)
			}
		
		case *ast.SwitchStmt:
			MutateRangeStmts(owner_list, index, n.Init, bm)
			bm(n.Body, MutateRangeStmts)
		
		case *ast.CaseClause:
			for i, stmt := range n.Body {
				MutateRangeStmts(&n.Body, i, stmt, bm)
			}
		
		case *ast.RangeStmt:
			if n.Key != nil {
				if iden, ok := n.Key.(*ast.Ident); ok && iden.Name=="_" {
					n.Key = ast.NewIdent(fmt.Sprintf("%s_iter%d", ASTCtxt.CurrFunc.Name.Name, ASTCtxt.RangeIter))
					ASTCtxt.RangeIter++
				}
			} else {
				n.Key = ast.NewIdent(fmt.Sprintf("%s_iter%d", ASTCtxt.CurrFunc.Name.Name, ASTCtxt.RangeIter))
				ASTCtxt.RangeIter++
			}
			
			if n.Value != nil {
				if iden, ok := n.Value.(*ast.Ident); ok && iden.Name=="_" {
					n.Value = nil
				} else {
					switch n.Tok {
						case token.DEFINE:
							decl_stmt := MakeVarDecl([]*ast.Ident{n.Value.(*ast.Ident)}, n.Value, nil)
							n.Body.List = InsertStmt(n.Body.List, 0, decl_stmt)
							
							assign := MakeAssign(false)
							assign.Lhs = append(assign.Lhs, n.Value)
							
							get_index := MakeIndex(n.Key, n.X)
							assign.Rhs = append(assign.Rhs, get_index)
							
							n.Body.List = InsertStmt(n.Body.List, 1, assign)
							n.Value = nil
					}
				}
			}
			bm(n.Body, MutateRangeStmts)
	}
}

func MutateNoRetCallStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator) {
	switch n := s.(type) {
		case *ast.BlockStmt:
			bm(n, MutateNoRetCallStmts)
		
		case *ast.ForStmt:
			if n.Init != nil {
				MutateNoRetCallStmts(owner_list, index, n.Init, bm)
			}
			if n.Post != nil {
				MutateNoRetCallStmts(owner_list, index, n.Post, bm)
			}
			bm(n.Body, MutateNoRetCallStmts)
		
		case *ast.IfStmt:
			if n.Init != nil {
				MutateNoRetCallStmts(owner_list, index, n.Init, bm)
			}
			bm(n.Body, MutateNoRetCallStmts)
			if n.Else != nil {
				MutateNoRetCallStmts(owner_list, index, n.Else, bm)
			}
			
		case *ast.SwitchStmt:
			MutateNoRetCallStmts(owner_list, index, n.Init, bm)
			bm(n.Body, MutateNoRetCallStmts)
		
		case *ast.CaseClause:
			for i, stmt := range n.Body {
				MutateNoRetCallStmts(&n.Body, i, stmt, bm)
			}
		
		case *ast.RangeStmt:
			bm(n.Body, MutateNoRetCallStmts)
		
		case *ast.ExprStmt:
			if fn, is_func_call := n.X.(*ast.CallExpr); is_func_call {
				fmt.Printf("ast.ExprStmt Mutator :: %+v\n", fn.Fun)
				if typ := ASTCtxt.TypeInfo.TypeOf(fn); typ != nil {
					switch t := typ.(type) {
						case *types.Tuple:
							extra_args := t.Len()
							if IsFuncPtr(fn) {
								if extra_args > 1 {
									retvals := make([]ast.Expr, 0)
									rettypes := make([]types.Type, 0)
									for i:=0; i<extra_args; i++ {
										ret_tmp := ast.NewIdent(fmt.Sprintf("fptr_temp%d", ASTCtxt.TmpVar))
										ASTCtxt.TmpVar++
										declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, nil, t.At(i).Type())
										*owner_list = InsertStmt(*owner_list, 0, declstmt)
										retvals = append(retvals, ret_tmp)
										rettypes = append(rettypes, t.At(i).Type())
									}
									
									calls := ExpandFuncPtrCalls(fn, retvals, rettypes)
									for i := len(calls)-1; i>0; i-- {
										*owner_list = InsertStmt(*owner_list, FindStmt(*owner_list, n)+1, calls[i])
									}
									n.X = calls[0].(*ast.ExprStmt).X
								} else {
									calls := ExpandFuncPtrCalls(fn, nil, nil)
									n.X = calls[0].(*ast.ExprStmt).X
									for i := len(calls)-1; i>0; i-- {
										*owner_list = InsertStmt(*owner_list, FindStmt(*owner_list, n)+1, calls[i])
									}
								}
							} else {
								if extra_args > 1 {
									for i:=1; i<extra_args; i++ {
										ret_tmp := ast.NewIdent(fmt.Sprintf("fn_temp%d", ASTCtxt.TmpVar))
										ASTCtxt.TmpVar++
										declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, nil, t.At(i).Type())
										*owner_list = InsertStmt(*owner_list, 0, declstmt)
										fn.Args = append(fn.Args, MakeReference(ret_tmp))
									}
								}
							}
						
						default:
							if IsFuncPtr(fn) {
								calls := ExpandFuncPtrCalls(fn, nil, nil)
								n.X = calls[0].(*ast.ExprStmt).X
								for i := len(calls)-1; i>0; i-- {
									*owner_list = InsertStmt(*owner_list, FindStmt(*owner_list, n)+1, calls[i])
								}
							}
					}
				}
			}
	}
}

/// a &^ b ==> a & ^(b) ==> a & ~(b) in sp.
func MutateAndNotExpr(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch x := n.(type) {
				case *ast.BinaryExpr:
					if x.Op==token.AND_NOT {
						x.Op = token.AND
						x.Y = MakeBitNotExpr(MakeParenExpr(x.Y))
					}
				case *ast.AssignStmt:
					if x.Tok==token.AND_NOT_ASSIGN {
						x.Tok = token.AND_ASSIGN
						for i := range x.Rhs {
							x.Rhs[i] = MakeBitNotExpr(MakeParenExpr(x.Rhs[i]))
						}
					}
			}
		}
		return true
	})
}


/// func(params){code}(args) => func _srcgo_func#(params){code} ... _srcgo_func#(args)
func MutateFuncLit(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator) {
	switch n := s.(type) {
		case *ast.BlockStmt:
			bm(n, MutateFuncLit)
		
		case *ast.ForStmt:
			if n.Init != nil {
				MutateFuncLit(owner_list, index, n.Init, bm)
			}
			if n.Cond != nil {
				MutateFuncLitExprs(&n.Cond, nil)
			}
			if n.Post != nil {
				MutateFuncLit(owner_list, index, n.Post, bm)
			}
			bm(n.Body, MutateFuncLit)
		
		case *ast.IfStmt:
			if n.Init != nil {
				MutateFuncLit(owner_list, index, n.Init, bm)
			}
			MutateFuncLitExprs(&n.Cond, nil)
			bm(n.Body, MutateFuncLit)
			if n.Else != nil {
				MutateFuncLit(owner_list, index, n.Else, bm)
			}
		
		case *ast.SwitchStmt:
			MutateFuncLit(owner_list, index, n.Init, bm)
			MutateFuncLitExprs(&n.Tag, nil)
			bm(n.Body, MutateFuncLit)
		
		case *ast.CaseClause:
			for j := range n.List {
				MutateFuncLitExprs(&n.List[j], nil)
			}
			for i, stmt := range n.Body {
				MutateFuncLit(&n.Body, i, stmt, bm)
			}
		
		case *ast.RangeStmt:
			MutateFuncLitExprs(&n.Key, nil)
			MutateFuncLitExprs(&n.Value, nil)
			MutateFuncLitExprs(&n.X, nil)
			bm(n.Body, MutateFuncLit)
		
		case *ast.ExprStmt:
			MutateFuncLitExprs(&n.X, nil)
		
		case *ast.ReturnStmt:
			for i := range n.Results {
				MutateFuncLitExprs(&n.Results[i], nil)
			}
		
		case *ast.AssignStmt:
			for i := range n.Rhs {
				MutateFuncLitExprs(&n.Rhs[i], nil)
			}
			for i := range n.Lhs {
				MutateFuncLitExprs(&n.Lhs[i], nil)
			}
		
		case *ast.DeclStmt:
			g := n.Decl.(*ast.GenDecl)
			for _, d := range g.Specs {
				switch g.Tok {
					case token.CONST, token.VAR:
						v := d.(*ast.ValueSpec)
						for expr := range v.Values {
							MutateFuncLitExprs(&v.Values[expr], nil)
						}
				}
			}
	}
}

func MutateFuncLitExprs(e, prev *ast.Expr) {
	if e==nil || *e == nil {
		return
	}
	switch n := (*e).(type) {
		case *ast.BinaryExpr:
			MutateFuncLitExprs(&n.X, e)
			MutateFuncLitExprs(&n.Y, e)
		
		case *ast.CallExpr:
			MutateFuncLitExprs(&n.Fun, e)
			for i := range n.Args {
				MutateFuncLitExprs(&n.Args[i], e)
			}
		
		case *ast.KeyValueExpr:
			MutateFuncLitExprs(&n.Key, e)
			MutateFuncLitExprs(&n.Value, e)
		
		case *ast.IndexExpr:
			MutateFuncLitExprs(&n.X, e)
			MutateFuncLitExprs(&n.Index, e)
		
		case *ast.UnaryExpr:
			MutateFuncLitExprs(&n.X, e)
		
		case *ast.FuncLit:
			tmp_func_name := ast.NewIdent(fmt.Sprintf("SrcGoTmpFunc%d", ASTCtxt.TmpFunc))
			ASTCtxt.TmpFunc++
			fn_decl := new(ast.FuncDecl)
			fn_decl.Name = tmp_func_name
			fn_decl.Type = n.Type
			n.Type = nil
			fn_decl.Body = n.Body
			n.Body = nil
			*e = tmp_func_name
			
			(*ASTCtxt.CurrFile).Decls = append((*ASTCtxt.CurrFile).Decls, fn_decl)
	}
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