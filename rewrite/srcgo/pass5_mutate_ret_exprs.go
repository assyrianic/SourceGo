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
	"strings"
	"go/types"
	"go/ast"
	///"go/token"
)


type pass5Ctxt struct {
	funcMap  map[string]*ast.FuncDecl
	currFunc *ast.FuncDecl
	tyinfo   *types.Info
	tmpVar    uint
}

/* Case Studies of Returning statements to transform:
 * 
 * return m3() /// func m3() (type, type, type)
 * 
 * return int, float, m1() /// func m1() type
 */
func (a *AstTransmitter) MutateRetExprs(file *ast.File) {
	context := pass5Ctxt{ funcMap: CollectFuncNames(file), tyinfo: a.TypeInfo, tmpVar: a.TmpVar }
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			context.currFunc = d
			if d.Body != nil {
				MutateBlock(d.Body, MutateRetStmts, &context)
			}
			context.currFunc = nil
		}
	}
	a.TmpVar = context.tmpVar
}

func MutateRetStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator, ctxt any) {
	switch n := s.(type) {
	case *ast.BlockStmt:
		bm(n, MutateRetStmts, ctxt)
	
	case *ast.ForStmt:
		if n.Init != nil {
			MutateRetStmts(owner_list, index, n.Init, bm, ctxt)
		}
		if n.Post != nil {
			MutateRetStmts(owner_list, index, n.Post, bm, ctxt)
		}
		bm(n.Body, MutateRetStmts, ctxt)
	
	case *ast.IfStmt:
		if n.Init != nil {
			MutateRetStmts(owner_list, index, n.Init, bm, ctxt)
		}
		bm(n.Body, MutateRetStmts, ctxt)
		if n.Else != nil {
			MutateRetStmts(owner_list, index, n.Else, bm, ctxt)
		}
	
	case *ast.ReturnStmt:
		index := func(a []ast.Stmt, x ast.Stmt) int {
			for i, n := range a {
				if x==n {
					return i
				}
			}
			return -1
		}(*owner_list, s)
		
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
		context := ctxt.(*pass5Ctxt)
		for i:=1; i<res; i++ {
			call := func_calls[i]
			if call != nil && IsFuncPtr(call.Fun, context.tyinfo) {
				calls := ExpandFuncPtrCalls(call, []ast.Expr{ast.NewIdent(fmt.Sprintf("%s_param%d", context.currFunc.Name.Name, i))}, nil, context.tyinfo, context.currFunc)
				for i:=len(calls)-1; i>=0; i-- {
					*owner_list = InsertToIndex[ast.Stmt](*owner_list, index, calls[i])
				}
			} else {
				ptr_deref := PtrizeExpr(ast.NewIdent(fmt.Sprintf("%s_param%d", context.currFunc.Name.Name, i)))
				assign := MakeAssign(false)
				assign.Lhs = append(assign.Lhs, ptr_deref)
				assign.Rhs = append(assign.Rhs, n.Results[i])
				*owner_list = InsertToIndex[ast.Stmt](*owner_list, index, assign)
				index++
			}
		}
		
		if res_len > 1 {
			n.Results = n.Results[:1]
			MutateRetStmts(owner_list, index, s, bm, ctxt)
		} else if res_len==1 {
			// check for function call.
			ast.Inspect(n.Results[0], func(node ast.Node) bool {
				if node != nil {
					switch call := node.(type) {
					case *ast.CallExpr:
						if IsFuncPtr(call.Fun, context.tyinfo) {
							ret_tmp := ast.NewIdent(fmt.Sprintf("fptr_temp%d", context.tmpVar))
							context.tmpVar++
							declstmt := MakeVarDecl([]*ast.Ident{ret_tmp}, n.Results[0], nil, context.tyinfo)
							calls := ExpandFuncPtrCalls(call, []ast.Expr{ret_tmp}, nil, context.tyinfo, context.currFunc)
							calls = InsertToIndex[ast.Stmt](calls, 0, declstmt)
							for i:=len(calls)-1; i>=0; i-- {
								*owner_list = InsertToIndex[ast.Stmt](*owner_list, index, calls[i])
							}
							n.Results[0] = ret_tmp
						} else if fn, found := context.funcMap[GetFuncName(call.Fun)]; found {
							if arg_count, param_count := len(call.Args), len(fn.Type.Params.List); arg_count < param_count {
								for _, param := range context.currFunc.Type.Params.List {
									for _, name := range param.Names {
										if strings.HasPrefix(name.Name, context.currFunc.Name.Name + "_param") {
											call.Args = append(call.Args, name)
										}
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
		MutateRetStmts(owner_list, index, n.Init, bm, ctxt)
		bm(n.Body, MutateRetStmts, ctxt)
	
	case *ast.CaseClause:
		for i, stmt := range n.Body {
			MutateRetStmts(&n.Body, i, stmt, bm, ctxt)
		}
	
	case *ast.RangeStmt:
		bm(n.Body, MutateRetStmts, ctxt)
	}
}