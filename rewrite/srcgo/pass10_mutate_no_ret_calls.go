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
	"go/types"
	"go/ast"
	///"go/token"
)

type pass10Ctxt struct {
	currFunc *ast.FuncDecl
	transmit *AstTransmitter
}

func (a *AstTransmitter) MutateNoRetCalls(file *ast.File) {
	context := pass10Ctxt{ transmit: a }
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			context.currFunc = d
			if d.Body != nil {
				MutateBlock(d.Body, MutateNoRetCallStmts, &context)
			}
			context.currFunc = nil
		}
	}
}

func MutateNoRetCallStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator, ctxt any) {
	switch n := s.(type) {
	case *ast.BlockStmt:
		bm(n, MutateNoRetCallStmts, ctxt)
	
	case *ast.ForStmt:
		if n.Init != nil {
			MutateNoRetCallStmts(owner_list, index, n.Init, bm, ctxt)
		}
		if n.Post != nil {
			MutateNoRetCallStmts(owner_list, index, n.Post, bm, ctxt)
		}
		bm(n.Body, MutateNoRetCallStmts, ctxt)
	
	case *ast.IfStmt:
		if n.Init != nil {
			MutateNoRetCallStmts(owner_list, index, n.Init, bm, ctxt)
		}
		bm(n.Body, MutateNoRetCallStmts, ctxt)
		if n.Else != nil {
			MutateNoRetCallStmts(owner_list, index, n.Else, bm, ctxt)
		}
	
	case *ast.SwitchStmt:
		MutateNoRetCallStmts(owner_list, index, n.Init, bm, ctxt)
		bm(n.Body, MutateNoRetCallStmts, ctxt)
	
	case *ast.CaseClause:
		for i, stmt := range n.Body {
			MutateNoRetCallStmts(&n.Body, i, stmt, bm, ctxt)
		}
	
	case *ast.RangeStmt:
		bm(n.Body, MutateNoRetCallStmts, ctxt)
	
	case *ast.ExprStmt:
		context := ctxt.(*pass10Ctxt)
		if fn, is_func_call := n.X.(*ast.CallExpr); is_func_call {
			if typ := context.transmit.TypeInfo.TypeOf(fn); typ != nil {
				switch t := typ.(type) {
				case *types.Tuple:
					extra_args := t.Len()
					if IsFuncPtr(fn, context.transmit.TypeInfo) {
						if extra_args > 1 {
							retvals := make([]ast.Expr, 0)
							rettypes := make([]types.Type, 0)
							for i := 0; i < extra_args; i++ {
								ret_tmp := ast.NewIdent(fmt.Sprintf("fptr_temp%d", context.transmit.TmpVar))
								context.transmit.TmpVar++
								declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, nil, t.At(i).Type(), context.transmit.TypeInfo)
								*owner_list = InsertToIndex[ast.Stmt](*owner_list, 0, declstmt)
								retvals, rettypes = append(retvals, ret_tmp), append(rettypes, t.At(i).Type())
							}
							
							calls := ExpandFuncPtrCalls(fn, retvals, rettypes, context.transmit.TypeInfo, context.currFunc)
							for i := len(calls)-1; i > 0; i-- {
								*owner_list = InsertToIndex[ast.Stmt](*owner_list, FindStmt(*owner_list, n)+1, calls[i])
							}
							n.X = calls[0].(*ast.ExprStmt).X
						} else {
							calls := ExpandFuncPtrCalls(fn, nil, nil, context.transmit.TypeInfo, context.currFunc)
							n.X = calls[0].(*ast.ExprStmt).X
							for i := len(calls)-1; i > 0; i-- {
								*owner_list = InsertToIndex[ast.Stmt](*owner_list, FindStmt(*owner_list, n)+1, calls[i])
							}
						}
					} else {
						if extra_args > 1 {
							for i := 1; i < extra_args; i++ {
								ret_tmp := ast.NewIdent(fmt.Sprintf("fn_temp%d", context.transmit.TmpVar))
								context.transmit.TmpVar++
								declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, nil, t.At(i).Type(), context.transmit.TypeInfo)
								*owner_list = InsertToIndex[ast.Stmt](*owner_list, 0, declstmt)
								fn.Args = append(fn.Args, MakeReference(ret_tmp))
							}
						}
					}
				
				default:
					if IsFuncPtr(fn, context.transmit.TypeInfo) {
						ret_tmp := ast.NewIdent(fmt.Sprintf("fptr_temp%d", context.transmit.TmpVar))
						context.transmit.TmpVar++
						declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, nil, t, context.transmit.TypeInfo)
						*owner_list = InsertToIndex[ast.Stmt](*owner_list, 0, declstmt)
						
						calls := ExpandFuncPtrCalls(fn, []ast.Expr{ret_tmp}, []types.Type{t}, context.transmit.TypeInfo, context.currFunc)
						n.X = calls[0].(*ast.ExprStmt).X
						for i := len(calls)-1; i > 0; i-- {
							*owner_list = InsertToIndex[ast.Stmt](*owner_list, FindStmt(*owner_list, n)+1, calls[i])
						}
					}
				}
			}
		}
	}
}