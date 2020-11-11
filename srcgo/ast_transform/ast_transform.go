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
	"unicode"
	"errors"
	"go/token"
	"go/ast"
	"go/types"
	"go/format"
	"go/constant"
)


var ASTCtxt struct {
	TypeInfo      *types.Info
	CurrFunc      *ast.FuncDecl
	FSet          *token.FileSet
	BuiltInTypes  map[string]types.Object
	Err           func(err error)
	RangeIter     uint
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

func MakeIdent(n string) *ast.Ident {
	id := new(ast.Ident)
	id.Name = n
	return id
}

func Arrayify(typ ast.Expr, len ast.Expr) *ast.ArrayType {
	a := new(ast.ArrayType)
	a.Len = len
	a.Elt = typ
	return a
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
				x = MakeIdent(t.String())
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

func CheckReturnTypes(n ast.Node) bool {
	var list *ast.FieldList
	switch f := n.(type) {
		case *ast.FuncDecl:
			list = f.Type.Results
		case *ast.FuncType:
			list = f.Results
	}
	if list != nil {
		for i, ret := range list.List {
			if ptr, is_ptr := ret.Type.(*ast.StarExpr); is_ptr {
				PrintSrcGoErr(ptr.Pos(), "Returning Pointers isn't Allowed." + fmt.Sprintf(" Param %d is a pointer", i))
				return false
			}
		}
	}
	return true
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
	 * The typ argument may be a defined (Named) type or an alias type. It may also be nil such that the returned TypeName can be used as argument for NewNamed, which will set the TypeName's type as a side- effect.
	 * 
	 * 
	 * func NewNamed(obj *TypeName, underlying Type, methods []*Func) *Named
	 * 
	 * NewNamed returns a new named type for the given type name, underlying type, and associated methods. If the given type name obj doesn't have a type yet, its type is set to the returned named type. The underlying type must not be a *Named
	 */
	ASTCtxt.BuiltInTypes = make(map[string]types.Object)
	
	/// Basic types for SourcePawn
	MakeTypeAlias("char",    types.Typ[types.Int8], true)
	MakeTypeAlias("Entity",  types.Typ[types.Int], false)
	MakeTypeAlias("Address", types.Typ[types.Int], true)
	MakeTypeAlias("float",   types.Typ[types.Float32], false)
	
	/// Array types.
	vec3_array := types.NewArray(types.Typ[types.Float32], 3)
	vec3_type_name := types.NewTypeName(token.NoPos, nil, "Vec3", vec3_array)
	types.Universe.Insert(vec3_type_name)
	
	plugin_reg_struc := types.NewStruct([]*types.Var{
		types.NewField(token.NoPos, nil, "name",        types.Typ[types.String], false),
		types.NewField(token.NoPos, nil, "description", types.Typ[types.String], false),
		types.NewField(token.NoPos, nil, "author",      types.Typ[types.String], false),
		types.NewField(token.NoPos, nil, "version",     types.Typ[types.String], false),
		types.NewField(token.NoPos, nil, "url",         types.Typ[types.String], false),
	}, nil)
	plugin_reg_type_name := types.NewTypeName(token.NoPos, nil, "Plugin", nil)
	types.NewNamed(plugin_reg_type_name, plugin_reg_struc, nil)
	types.Universe.Insert(plugin_reg_type_name)
	
	/// per request of JoinedSenses.
	ext_struc := types.NewStruct([]*types.Var{
		types.NewField(token.NoPos, nil, "name",     types.Typ[types.String], false),
		types.NewField(token.NoPos, nil, "file",     types.Typ[types.String], false),
		types.NewField(token.NoPos, nil, "autoload", types.Typ[types.Int], false),
		types.NewField(token.NoPos, nil, "required", types.Typ[types.Int], false),
	}, nil)
	ext_type_name := types.NewTypeName(token.NoPos, nil, "Extension", nil)
	types.NewNamed(ext_type_name, ext_struc, nil)
	types.Universe.Insert(ext_type_name)
	
	/// Action
	MakeEnumType("Action", []string{"Plugin_Continue", "Plugin_Changed", "Plugin_Handled", "Plugin_Stop"}, []int64{0,1,2,3})
	
	/// Handle and derived types.
	MakeNamedType("Handle", types.Typ[types.UnsafePointer], nil)
	MakeNamedType("Map",    types.Typ[types.UnsafePointer], nil)
	MakeNamedType("Array",  types.Typ[types.UnsafePointer], nil)
	MakeNamedType("Event",  types.Typ[types.UnsafePointer], nil)
	
	/// defined constants.
	MakeIntVar("MaxClients")
	MakeIntConst("MAXPLAYERS", 65)
	MakeIntConst("MAXPLAYERS", 2048)
	
	/// Functions
	MakeFunc("IsClientInGame", nil, MakeParams([]string{"client"}, []types.Type{types.Typ[types.Int]}), MakeRet([]types.Type{types.Typ[types.Bool]}), false)
}


func AnalyzeFile(f *ast.File, info *types.Info, err_fn func(err error)) {
	ASTCtxt.TypeInfo = info
	ASTCtxt.Err = err_fn
	
	for _, decl := range f.Decls {
		ManageDeclNode(decl)
	}
}


/** Top Level of the Grammar
 * There's 4 types of nodes in Golang and their hierarchy:
 * 
 * Decl (Declaration) nodes.
 * Spec (Specification) nodes.
 * Stmt (Statement) nodes.
 * Expr (Expression) nodes.
 */

func ManageDeclNode(d ast.Decl) {
	switch decl := d.(type) {
		case *ast.GenDecl:
			AnalyzeGenDecl(decl)
		case *ast.FuncDecl:
			AnalyzeFuncDecl(decl)
	}
}

/// Generic Declaration Node
func AnalyzeGenDecl(g *ast.GenDecl) {
	for _, spec := range g.Specs {
		switch s := spec.(type) {
			case *ast.ImportSpec:
			case *ast.ValueSpec:
				for _, n := range s.Names {
					ManageExprNode(nil, nil, n)
				}
				ManageExprNode(nil, nil, s.Type)
				for _, e := range s.Values {
					ManageExprNode(nil, nil, e)
				}
			
			case *ast.TypeSpec:
				ManageExprNode(nil, nil, s.Name)
				/// make sure struct fields are not pointers or slices.
				switch t := s.Type.(type) {
					case *ast.StructType:
						if !t.Incomplete {
							for _, f := range t.Fields.List {
								switch t := f.Type.(type) {
									case *ast.StarExpr:
										PrintSrcGoErr(t.Pos(), "Pointers are not allowed in Structs.")
									case *ast.ArrayType:
										if t.Len==nil {
											PrintSrcGoErr(t.Pos(), "Arrays of unknown size are not allowed in Structs.")
										}
								}
							}
						}
					
					case *ast.FuncType:
						if t.Results != nil && len(t.Results.List) > 1 {
							/// too many return values, add them as pointer params!
							results := len(t.Results.List)
							for i:=1; i<results; i++ {
								ret := t.Results.List[i]
								/// if they're named, treat as reference types.
								if ret.Names != nil && len(ret.Names) > 1 {
									ret.Type = PtrizeExpr(ret.Type)
									t.Params.List = append(t.Params.List, ret)
								} else {
									ret.Names = append(ret.Names, ast.NewIdent(fmt.Sprintf("%s_param%d", s.Name.Name, i)))
									ret.Type = PtrizeExpr(ret.Type)
									t.Params.List = append(t.Params.List, ret)
								}
							}
							t.Results.List = t.Results.List[:1]
						} else if len(t.Results.List)==1 && t.Results.List[0].Names != nil && len(t.Results.List[0].Names) > 1 {
							t.Results.List[0].Type = PtrizeExpr(t.Results.List[0].Type)
							t.Params.List = append(t.Params.List, t.Results.List[0])
							t.Results = nil
						}
					
					/// TODO: make interface compile to a typeset?
					case *ast.InterfaceType:
				}
		}
	}
}


func AnalyzeFuncDecl(f *ast.FuncDecl) {
	ASTCtxt.CurrFunc = f
	new_params := make([]*ast.Field, 0)
	if f.Recv != nil {
		if len(f.Recv.List) > 1 {
			PrintSrcGoErr(f.Pos(), "Multiple Receiver Params are not allowed in Functions.")
		} else {
			/// merge receiver with the params and nullify it.
			new_params = append(new_params, f.Recv.List[0])
			if type_expr := ASTCtxt.TypeInfo.TypeOf(f.Recv.List[0].Type); type_expr != nil {
				type_name := type_expr.String()
				type_name = strings.Replace(type_name, ".", "_", -1)
				type_name = strings.Replace(type_name, " ", "_", -1)
				
				type_name = strings.TrimFunc(type_name, func(r rune) bool {
					return !unicode.IsLetter(r) && r != rune('_')
				})
				f.Name.Name = type_name + "_" + f.Name.Name
			}
			f.Recv = nil
		}
	}
	
	for _, param := range f.Type.Params.List {
		new_params = append(new_params, param)
	}
	
	if f.Type.Results != nil {
		if !CheckReturnTypes(f) {
			return
		}
		
		/// func f() (int, float) {} => func f(_param1 *float) int {}
		results := len(f.Type.Results.List)
		if results > 1 {
			for i:=1; i<results; i++ {
				ret := f.Type.Results.List[i]
				/// if they're named, treat as reference types.
				if ret.Names != nil && len(ret.Names) > 1 {
					ret.Type = PtrizeExpr(ret.Type)
					new_params = append(new_params, ret)
				} else {
					ret.Names = append(ret.Names, ast.NewIdent(fmt.Sprintf("%s_param%d", f.Name.Name, i)))
					ret.Type = PtrizeExpr(ret.Type)
					new_params = append(new_params, ret)
				}
			}
			f.Type.Results.List = f.Type.Results.List[:1]
		} else if results==1 && f.Type.Results.List[0].Names != nil && len(f.Type.Results.List[0].Names) > 1 {
			f.Type.Results.List[0].Type = PtrizeExpr(f.Type.Results.List[0].Type)
			new_params = append(new_params, f.Type.Results.List[0])
			f.Type.Results.List = nil
		}
	}
	f.Type.Params.List = new_params
	
	if f.Body != nil {
		AnalyzeBlockStmt(f.Body);
	}
	ASTCtxt.CurrFunc = nil
	ASTCtxt.RangeIter = 0
}

func ManageStmtNode(owner_list *[]ast.Stmt, index int, s ast.Stmt) {
	switch n := s.(type) {
		case *ast.AssignStmt:
			left_len, rite_len := len(n.Lhs), len(n.Rhs)
			funct, is_func_call := n.Rhs[0].(*ast.CallExpr)
			if rite_len==1 && left_len >= rite_len && is_func_call {
				/// a func call returning multiple items.
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
								PrintSrcGoErr(n.TokPos, "Failed to expand assignments.")
							}
						}
						
						for key, val := range var_map {
							val_spec := new(ast.ValueSpec)
							for _, name := range val {
								val_spec.Names = append(val_spec.Names, name.(*ast.Ident))
							}
							type_expr := TypeToASTExpr(key)
							val_spec.Type = type_expr
							gen_decl.Specs = append(gen_decl.Specs, val_spec)
						}
						
						decl_stmt.Decl = gen_decl
						AnalyzeGenDecl(gen_decl)
						
						*owner_list = InsertStmt(*owner_list, index, decl_stmt)
						n.Tok = token.ASSIGN
						ManageStmtNode(owner_list, index, s)
					
					case token.ASSIGN: /// transform the tuple return into a single return + pass by ref.
						if is_func_call {
							for i:=1; i<left_len; i++ {
								switch e := n.Lhs[i].(type) {
									case *ast.Ident:
										funct.Args = append(funct.Args, MakeReference(e))
								}
							}
							n.Lhs = n.Lhs[:1]
						}
				}
			}
			for _, e := range n.Lhs {
				ManageExprNode(owner_list, s, e)
			}
			for _, e := range n.Rhs {
				ManageExprNode(owner_list, s, e)
			}
		
		case *ast.BlockStmt:
			AnalyzeBlockStmt(n)
		
		case *ast.BranchStmt:
			if n.Tok==token.GOTO || n.Tok==token.FALLTHROUGH {
				PrintSrcGoErr(n.Pos(), fmt.Sprintf(" %s is Illegal.", n.Tok.String()))
			} else if n.Label != nil {
				PrintSrcGoErr(n.Pos(), "Branched Labels are Illegal.")
			}
		
		case *ast.DeclStmt:
			ManageDeclNode(n.Decl)
			
		case *ast.EmptyStmt:
			
		case *ast.ExprStmt:
			ManageExprNode(owner_list, s, n.X)
			
		case *ast.ForStmt:
			/// in Golang, 'for' replaces both for and while-loops.
			/// we'll have to replace while-loop like constructs with a degenerate for-loop
			if n.Init != nil { /// initialization statement; or nil
				ManageStmtNode(owner_list, index, n.Init)
			}
			if n.Cond != nil { /// condition; or nil
				ManageExprNode(owner_list, s, n.Cond)
			}
			if n.Post != nil { /// post iteration statement; or nil
				ManageStmtNode(owner_list, index, n.Post)
			}
			AnalyzeBlockStmt(n.Body)
			
		case *ast.IfStmt:
			/// TODO: have initializer stmt prior to the if stmt.
			if n.Init != nil {
				ManageStmtNode(owner_list, index, n.Init)
			}
			ManageExprNode(owner_list, s, n.Cond)
			AnalyzeBlockStmt(n.Body)
			if n.Else != nil {
				ManageStmtNode(owner_list, index, n.Else)
			}
			
		case *ast.IncDecStmt:
			ManageExprNode(owner_list, s, n.X)
			
		case *ast.ReturnStmt:
			/// change multiple var returns into passing by reference.
			index := FindStmt(*owner_list, s)
			res_len := len(n.Results)
			for i:=1; i<res_len; i++ {
				ptr_deref := PtrizeExpr(ast.NewIdent(fmt.Sprintf("%s_param%d", ASTCtxt.CurrFunc.Name.Name, i)))
				assign := new(ast.AssignStmt)
				assign.Lhs = append(assign.Lhs, ptr_deref)
				assign.Tok = token.ASSIGN
				assign.Rhs = append(assign.Rhs, n.Results[i])
				*owner_list = InsertStmt(*owner_list, index, assign)
				index++
			}
			if res_len > 1 {
				n.Results = n.Results[:1]
				ManageStmtNode(owner_list, index, s) /// reanalyze
			}
		
		case *ast.SwitchStmt:
			ManageStmtNode(owner_list, index, n.Init)
			ManageExprNode(owner_list, s, n.Tag)
			AnalyzeBlockStmt(n.Body)
		
		case *ast.CaseClause:
			for _, expr := range n.List {
				ManageExprNode(owner_list, s, expr)
			}
			for i, stmt := range n.Body {
				ManageStmtNode(&n.Body, i, stmt)
			}
		
		case *ast.RangeStmt: /// TODO: adapt range statement for ArrayLists and other containers.
			if n.Key != nil {
				if iden, ok := n.Key.(*ast.Ident); ok && iden.Name=="_" {
					n.Key = MakeIdent(fmt.Sprintf("%s_Iter%d", ASTCtxt.CurrFunc.Name.Name, ASTCtxt.RangeIter))
					ASTCtxt.RangeIter++
				}
			} else {
				n.Key = MakeIdent(fmt.Sprintf("%s_Iter%d", ASTCtxt.CurrFunc.Name.Name, ASTCtxt.RangeIter))
				ASTCtxt.RangeIter++
			}
			ManageExprNode(owner_list, s, n.Key)
			
			if n.Value != nil {
				if iden, ok := n.Value.(*ast.Ident); ok && iden.Name=="_" {
					n.Value = nil
				} else {
					switch n.Tok {
						case token.DEFINE:
							decl_stmt := new(ast.DeclStmt)
							gen_decl := new(ast.GenDecl)
							gen_decl.Tok = token.VAR
							
							val_spec := new(ast.ValueSpec)
							id := n.Value.(*ast.Ident)
							val_spec.Names = append(val_spec.Names, id)
							
							val_spec.Type = ValueToTypeExpr(n.Value)
							gen_decl.Specs = append(gen_decl.Specs, val_spec)
							
							decl_stmt.Decl = gen_decl
							n.Body.List = InsertStmt(n.Body.List, 0, decl_stmt)
							
							assign := new(ast.AssignStmt)
							assign.Lhs = append(assign.Lhs, n.Value)
							
							get_index := new(ast.IndexExpr)
							get_index.Index = n.Key
							get_index.X = n.X
							assign.Rhs = append(assign.Rhs, get_index)
							assign.Tok = token.ASSIGN
							n.Body.List = InsertStmt(n.Body.List, 1, assign)
							n.Value = nil
					}
				}
			}
			ManageExprNode(owner_list, s, n.X)
			AnalyzeBlockStmt(n.Body)
		
		case *ast.CommClause:
			PrintSrcGoErr(n.Pos(), "Comm Select Cases are Illegal.")
		case *ast.DeferStmt:
			PrintSrcGoErr(n.Pos(), "Defer Statements are Illegal.")
		case *ast.TypeSwitchStmt:
			PrintSrcGoErr(n.Pos(), "Type-Switches are Illegal.")
		case *ast.LabeledStmt:
			PrintSrcGoErr(n.Pos(), "Labels are Illegal.")
		case *ast.GoStmt:
			PrintSrcGoErr(n.Pos(), "Goroutines are Illegal.")
		case *ast.SelectStmt:
			PrintSrcGoErr(n.Pos(), "Select Statements are Illegal.")
		case *ast.SendStmt:
			PrintSrcGoErr(n.Pos(), "Send Statements are Illegal.")
	}
}

