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

package SrcGo_ASTMod


import (
	"fmt"
	"bytes"
	"strings"
	"unicode"
	"go/token"
	"go/ast"
	"go/types"
	"go/format"
	"go/constant"
)


var (
	SrcGoTypeInfo *types.Info
)


func PtrizeExpr(x ast.Expr) *ast.StarExpr {
	ptr := new(ast.StarExpr)
	ptr.X = x
	return ptr
}

func Arrayify(typ ast.Expr, len ast.Expr) *ast.ArrayType {
	a := new(ast.ArrayType)
	a.Len = len
	a.Elt = typ
	return a
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
	types.Universe.Insert(types.NewTypeName(token.NoPos, nil, "Action", types.Typ[types.Int]))
	types.Universe.Insert(types.NewTypeName(token.NoPos, nil, "char", types.Typ[types.Int8]))
	types.Universe.Insert(types.NewTypeName(token.NoPos, nil, "Entity", types.Typ[types.Int]))
	types.Universe.Insert(types.NewTypeName(token.NoPos, nil, "float", types.Typ[types.Float32]))
	
	vec3_array := types.NewArray(types.Typ[types.Float32], 3)
	vec3_type_name := types.NewTypeName(token.NoPos, nil, "Vec3", vec3_array)
	types.Universe.Insert(vec3_type_name)
	
	/// TODO: define methods for the Handle types, Vec3, and Entity.
	/// also TODO: Add QAngle, AngularImpulse as [3]float like Vec3
	handle_type_name := types.NewTypeName(token.NoPos, nil, "Handle", nil)
	handle_type := types.NewNamed(handle_type_name, types.Typ[types.UnsafePointer], nil)
	types.Universe.Insert(handle_type_name)
	
	strmap := types.NewTypeName(token.NoPos, nil, "Map", handle_type)
	types.NewNamed(strmap, nil, nil)
	types.Universe.Insert(strmap)
	
	/// defined constants.
	types.Universe.Insert(types.NewVar(token.NoPos, nil, "MaxClients", types.Typ[types.Int]))
	
	types.Universe.Insert(types.NewConst(token.NoPos, nil, "MAXPLAYERS", types.Typ[types.Int], constant.MakeInt64(65)))
	
	types.Universe.Insert(types.NewConst(token.NoPos, nil, "MAXENTS", types.Typ[types.Int], constant.MakeInt64(2048)))
}


func AnalyzeFile(f *ast.File, info *types.Info) {
	SrcGoTypeInfo = info
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
					ManageExprNode(nil, n)
				}
				ManageExprNode(nil, s.Type)
				for _, e := range s.Values {
					ManageExprNode(nil, e)
				}
			
			case *ast.TypeSpec:
				ManageExprNode(nil, s.Name)
				/// make sure struct fields are not pointers or slices.
				if struc, is_struct := s.Type.(*ast.StructType); is_struct {
					if !struc.Incomplete {
						for _, f := range struc.Fields.List {
							switch t := f.Type.(type) {
								case *ast.StarExpr:
									panic("SourceGo: you can't have pointers in structs.")
								case *ast.ArrayType:
									if t.Len==nil {
										panic("SourceGo: you can't have slices in structs.")
									}
							}
						}
					}
				}
		}
	}
}


