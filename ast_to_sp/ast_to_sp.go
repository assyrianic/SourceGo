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

package AST2SP


import (
	"fmt"
	"bytes"
	"go/token"
	"go/ast"
	"go/types"   /// Golang's Type system.
	"go/format"
)

var type_info *types.Info

func PtrizeExpr(x ast.Expr) *ast.StarExpr {
	ptr := new(ast.StarExpr)
	ptr.X = x
	return ptr
}

func Arrayify(typ ast.Expr) *ast.ArrayType {
	a := new(ast.ArrayType)
	a.Len = nil
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


func AnalyzeFile(f *ast.File, info *types.Info) {
	/*
	type File struct {
		Doc        *CommentGroup   // associated documentation; or nil
		Package    token.Pos       // position of "package" keyword
		Name       *Ident          // package name
		Decls      []Decl          // top-level declarations; or nil
		Scope      *Scope          // package scope (this file only)
		Imports    []*ImportSpec   // imports in this file
		Unresolved []*Ident        // unresolved identifiers in this file
		Comments   []*CommentGroup // list of all comments in the source file
	}
	 */
	type_info = info
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
		switch spec.(type) {
			case *ast.ImportSpec:
			case *ast.ValueSpec:
			/* ConstSpec or VarSpec production
				Doc     *CommentGroup // associated documentation; or nil
				Names   []*Ident      // value names (len(Names) > 0)
				Type    Expr          // value type; or nil
				Values  []Expr        // initial values; or nil
				Comment *CommentGroup // line comments; or nil
			 */
			case *ast.TypeSpec:
			/*
				Doc     *CommentGroup // associated documentation; or nil
				Name    *Ident        // type name
				Assign  token.Pos     // position of '=', if any; added in Go 1.9
				Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
				Comment *CommentGroup // line comments; or nil
			 */
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
		
		/// func f() (int, float) {} => func f(_1 *float) int {}
		results := len(f.Type.Results.List)
		if results > 1 {
			for i:=results-1; i>=0; i-- {
				ret := f.Type.Results.List[i]
				/// if they're named, treat as reference types.
				if ret.Names != nil && len(ret.Names) > 1 {
					ret.Type = PtrizeExpr(ret.Type)
					new_params = append(new_params, ret)
					copy(f.Type.Results.List[i:], f.Type.Results.List[i+1:])
					f.Type.Results.List = f.Type.Results.List[:len(f.Type.Results.List)-1]
				} else if i != 0 {
					ret.Names = append(ret.Names, ast.NewIdent(fmt.Sprintf("_%d", results - i)))
					ret.Type = PtrizeExpr(ret.Type)
					new_params = append(new_params, ret)
					copy(f.Type.Results.List[i:], f.Type.Results.List[i+1:])
					f.Type.Results.List = f.Type.Results.List[:len(f.Type.Results.List)-1]
				}
			}
		} else if results==1 && f.Type.Results.List[0].Names != nil && len(f.Type.Results.List[0].Names) > 1 {
			/// TODO: analyze if the return values are all the same type and group it to return an array EXCEPT if the function is a native...
			//arr := Arrayify(f.Type.Results.List[0].Type)
			//f.Type.Results.List[0].Type = arr
			f.Type.Results.List[0].Type = PtrizeExpr(f.Type.Results.List[0].Type)
			new_params = append(new_params, f.Type.Results.List[0])
			f.Type.Results.List = nil
			//f.Type.Results.List[0].Names = nil
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
			same_len := left_len==rite_len
			if !same_len && left_len > rite_len && rite_len==1 {
				/// probably a func call returning multiple items.
				switch n.Tok {
					case token.DEFINE: /// TODO: break this down into decl + assigns
						decl_stmt := new(ast.DeclStmt)
						gen_decl := new(ast.GenDecl)
						gen_decl.Tok = token.VAR
						
						val_spec := new(ast.ValueSpec)
						val_spec.Names = make([]*ast.Ident, 0)
						for _, e := range n.Lhs {
							n := e.(*ast.Ident)
							val_spec.Names = append(val_spec.Names, n)
							if type_of_expr := type_info.TypeOf(e); type_of_expr != nil {
								type_expr := new(ast.Ident)
								type_expr.Name = type_of_expr.String()
								val_spec.Type = type_expr
							} else {
								panic("SourceGo: failed to space out assignments.")
							}
						}
						//Values  []Expr        // initial values; or nil
						
						gen_decl.Specs = []ast.Spec{val_spec}
						decl_stmt.Decl = gen_decl
						owner_block.List = InsertStmt(owner_block.List, FindStmt(owner_block.List, s), decl_stmt)
						n.Tok = token.ASSIGN
						AnalyzeBlockStmt(owner_block)
					case token.ASSIGN: /// transform the tuple return into a single return + pass by ref.
						for i:=1; i<left_len; i++ {
							
						}
				}
			}
			for _, e := range n.Lhs {
				ManageExprNode(owner_block, e)
			}
			for _, e := range n.Rhs {
				ManageExprNode(owner_block, e)
			}
		
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
			/// assumes tabs have been written to string builder.
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
			for _, result := range n.Results {
				/// change multiple var returns into passing by reference.
				ManageExprNode(owner_block, result)
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
		case *ast.RangeStmt:
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
			/// change map access to GetTrieValue calls.
			/// func GetTrieValue(map Handle, key string, value *any) bool
			/// value, exists := map[key] => if( GetTrieValue(map, key, value) )...
			/// value = map[key] => GetTrieValue(map, key, value);
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
			ManageExprNode(owner_block, x.Fun)
			for _, arg := range x.Args {
				ManageExprNode(owner_block, arg)
			}
		
		case *ast.BinaryExpr:
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