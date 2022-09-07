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
	"fmt"
	///"go/types"
	"go/ast"
	"go/token"
)


type pass7Ctxt struct {
	currFunc *ast.FuncDecl
	transmit *AstTransmitter
}

/* Case Studies of Changing non-defining assignments:
 * 
 */
func (a *AstTransmitter) MutateAssigns(file *ast.File) {
	context := pass7Ctxt{ transmit: a }
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			context.currFunc = d
			if d.Body != nil {
				MutateBlock(d.Body, MutateAssignStmts, &context)
			}
			context.currFunc = nil
		}
	}
}

func MutateAssignStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator, ctxt any) {
	switch n := s.(type) {
	case *ast.BlockStmt:
		bm(n, MutateAssignStmts, ctxt)
	
	case *ast.ForStmt:
		if n.Init != nil {
			MutateAssignStmts(owner_list, index, n.Init, bm, ctxt)
		}
		if n.Post != nil {
			MutateAssignStmts(owner_list, index, n.Post, bm, ctxt)
		}
		bm(n.Body, MutateAssignStmts, ctxt)
	
	case *ast.IfStmt:
		if n.Init != nil {
			MutateAssignStmts(owner_list, index, n.Init, bm, ctxt)
		}
		bm(n.Body, MutateAssignStmts, ctxt)
		if n.Else != nil {
			MutateAssignStmts(owner_list, index, n.Else, bm, ctxt)
		}
	
	case *ast.SwitchStmt:
		MutateAssignStmts(owner_list, index, n.Init, bm, ctxt)
		bm(n.Body, MutateAssignStmts, ctxt)
	
	case *ast.CaseClause:
		for i, stmt := range n.Body {
			MutateAssignStmts(&n.Body, i, stmt, bm, ctxt)
		}
	
	case *ast.RangeStmt:
		bm(n.Body, MutateAssignStmts, ctxt)
	
	case *ast.AssignStmt:
		context := ctxt.(*pass7Ctxt)
		left_len, rite_len := len(n.Lhs), len(n.Rhs)
		if rite_len==1 && left_len >= rite_len {
			switch n.Tok {
			case token.ASSIGN:
				switch fn := n.Rhs[0].(type) {
				case *ast.CallExpr:
					// a func call returning multiple items as a decl + init.
					if IsFuncPtr(fn, context.transmit.TypeInfo) {
						ret_tmp := ast.NewIdent(fmt.Sprintf("fptr_temp%d", context.transmit.TmpVar))
						context.transmit.TmpVar++
						declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, n.Lhs[0], nil, context.transmit.TypeInfo)
						
						retvals := make([]ast.Expr, 0)
						retvals = append(retvals, ret_tmp)
						for i := 1; i < left_len; i++ {
							retvals = append(retvals, n.Lhs[i])
						}
						
						calls := ExpandFuncPtrCalls(fn, retvals, nil, context.transmit.TypeInfo, context.currFunc)
						calls = InsertToIndex[ast.Stmt](calls, 0, declstmt)
						for i := len(calls)-1; i >= 0; i-- {
							*owner_list = InsertToIndex[ast.Stmt](*owner_list, index, calls[i])
						}
						
						n.Lhs = n.Lhs[:1]
						n.Rhs[0] = ret_tmp
					} else {
						// transform the tuple return into a single return + pass by ref.
						for i := 1; i < left_len; i++ {
							switch e := n.Lhs[i].(type) {
							case *ast.Ident:
								fn.Args = append(fn.Args, MakeReference(e))
							}
						}
						if left_len > 1 {
							n.Lhs = n.Lhs[:1]
						}
					}
				/**
				case *ast.IndexExpr:
					/// value, found = map[str]
					/// Becomes: found = map.GetValue(str, &value)
					if typ := context.TypeInfo.TypeOf(fn.X); typ != nil {
						switch t := typ.(type) {
						case *types.Map:
							ret_tmp := ast.NewIdent(fmt.Sprintf("map_value%d", context.TmpVar))
							context.TmpVar++
							declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, nil, t.Elem())
							
							new_stmts := make([]ast.Stmt, 0)
							map_access := new(ast.CallExpr)
							switch elem_typ := t.Elem().(type) {
								case 
								map_access.Fun = ast.NewIdent("GetValue")
							}
							map_access.Args = append(Call_StartFunction.Args, ast.NewIdent("nil"))
							map_access.Args = append(Call_StartFunction.Args, x.Fun)
						}
					}
				*/
				}
			}
		}
	}
}