/* 
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
	"go/types"
	"go/ast"
	"go/token"
)

/* Case Studies of Changing assignment definitions:
 *
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
 * 
 */
func (a *AstTransmitter) MutateAssignDecls(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Body != nil {
				MutateBlock(d.Body, MutateAssignDefStmts, a)
			}
		}
	}
}

func MutateAssignDefStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator, ctxt any) {
	switch n := s.(type) {
	case *ast.BlockStmt:
		bm(n, MutateAssignDefStmts, ctxt)
	
	case *ast.ForStmt:
		if n.Init != nil {
			MutateAssignDefStmts(owner_list, index, n.Init, bm, ctxt)
		}
		if n.Post != nil {
			MutateAssignDefStmts(owner_list, index, n.Post, bm, ctxt)
		}
		bm(n.Body, MutateAssignDefStmts, ctxt)
	
	case *ast.IfStmt:
		if n.Init != nil {
			MutateAssignDefStmts(owner_list, index, n.Init, bm, ctxt)
		}
		bm(n.Body, MutateAssignDefStmts, ctxt)
		if n.Else != nil {
			MutateAssignDefStmts(owner_list, index, n.Else, bm, ctxt)
		}
	
	case *ast.SwitchStmt:
		MutateAssignDefStmts(owner_list, index, n.Init, bm, ctxt)
		bm(n.Body, MutateAssignDefStmts, ctxt)
	
	case *ast.CaseClause:
		for i, stmt := range n.Body {
			MutateAssignDefStmts(&n.Body, i, stmt, bm, ctxt)
		}
	
	case *ast.RangeStmt:
		bm(n.Body, MutateAssignDefStmts, ctxt)
	
	case *ast.AssignStmt:
		context := ctxt.(*AstTransmitter)
		left_len, rite_len := len(n.Lhs), len(n.Rhs)
		if rite_len==1 && left_len >= rite_len {
			switch e := n.Rhs[0].(type) {
			case *ast.CallExpr:
				if iden, is_ident := e.Fun.(*ast.Ident); is_ident && iden.Name=="make" {
					arg_len := len(e.Args)
					switch {
					case arg_len > 2:
						context.PrintErr(e.Pos(), "'make' has too many arguments.")
					case arg_len < 2:
						context.PrintErr(e.Pos(), "'make' has too few arguments.")
					case left_len > 1:
						context.PrintErr(e.Pos(), "'make' only returns one value.")
					}
					return
				}
				
				// a func call returning multiple items as a decl + init.
				switch n.Tok {
				case token.DEFINE:
					decl_stmt := new(ast.DeclStmt)
					gen_decl := new(ast.GenDecl)
					gen_decl.Tok = token.VAR
					
					// first we get each name of a var and then map them to a type.
					var_map := make(map[types.Type][]ast.Expr)
					for _, e := range n.Lhs {
						if type_expr := context.TypeInfo.TypeOf(e); type_expr != nil {
							var_map[type_expr] = append(var_map[type_expr], e)
						} else {
							context.PrintErr(n.TokPos, "Failed to expand assignment statement.")
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
					*owner_list = InsertToIndex[ast.Stmt](*owner_list, 0, decl_stmt)
					n.Tok = token.ASSIGN
				}
			}
		}
	}
}