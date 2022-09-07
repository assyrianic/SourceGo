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
	"go/token"
	"go/ast"
)


type passTwoCtxt struct {
	newDecls []ast.Decl
	tmpFunc    uint
}

/*
 * Pass #2 - Name Anonymous Functions
 * Any unnamed function literal is to be named so they can be generated
 * as an incrementally named function in SourcePawn.
 *
 * This AST pass has to be second because we later
 * mutate the multi-return types for functions!
 *
 * So best get to naming the lambdas before things get really serious!
 */
func NameAnonFuncs(file *ast.File) {
	/**
	 * Function Literals can be represented in different ways:
	 * Assigned:
	 * f := func(params){
	 *     code
	 * }
	 * f(args)
	 * 
	 * OR
	 * Inline-Called:
	 * func(params){
	 *     code
	 * }(args)
	 */
	context := passTwoCtxt{ newDecls: make([]ast.Decl, 0) }
	
	for _, decl := range file.Decls {
		context.newDecls = append(context.newDecls, decl)
	}
	
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Body != nil {
				MutateBlock(d.Body, MutateFuncLit, &context)
			}
		}
	}
	if len(file.Decls) < len(context.newDecls) {
		file.Decls = context.newDecls
	}
	context.newDecls = nil
}


/// func(params){code}(args) => func _srcgo_func#(params){code} ... _srcgo_func#(args)
func MutateFuncLit(owner_list *[]ast.Stmt, index int, s ast.Stmt, bm BlockMutator, ctxt any) {
	switch n := s.(type) {
	case *ast.BlockStmt:
		bm(n, MutateFuncLit, ctxt)
	
	case *ast.ForStmt:
		if n.Init != nil {
			MutateFuncLit(owner_list, index, n.Init, bm, ctxt)
		}
		if n.Cond != nil {
			MutateFuncLitExprs(&n.Cond, ctxt)
		}
		if n.Post != nil {
			MutateFuncLit(owner_list, index, n.Post, bm, ctxt)
		}
		bm(n.Body, MutateFuncLit, ctxt)
	
	case *ast.IfStmt:
		if n.Init != nil {
			MutateFuncLit(owner_list, index, n.Init, bm, ctxt)
		}
		MutateFuncLitExprs(&n.Cond, ctxt)
		bm(n.Body, MutateFuncLit, ctxt)
		if n.Else != nil {
			MutateFuncLit(owner_list, index, n.Else, bm, ctxt)
		}
	
	case *ast.SwitchStmt:
		MutateFuncLit(owner_list, index, n.Init, bm, ctxt)
		MutateFuncLitExprs(&n.Tag, ctxt)
		bm(n.Body, MutateFuncLit, ctxt)
	
	case *ast.CaseClause:
		for j := range n.List {
			MutateFuncLitExprs(&n.List[j], ctxt)
		}
		for i, stmt := range n.Body {
			MutateFuncLit(&n.Body, i, stmt, bm, ctxt)
		}
	
	case *ast.RangeStmt:
		MutateFuncLitExprs(&n.Key, ctxt)
		MutateFuncLitExprs(&n.Value, ctxt)
		MutateFuncLitExprs(&n.X, ctxt)
		bm(n.Body, MutateFuncLit, ctxt)
	
	case *ast.ExprStmt:
		MutateFuncLitExprs(&n.X, ctxt)
	
	case *ast.ReturnStmt:
		for i := range n.Results {
			MutateFuncLitExprs(&n.Results[i], ctxt)
		}
	
	case *ast.AssignStmt:
		for i := range n.Rhs {
			MutateFuncLitExprs(&n.Rhs[i], ctxt)
		}
		for i := range n.Lhs {
			MutateFuncLitExprs(&n.Lhs[i], ctxt)
		}
	
	case *ast.DeclStmt:
		g := n.Decl.(*ast.GenDecl)
		for _, d := range g.Specs {
			switch g.Tok {
			case token.CONST, token.VAR:
				v := d.(*ast.ValueSpec)
				for expr := range v.Values {
					MutateFuncLitExprs(&v.Values[expr], ctxt)
				}
			}
		}
	}
}

func MutateFuncLitExprs(e *ast.Expr, ctxt any) {
	if e==nil || *e == nil {
		return
	}
	switch n := (*e).(type) {
	case *ast.BinaryExpr:
		MutateFuncLitExprs(&n.X, ctxt)
		MutateFuncLitExprs(&n.Y, ctxt)
	
	case *ast.CallExpr:
		MutateFuncLitExprs(&n.Fun, ctxt)
		for i := range n.Args {
			MutateFuncLitExprs(&n.Args[i], ctxt)
		}
	
	case *ast.KeyValueExpr:
		MutateFuncLitExprs(&n.Key, ctxt)
		MutateFuncLitExprs(&n.Value, ctxt)
	
	case *ast.IndexExpr:
		MutateFuncLitExprs(&n.X, ctxt)
		MutateFuncLitExprs(&n.Index, ctxt)
	
	case *ast.UnaryExpr:
		MutateFuncLitExprs(&n.X, ctxt)
	
	case *ast.FuncLit:
		context := ctxt.(*passTwoCtxt)
		tmp_func_name := ast.NewIdent(fmt.Sprintf("srcGoAnonFunc%d", context.tmpFunc))
		context.tmpFunc++
		fn_decl := new(ast.FuncDecl)
		fn_decl.Name = tmp_func_name
		fn_decl.Type = n.Type
		n.Type = nil
		fn_decl.Body = n.Body
		n.Body = nil
		*e = tmp_func_name
		context.newDecls = append(context.newDecls, fn_decl)
	}
}