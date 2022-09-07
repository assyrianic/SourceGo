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
	"go/token"
)

type pass9Ctxt struct {
	currFunc *ast.FuncDecl
	ti       *types.Info
	rangeIter uint
}

func (a *AstTransmitter) MutateRanges(file *ast.File) {
	context := pass9Ctxt{ ti: a.TypeInfo } 
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			context.currFunc = d
			if d.Body != nil {
				MutateBlock(d.Body, MutateRangeStmts, &context)
			}
			context.currFunc = nil
		}
	}
}

func MutateRangeStmts(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator, ctxt any) {
	switch n := s.(type) {
	case *ast.BlockStmt:
		bm(n, MutateRangeStmts, ctxt)
	
	case *ast.ForStmt:
		if n.Init != nil {
			MutateRangeStmts(owner_list, index, n.Init, bm, ctxt)
		}
		if n.Post != nil {
			MutateRangeStmts(owner_list, index, n.Post, bm, ctxt)
		}
		bm(n.Body, MutateRangeStmts, ctxt)
	
	case *ast.IfStmt:
		if n.Init != nil {
			MutateRangeStmts(owner_list, index, n.Init, bm, ctxt)
		}
		bm(n.Body, MutateRangeStmts, ctxt)
		if n.Else != nil {
			MutateRangeStmts(owner_list, index, n.Else, bm, ctxt)
		}
	
	case *ast.SwitchStmt:
		MutateRangeStmts(owner_list, index, n.Init, bm, ctxt)
		bm(n.Body, MutateRangeStmts, ctxt)
	
	case *ast.CaseClause:
		for i, stmt := range n.Body {
			MutateRangeStmts(&n.Body, i, stmt, bm, ctxt)
		}
	
	case *ast.RangeStmt:
		context := ctxt.(*pass9Ctxt)
		if n.Key != nil {
			if iden, ok := n.Key.(*ast.Ident); ok && iden.Name=="_" {
				n.Key = ast.NewIdent(fmt.Sprintf("%s_iter%d", context.currFunc.Name.Name, context.rangeIter))
				context.rangeIter++
			}
		} else {
			n.Key = ast.NewIdent(fmt.Sprintf("%s_iter%d", context.currFunc.Name.Name, context.rangeIter))
			context.rangeIter++
		}
		
		if n.Value != nil {
			if iden, ok := n.Value.(*ast.Ident); ok && iden.Name=="_" {
				n.Value = nil
			} else {
				switch n.Tok {
				case token.DEFINE:
					decl_stmt := MakeVarDecl([]*ast.Ident{n.Value.(*ast.Ident)}, n.Value, nil, context.ti)
					n.Body.List = InsertToIndex[ast.Stmt](n.Body.List, 0, decl_stmt)
					
					assign := MakeAssign(false)
					assign.Lhs = append(assign.Lhs, n.Value)
					
					get_index := MakeIndex(n.Key, n.X)
					assign.Rhs = append(assign.Rhs, get_index)
					
					n.Body.List = InsertToIndex[ast.Stmt](n.Body.List, 1, assign)
					n.Value = nil
				}
			}
		}
		bm(n.Body, MutateRangeStmts, ctxt)
	}
}