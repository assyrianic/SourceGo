/**
 * ast_to_sp.go
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

package ASTtoSP


import (
	"fmt"
	"strings"
	"go/token"
	//"go/scanner"
	//"go/parser"
	"go/ast"
	"go/types"   /// Golang's Type system.
)


type SPGen struct {
	Includes, Vars, Funcs  strings.Builder
	SrcGoAST               *ast.File
	Types, Patterns        map[string]string
	Syms                   map[string]types.Type
	Tabs                   uint
	GenInclude             bool
}

var sp_gen SPGen = SPGen{}

func AddTabs(sb *strings.Builder, tabs uint) {
	for i:=uint(0); i<tabs; i++ {
		sb.WriteString(" ")
	}
}


func (sp_gen *SPGen) AnalyzeFile() {
	/// basic types.
	sp_gen.Types = make(map[string]string)
	sp_gen.Types["int"] = "int"
	sp_gen.Types["int32"] = "int"
	sp_gen.Types["Entity"] = "int"
	sp_gen.Types["float"] = "float"
	sp_gen.Types["float32"] = "float"
	sp_gen.Types["bool"] = "bool"
	sp_gen.Types["string"] = "char[]"
	sp_gen.Types["char"] = "char"
	sp_gen.Types["int8"] = "char"
	sp_gen.Types["rune"] = "char"
	
	sp_gen.Types["Vec3"] = "float <id>[3];" /// <id> is replaced with identifier name.
	sp_gen.Types["Map"] = "StringMap"
	sp_gen.Types["Array"] = "ArrayList"
	sp_gen.Types["Obj"] = "Handle"
	
	sp_gen.Syms = make(map[string]types.Type)
	
	sp_gen.Includes.WriteString("/**\n * file generated by the GoToSourcePawn Transpiler v0.6a\n * Copyright 2020 (C) Kevin Yonan aka Nergal, Assyrianic.\n * GoToSourcePawn Project is licensed under MIT.\n * link: 'https://github.com/assyrianic/Go2SourcePawn'\n */\n\n")
	sp_gen.Includes.WriteString("#include <sourcemod>\n")
	
	sp_gen.Vars.WriteString("#pragma semicolon    1\n#pragma newdecls     required\n\n\n")
	
	for _, decl := range sp_gen.SrcGoAST.Decls {
		sp_gen.ManageDeclNode(decl)
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

func (sp_gen *SPGen) ManageDeclNode(d ast.Decl) {
	switch d.(type) {
		case *ast.GenDecl:
			sp_gen.AnalyzeGenDecl(d.(*ast.GenDecl))
		case *ast.FuncDecl:
			sp_gen.AnalyzeFuncDecl(d.(*ast.FuncDecl))
	}
}

/// Generic Declaration Node
func (sp_gen *SPGen) AnalyzeGenDecl(g *ast.GenDecl) {
	switch g.Tok {
		case token.IMPORT:
			for _, imp := range g.Specs {
				sp_gen.ReadImport(imp.(*ast.ImportSpec))
			}
			sp_gen.Includes.WriteString("\n")
		//case token.CONST:   *ast.ValueSpec
			/// careful on this one. strings have to be placed in defines.
		
		//case token.TYPE:    *ast.TypeSpec
		//case token.VAR:     *ast.ValueSpec
	}
}

func (sp_gen *SPGen) ReadImport(imp *ast.ImportSpec) {
	/// if we have a dot, assume it relative.
	if imp.Path.Value[1]=='.' {
		sp_gen.Includes.WriteString("#include \"" + imp.Path.Value[2:] + "\n")
	} else {
		sp_gen.Includes.WriteString("#include <" + imp.Path.Value[1 : len(imp.Path.Value)-1] + ">\n")
	}
}


func (sp_gen *SPGen) AnalyzeFuncDecl(f *ast.FuncDecl) {
	sp_gen.Funcs.WriteString("public ")
	if f.Type.Results==nil {
		sp_gen.Funcs.WriteString("void")
	} else {
		rettype := f.Type.Results.List[0].Type
		switch rettype.(type) {
			case *ast.Ident:
				iden := rettype.(*ast.Ident)
				sp_gen.Funcs.WriteString(iden.Name)
			case *ast.StarExpr:
				panic("SourceGo: You can't return a pointer/reference type.")
		}
	}
	sp_gen.Funcs.WriteString(" " + f.Name.Name)
	sp_gen.Funcs.WriteString("(")
	
	if f.Recv != nil {
		/// receivers act as first params.
		if len(f.Recv.List) > 1 {
			panic("SourceGo: receivers can only have one param.")
		} else {
			//sp_gen.ManageExprNode(f.Recv.List[0].Type)
			//sp_gen.Funcs.WriteString(" " + f.Recv.List[0].Names[0].Name)
			
			param_type := f.Recv.List[0].Type
			switch p := param_type.(type) {
				case *ast.Ident:
					sp_gen.Funcs.WriteString(p.Name + " " + f.Recv.List[0].Names[0].Name)
				
				case *ast.StarExpr:
					/// remember that enum structs don't take '&'
					iden := p.X.(*ast.Ident)
					sp_gen.Funcs.WriteString(iden.Name + "& " + f.Recv.List[0].Names[0].Name)
			}
			
			if len(f.Type.Params.List) != 0 {
				sp_gen.Funcs.WriteString(", ")
			}
		}
	}
	
	/// switch param names backwards and check for params combined as the same type.
	param_count := len(f.Type.Params.List)
	for n:=0; n<param_count; n++ {
		var type_str string
		param := f.Type.Params.List[n]
		switch t := param.Type.(type) {
			case *ast.Ident:
				type_str = t.Name
			case *ast.StarExpr:
				iden := t.X.(*ast.Ident)
				type_str = iden.Name + "&"
		}
		arg_count := len(param.Names)
		for i:=0; i<arg_count; i++ {
			sp_gen.Funcs.WriteString(type_str + " " + param.Names[i].Name)
			if i+1 != arg_count {
				sp_gen.Funcs.WriteString(", ")
			}
		}
		if n+1 != param_count {
			sp_gen.Funcs.WriteString(", ")
		}
	}
	sp_gen.Funcs.WriteString(")")
	if f.Body==nil {
		sp_gen.Funcs.WriteString(";")
	} else {
		sp_gen.AnalyzeBlockStmt(f.Body);
	}
	sp_gen.Funcs.WriteString("\n\n")
}

func (sp_gen *SPGen) ManageStmtNode(s ast.Stmt) {
	AddTabs(&sp_gen.Funcs, sp_gen.Tabs)
	switch n := s.(type) {
		case *ast.AssignStmt:
			sp_gen.AnalyzeAssignStmt(n)
		
		case *ast.BlockStmt:
			sp_gen.AnalyzeBlockStmt(n)
		
		case *ast.BranchStmt:
			sp_gen.AnalyzeBranchStmt(n)
		
		case *ast.DeclStmt:
			sp_gen.ManageDeclNode(n.Decl)
			sp_gen.Funcs.WriteString(";")
			
		case *ast.EmptyStmt:
			sp_gen.Funcs.WriteString(";")
			
		case *ast.ExprStmt:
			sp_gen.ManageExprNode(n.X)
			sp_gen.Funcs.WriteString(";")
			
		case *ast.ForStmt:
			sp_gen.AnalyzeForStmt(n)
			
		case *ast.IfStmt:
			sp_gen.AnalyzeIfStmt(n)
			
		case *ast.IncDecStmt:
			sp_gen.ManageExprNode(n.X)
			sp_gen.Funcs.WriteString(n.Tok.String() + ";")
			
		case *ast.ReturnStmt:
			sp_gen.Funcs.WriteString("return")
			if n.Results != nil {
				sp_gen.Funcs.WriteString(" ")
				/// change multiple var returns into passing by reference.
				sp_gen.ManageExprNode(n.Results[0])
			}
			sp_gen.Funcs.WriteString(";")
		
		case *ast.SwitchStmt:
			sp_gen.AnalyzeSwitch(n)
		
		case *ast.CaseClause:
			sp_gen.AnalyzeCaseClause(n)
		
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

func (sp_gen *SPGen) AnalyzeBlockStmt(b *ast.BlockStmt) {
	if b==nil {
		return
	}
	sp_gen.Funcs.WriteString(" {\n")
	sp_gen.Tabs++
	for _, stmt := range b.List {
		sp_gen.ManageStmtNode(stmt)
	}
	sp_gen.Tabs--
	sp_gen.Funcs.WriteString("\n")
	AddTabs(&sp_gen.Funcs, sp_gen.Tabs)
	sp_gen.Funcs.WriteString("}")
	if sp_gen.Tabs > 1 {
		sp_gen.Funcs.WriteString("\n")
	}
}

func (sp_gen *SPGen) AnalyzeForStmt(for_stmt *ast.ForStmt) {
	/// in Golang, 'for' replaces both for and while-loops.
	/// we'll have to replace while-loop like constructs with a degenerate for-loop
	sp_gen.Funcs.WriteString("for(")
	if for_stmt.Init != nil { /// initialization statement; or nil
		sp_gen.ManageStmtNode(for_stmt.Init)
	}
	sp_gen.Funcs.WriteString("; ")
	if for_stmt.Cond != nil { /// condition; or nil
		sp_gen.ManageExprNode(for_stmt.Cond)
	}
	sp_gen.Funcs.WriteString("; ")
	if for_stmt.Post != nil { /// post iteration statement; or nil
		sp_gen.ManageStmtNode(for_stmt.Post)
	}
	sp_gen.Funcs.WriteString(")")
	sp_gen.AnalyzeBlockStmt(for_stmt.Body)
}

func (sp_gen *SPGen) AnalyzeIfStmt(if_stmt *ast.IfStmt) {
	/// assumes tabs have been written to string builder.
	sp_gen.Funcs.WriteString("if(")
	if if_stmt.Init != nil { /// initialization statement; or nil
		sp_gen.ManageStmtNode(if_stmt.Init)
	}
	sp_gen.ManageExprNode(if_stmt.Cond)
	sp_gen.Funcs.WriteString(")")
	sp_gen.AnalyzeBlockStmt(if_stmt.Body)
	if if_stmt.Else != nil {
		sp_gen.Funcs.WriteString(" else ")
		sp_gen.ManageStmtNode(if_stmt.Else)
	}
}

func (sp_gen *SPGen) AnalyzeBranchStmt(b *ast.BranchStmt) {
	if b.Tok==token.GOTO || b.Tok==token.FALLTHROUGH {
		panic("SourceGo: " + fmt.Sprintf("%s is illegal.", b.Tok.String()))
	} else if b.Label != nil {
		panic("SourceGo: Branched Labels are illegal.")
	}
	sp_gen.Funcs.WriteString(strings.ToLower(b.Tok.String()) + ";")
}

func (sp_gen *SPGen) AnalyzeSwitch(s *ast.SwitchStmt) {
	sp_gen.Funcs.WriteString("switch(")
	//sp_gen.ManageStmtNode(s.Init)
	sp_gen.ManageExprNode(s.Tag)
	sp_gen.Funcs.WriteString(")")
	sp_gen.AnalyzeBlockStmt(s.Body)
	
	/**
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
	 * See if we can transform a true-switch into that for SourcePawn.
	 */
}

func (sp_gen *SPGen) AnalyzeCaseClause(c *ast.CaseClause) {
	exprs := len(c.List)
	if exprs > 0 {
		sp_gen.Funcs.WriteString("case ")
		for i:=0; i<exprs; i++ {
			sp_gen.ManageExprNode(c.List[i])
			if i+1 != exprs {
				sp_gen.Funcs.WriteString(", ")
			}
		}
		sp_gen.Funcs.WriteString(": {\n")
	} else {
		sp_gen.Funcs.WriteString("default: {\n")
	}
	
	sp_gen.Tabs++
	for _, stmt := range c.Body {
		sp_gen.ManageStmtNode(stmt)
		sp_gen.Funcs.WriteString("\n")
	}
	sp_gen.Tabs--
	AddTabs(&sp_gen.Funcs, sp_gen.Tabs)
	sp_gen.Funcs.WriteString("}")
}

func (sp_gen *SPGen) AnalyzeAssignStmt(a *ast.AssignStmt) {
	/// make sure to check if len(rhs) <= len(lhs).
	/// also check if rhs is function call expr.
	for _, e := range a.Lhs {
		sp_gen.ManageExprNode(e)
	}
	sp_gen.Funcs.WriteString(" " + a.Tok.String() + " ")
	for _, e := range a.Rhs {
		sp_gen.ManageExprNode(e)
	}
	sp_gen.Funcs.WriteString(";")
}

func (sp_gen *SPGen) ManageExprNode(e ast.Expr) {
	switch x := e.(type) {
		case *ast.IndexExpr:
			sp_gen.ManageExprNode(x.X)
			sp_gen.Funcs.WriteString("[")
			sp_gen.ManageExprNode(x.Index)
			sp_gen.Funcs.WriteString("]")
		
		case *ast.KeyValueExpr: /// for Map,Struct types.
			//Key   Expr
			//Value Expr
		
		case *ast.ParenExpr:
			sp_gen.Funcs.WriteString("(")
			sp_gen.ManageExprNode(x.X)
			sp_gen.Funcs.WriteString(")")
		
		case *ast.StarExpr:
			sp_gen.ManageExprNode(x.X)
		
		case *ast.UnaryExpr:
			sp_gen.Funcs.WriteString(x.Op.String())
			sp_gen.ManageExprNode(x.X)
		
		case *ast.CallExpr:
			sp_gen.ManageExprNode(x.Fun)
			sp_gen.Funcs.WriteString("(")
			args := len(x.Args)
			for i:=0; i<args; i++ {
				sp_gen.ManageExprNode(x.Args[i])
				if i+1 != args {
					sp_gen.Funcs.WriteString(", ")
				}
			}
			sp_gen.Funcs.WriteString(")")
		
		case *ast.BinaryExpr:
			sp_gen.ManageExprNode(x.X)
			if x.Op==token.AND_NOT {
				sp_gen.Funcs.WriteString(" & ~(")
				sp_gen.ManageExprNode(x.Y)
				sp_gen.Funcs.WriteString(")")
			} else {
				sp_gen.Funcs.WriteString(" " + x.Op.String() + " ")
				sp_gen.ManageExprNode(x.Y)
			}
		
		case *ast.Ident:
			//Obj *Object   // denoted object; or nil
			sp_gen.Funcs.WriteString(x.Name)
		
		case *ast.BasicLit:
			if x.Kind==token.IMAG {
				panic("SourceGo: Imaginary numbers are illegal.")
			}
			sp_gen.Funcs.WriteString(x.Value)
		
		case *ast.TypeAssertExpr:
			panic("SourceGo: Type Assertions are illegal.")
		case *ast.SelectorExpr:
			panic("SourceGo: Selectors are illegal.")
		case *ast.SliceExpr:
			panic("SourceGo: Slice Expressions are illegal.")
	}
}

func (sp_gen *SPGen) PrintAST() {
	ast.Inspect(sp_gen.SrcGoAST, func(n ast.Node) bool {
		if n != nil {
			fmt.Println(fmt.Sprintf("%T:\t\t", n), n)
		}
		return true
	})
}

func (sp_gen *SPGen) Finalize() string {
	return sp_gen.Includes.String() + sp_gen.Vars.String() + sp_gen.Funcs.String()
}