func AnalyzeBlockStmt(b *ast.BlockStmt) {
	for i, stmt := range b.List {
		ManageStmtNode(&b.List, i, stmt)
	}
}


func ManageExprNode(owner_list *[]ast.Stmt, owner_stmt ast.Stmt, e ast.Expr) {
	switch x := e.(type) {
		case *ast.IndexExpr:
			ManageExprNode(owner_list, owner_stmt, x.X)
			ManageExprNode(owner_list, owner_stmt, x.Index)
		
		case *ast.KeyValueExpr:
			ManageExprNode(owner_list, owner_stmt, x.Key)
			ManageExprNode(owner_list, owner_stmt, x.Value)
		
		case *ast.ParenExpr:
			ManageExprNode(owner_list, owner_stmt, x.X)
		
		case *ast.StarExpr:
			/// in an ordinary block, we ignore the dereference since it'll become a reference.
			ManageExprNode(owner_list, owner_stmt, x.X)
		
		case *ast.UnaryExpr:
			ManageExprNode(owner_list, owner_stmt, x.X)
		
		case *ast.CallExpr:
			is_fptr := IsFuncPtr(x.Fun)
			/*
			if is_fptr {
				Call_StartFunction := ast.NewIdent("Call_StartFunction")
				new_args := []ast.Expr{ast.NewIdent("nil"), x.Fun}
				
				//call_start_stmt := new(ast.ExprStmt)
				//call_start_stmt.X = Call_StartFunction
				stmt_index := FindStmt(*owner_list, owner_stmt)
				x.Fun = Call_StartFunction
				//*owner_list = InsertStmt(*owner_list, stmt_index, call_start_stmt)
				
				for _, arg := range x.Args {
					if typ := ASTCtxt.TypeInfo.TypeOf(arg); typ != nil {
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
										*owner_list = InsertStmt(*owner_list, stmt_index+1, Call_PushString_stmt)
									default:
										Call_PushArrayEx := new(ast.CallExpr)
										Call_PushArrayEx.Fun = ast.NewIdent("Call_PushArrayEx")
										Call_PushArrayEx.Args = append(Call_PushArrayEx.Args, arg)
										
										Call_PushArrayEx.Args = append(Call_PushArrayEx.Args, MakeBasicLit(token.INT, fmt.Sprintf("%d", t.Len())))
										
										Call_PushArrayEx.Args = append(Call_PushArrayEx.Args, MakeBasicLit(token.INT, "1"))
										
										Call_PushArrayEx_stmt := new(ast.ExprStmt)
										Call_PushArrayEx_stmt.X = Call_PushArrayEx
										*owner_list = InsertStmt(*owner_list, stmt_index+1, Call_PushArrayEx_stmt)
								}
							case *types.Pointer:
								switch t.Elem().String() {
									case "float":
										Call_PushFloatRef := new(ast.CallExpr)
										Call_PushFloatRef.Fun = ast.NewIdent("Call_PushFloatRef")
										Call_PushFloatRef.Args = append(Call_PushFloatRef.Args, arg)
										
										Call_PushFloatRef_stmt := new(ast.ExprStmt)
										Call_PushFloatRef_stmt.X = Call_PushFloatRef
										*owner_list = InsertStmt(*owner_list, stmt_index+1, Call_PushFloatRef_stmt)
									default:
										Call_PushCellRef := new(ast.CallExpr)
										Call_PushCellRef.Fun = ast.NewIdent("Call_PushCellRef")
										Call_PushCellRef.Args = append(Call_PushCellRef.Args, arg)
										
										Call_PushCellRef_stmt := new(ast.ExprStmt)
										Call_PushCellRef_stmt.X = Call_PushCellRef
										*owner_list = InsertStmt(*owner_list, stmt_index+1, Call_PushCellRef_stmt)
								}
							case *types.Basic:
								switch t.Name() {
									case "string":
										Call_PushString := new(ast.CallExpr)
										Call_PushString.Fun = ast.NewIdent("Call_PushString")
										Call_PushString.Args = append(Call_PushString.Args, arg)
										
										Call_PushString_stmt := new(ast.ExprStmt)
										Call_PushString_stmt.X = Call_PushString
										*owner_list = InsertStmt(*owner_list, stmt_index+1, Call_PushString_stmt)
									case "float", "float32":
										Call_PushFloat := new(ast.CallExpr)
										Call_PushFloat.Fun = ast.NewIdent("Call_PushFloat")
										Call_PushFloat.Args = append(Call_PushFloat.Args, arg)
										
										Call_PushFloat_stmt := new(ast.ExprStmt)
										Call_PushFloat_stmt.X = Call_PushFloat
										*owner_list = InsertStmt(*owner_list, stmt_index+1, Call_PushFloat_stmt)
									default:
										//"Call_PushCell(%s); ", GetExprString(arg))
										Call_PushCell := new(ast.CallExpr)
										Call_PushCell.Fun = ast.NewIdent("Call_PushCell")
										Call_PushCell.Args = append(Call_PushCell.Args, arg)
										
										Call_PushCell_stmt := new(ast.ExprStmt)
										Call_PushCell_stmt.X = Call_PushCell
										*owner_list = InsertStmt(*owner_list, stmt_index+1, Call_PushCell_stmt)
								}
						}
					}
				}
				//"Call_Finish();")
				x.Args = new_args
			}*/
			if caller, is_method_call := x.Fun.(*ast.SelectorExpr); is_method_call && !is_fptr {
				/// check for *ast.SelectorExpr
				x.Args = InsertExpr(x.Args, 0, caller.X)
				if typ := ASTCtxt.TypeInfo.TypeOf(caller.X); typ != nil {
					caller.Sel.Name = typ.String() + "_" + caller.Sel.Name
				}
				x.Fun = caller.Sel
			}
			
			ManageExprNode(owner_list, owner_stmt, x.Fun)
			for _, arg := range x.Args {
				ManageExprNode(owner_list, owner_stmt, arg)
			}
		
		case *ast.BinaryExpr: /// a op b
			/// a &^ b ==> a & ^(b) ==> a & ~(b) in sp.
			if x.Op==token.AND_NOT {
				x.Op = token.AND
				
				n := new(ast.UnaryExpr)
				n.Op = token.XOR
				
				p := new(ast.ParenExpr)
				p.X = x.Y
				n.X = p
				x.Y = n
			} else if x.Op==token.AND_NOT_ASSIGN {
				x.Op = token.AND_ASSIGN
				
				n := new(ast.UnaryExpr)
				n.Op = token.XOR
				
				p := new(ast.ParenExpr)
				p.X = x.Y
				n.X = p
				x.Y = n
			}
			ManageExprNode(owner_list, owner_stmt, x.X)
			ManageExprNode(owner_list, owner_stmt, x.Y)
		
		case *ast.SelectorExpr:
			ManageExprNode(owner_list, owner_stmt, x.X)
			ManageExprNode(owner_list, owner_stmt, x.Sel)
		
		case *ast.Ident:
			//Obj *Object   // denoted object; or nil
		
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

func PrintNode(n ast.Node) string {
	var ast_str string
	ast.Inspect(n, func(n ast.Node) bool {
		if n != nil {
			ast_str += fmt.Sprintf("%p - %T:\t\t%+v", n, n, n) + "\n"
		}
		return true
	})
	return ast_str
}

func PrintAST(f *ast.File) string {
	var ast_str string
	ast.Inspect(f, func(n ast.Node) bool {
		if n != nil {
			ast_str += fmt.Sprintf("%p - %T:\t\t%+v", n, n, n) + "\n"
		}
		return true
	})
	return ast_str
}

func PrettyPrintAST(f ast.Node) string {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	err := format.Node(&buf, fset, f)
	if err != nil {
		fmt.Println(err)
	}
	return buf.String()
}