func AnalyzeFuncDecl(f *ast.FuncDecl) {
	new_params := make([]*ast.Field, 0)
	if f.Recv != nil {
		if len(f.Recv.List) > 1 {
			panic("SourceGo: You can't have multiple receivers.")
		} else {
			/// merge receiver with the params and nullify it.
			new_params = append(new_params, f.Recv.List[0])
			if type_expr := SrcGoTypeInfo.TypeOf(f.Recv.List[0].Type); type_expr != nil {
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
		for _, ret := range f.Type.Results.List {
			if _, is_ptr := ret.Type.(*ast.StarExpr); is_ptr {
				panic("SourceGo: You can't return pointers.")
			}
		}
		
		/// func f() (int, float) {} => func f(_param1 *float) int {}
		results := len(f.Type.Results.List)
		if results > 1 {
			/// TODO: group together moved var names under their same type if possible.
			for i:=1; i<results; i++ {
				ret := f.Type.Results.List[i]
				/// if they're named, treat as reference types.
				if ret.Names != nil && len(ret.Names) > 1 {
					ret.Type = PtrizeExpr(ret.Type)
					new_params = append(new_params, ret)
				} else {
					ret.Names = append(ret.Names, ast.NewIdent(fmt.Sprintf("srcgo_param%d", i)))
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
}

func ManageStmtNode(owner_block *ast.BlockStmt, s ast.Stmt) {
	switch n := s.(type) {
		case *ast.AssignStmt:
			/// TODO: make sure to check if len(rhs) <= len(lhs).
			/// also check if rhs is function call expr.
			left_len := len(n.Lhs)
			rite_len := len(n.Rhs)
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
							if type_expr := SrcGoTypeInfo.TypeOf(e); type_expr != nil {
								var_map[type_expr] = append(var_map[type_expr], e)
							} else {
								panic("SourceGo: failed to space out assignments.")
							}
						}
						
						for key, val := range var_map {
							val_spec := new(ast.ValueSpec)
							for _, name := range val {
								val_spec.Names = append(val_spec.Names, name.(*ast.Ident))
							}
							type_expr := new(ast.Ident)
							type_expr.Name = key.String()
							val_spec.Type = type_expr
							gen_decl.Specs = append(gen_decl.Specs, val_spec)
						}
						
						decl_stmt.Decl = gen_decl
						owner_block.List = InsertStmt(owner_block.List, FindStmt(owner_block.List, s), decl_stmt)
						n.Tok = token.ASSIGN
						AnalyzeBlockStmt(owner_block)
					
					case token.ASSIGN: /// transform the tuple return into a single return + pass by ref.
						if is_func_call {
							for i:=1; i<left_len; i++ {
								switch e := n.Lhs[i].(type) {
									case *ast.Ident:
										uexpr := new(ast.UnaryExpr)
										uexpr.X = e
										uexpr.Op = token.AND
										funct.Args = append(funct.Args, uexpr)
								}
							}
							n.Lhs = n.Lhs[:1]
						}
				}
			}
			for _, e := range n.Lhs {
				ManageExprNode(owner_block, e)
			}
			for _, e := range n.Rhs {
				ManageExprNode(owner_block, e)
			}
			//Values  []Expr        // initial values; or nil
		
		case *ast.BlockStmt:
			AnalyzeBlockStmt(n)
		
		case *ast.BranchStmt:
			if n.Tok==token.GOTO || n.Tok==token.FALLTHROUGH {
				panic("SourceGo: " + fmt.Sprintf("%s is illegal.", n.Tok.String()))
			} else if n.Label != nil {
				panic("SourceGo: Branched Labels are illegal.")
			}
		
		case *ast.DeclStmt:
			ManageDeclNode(n.Decl)
			
		case *ast.EmptyStmt:
			
		case *ast.ExprStmt:
			ManageExprNode(owner_block, n.X)
			
		case *ast.ForStmt:
			/// in Golang, 'for' replaces both for and while-loops.
			/// we'll have to replace while-loop like constructs with a degenerate for-loop
			if n.Init != nil { /// initialization statement; or nil
				ManageStmtNode(owner_block, n.Init)
			}
			if n.Cond != nil { /// condition; or nil
				ManageExprNode(owner_block, n.Cond)
			}
			if n.Post != nil { /// post iteration statement; or nil
				ManageStmtNode(owner_block, n.Post)
			}
			AnalyzeBlockStmt(n.Body)
			
		case *ast.IfStmt:
			/// TODO: have initializer stmt prior to the if stmt.
			if n.Init != nil {
				ManageStmtNode(owner_block, n.Init)
			}
			ManageExprNode(owner_block, n.Cond)
			AnalyzeBlockStmt(n.Body)
			if n.Else != nil {
				ManageStmtNode(owner_block, n.Else)
			}
			
		case *ast.IncDecStmt:
			ManageExprNode(owner_block, n.X)
			
		case *ast.ReturnStmt:
			/// change multiple var returns into passing by reference.
			index := FindStmt(owner_block.List, s)
			res_len := len(n.Results)
			for i:=1; i<res_len; i++ {
				ptr_deref := PtrizeExpr(ast.NewIdent(fmt.Sprintf("srcgo_param%d", i)))
				assign := new(ast.AssignStmt)
				assign.Lhs = append(assign.Lhs, ptr_deref)
				assign.Tok = token.ASSIGN
				assign.Rhs = append(assign.Rhs, n.Results[i])
				owner_block.List = InsertStmt(owner_block.List, index, assign)
				index++
			}
			n.Results = n.Results[:1]
			if res_len > 1 {
				AnalyzeBlockStmt(owner_block) /// reanalyze
			}
		
		case *ast.SwitchStmt:
			ManageStmtNode(owner_block, n.Init)
			ManageExprNode(owner_block, n.Tag)
			AnalyzeBlockStmt(n.Body)
			
			/** TODO:
			 * Switch statements can be "true" aka empty expression
			 * to work as a more compact if-else-if series:
			 * 
			 * switch {
			 *     case i < 10:
			 *         code()
			 *     case i > 10:
			 *         code()
			 * }
			 * 
			 * See if we can transform a true-switch into if-else-if for SourcePawn.
			 */
		
		case *ast.CaseClause:
			for _, expr := range n.List {
				ManageExprNode(owner_block, expr)
			}
			for _, stmt := range n.Body {
				ManageStmtNode(owner_block, stmt)
			}
		
		case *ast.CommClause:
			panic("SourceGo: Comm Select Cases are illegal.")
		case *ast.RangeStmt: /// TODO: allow ranges for fixed-sized arrays?
			panic("SourceGo: Ranges are illegal.")
		case *ast.DeferStmt:
			panic("SourceGo: Defer Statements are illegal.")
		case *ast.TypeSwitchStmt:
			panic("SourceGo: Type Switches are illegal.")
		case *ast.LabeledStmt:
			panic("SourceGo: Labels are illegal.")
		case *ast.GoStmt:
			panic("SourceGo: Goroutines are illegal.")
		case *ast.SelectStmt:
			panic("SourceGo: Select is illegal.")
		case *ast.SendStmt:
			panic("SourceGo: Send is illegal.")
	}
}

func AnalyzeBlockStmt(b *ast.BlockStmt) {
	for _, stmt := range b.List {
		ManageStmtNode(b, stmt)
	}
}


func ManageExprNode(owner_block *ast.BlockStmt, e ast.Expr) {
	switch x := e.(type) {
		case *ast.IndexExpr:
			ManageExprNode(owner_block, x.X)
			ManageExprNode(owner_block, x.Index)
		
		case *ast.KeyValueExpr:
			ManageExprNode(owner_block, x.Key)
			ManageExprNode(owner_block, x.Value)
		
		case *ast.ParenExpr:
			ManageExprNode(owner_block, x.X)
		
		case *ast.StarExpr:
			/// in an ordinary block, we ignore the dereference since it'll become a reference.
			ManageExprNode(owner_block, x.X)
		
		case *ast.UnaryExpr:
			ManageExprNode(owner_block, x.X)
		
		case *ast.CallExpr:
			/// check for *ast.SelectorExpr
			if caller, is_method_call := x.Fun.(*ast.SelectorExpr); is_method_call {
				x.Args = InsertExpr(x.Args, 0, caller.X)
				x.Fun = caller.Sel
			}
			ManageExprNode(owner_block, x.Fun)
			for _, arg := range x.Args {
				ManageExprNode(owner_block, arg)
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
			ManageExprNode(owner_block, x.X)
			ManageExprNode(owner_block, x.Y)
		
		case *ast.SelectorExpr:
			ManageExprNode(owner_block, x.X)
			ManageExprNode(owner_block, x.Sel)
		
		case *ast.Ident:
			//Obj *Object   // denoted object; or nil
		
		case *ast.BasicLit:
			if x.Kind==token.IMAG {
				panic("SourceGo: Imaginary numbers are illegal.")
			}
		
		case *ast.TypeAssertExpr:
			panic("SourceGo: Type Assertions are illegal.")
		case *ast.SliceExpr:
			/// TODO: make new, local array when slicing?
			panic("SourceGo: Slice Expressions are illegal.")
	}
}

func PrintAST(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		if n != nil {
			fmt.Println(fmt.Sprintf("%p - %T:\t\t", n, n), n)
		}
		return true
	})
	fmt.Println("\n")
}

func PrettyPrintAST(f *ast.File) {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	err := format.Node(&buf, fset, f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(buf.String())
}
