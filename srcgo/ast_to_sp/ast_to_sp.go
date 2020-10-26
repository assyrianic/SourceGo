/**
 * ast_to_so.go
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

package SrcGoSPGen


import (
	"strings"
	"go/token"
	"go/ast"
	"go/types"
	//"go/constant"
)


var type_info *types.Info


const (
	IncludeTemplate string = "#include <PATH>"
	VarTemplate string     = "<TYPE> <NAME>"
	FuncTemplate string    = "public <TYPE> <NAME>(<VARS>) <BLOCK>"
	BlockTemplate string   = "{ <STATEMENTS> }"
	StmtTemplate string    = "<STATEMENT>;"
	IfTemplate string      = "if(<COND>) <BLOCK>"
	ElIfTemplate string    = "else <IF_STMT> <BLOCK>"
	ElseTemplate string    = "else <BLOCK>"
	ForTemplate string     = "for(<INIT>;<COND>;<POST>) <BLOCK>"
	SwitchTemplate string  = "switch(<COND>) <BLOCK>"
	CaseTemplate string    = "case <EXPRS>: <BLOCK>"
	DefaultTemplate string = "default: <BLOCK>"
	
	ReturnTemplate string  = "return <EXPR>"
	BinExprTemplate string = "<EXPR> <OP> <EXPR>"
	SelectTemplate string  = "<EXPR>.<EXPR>"
	UnaryTemplate string   = "<OP><EXPR>"
	ArrayTemplate string   = "<EXPR>[<EXPR>]"
	CallTemplate string    = "<EXPR>(<ARGS>)"
)

/// final code generation stuffs.
type CodeGen struct {
	Includes, Globals, Funcs strings.Builder
	Tabs uint
}

type CodeBlock struct {
	Tab uint
	Stmts []string
}

type CodeFunc struct {
	CodeBlock
	Header string
}

func GenSPFile(f *ast.File) string {
	var c = CodeGen{}
	c.Includes.WriteString("#include <sourcemod>\n")
	for _, d := range f.Decls {
		switch decl := d.(type) {
			case *ast.GenDecl:
				c.GenerateGenDecl(decl)
			
			//case *ast.FuncDecl:
			//	c.GenerateFuncDecl(decl)
		}
	}
	return c.Includes.String() + "\n" + c.Globals.String() + "\n" + c.Funcs.String()
}

func (c *CodeGen) GenerateGenDecl(g *ast.GenDecl) {
	switch g.Tok {
		case token.IMPORT:
			for _, spec := range g.Specs {
				imp := spec.(*ast.ImportSpec)
				if imp.Path.Value[1] == '.' {
					c.Includes.WriteString( strings.Replace(IncludeTemplate, "<PATH>", "\"" + imp.Path.Value[2:], -1) )
					c.Includes.WriteString("\n")
				} else {
					c.Includes.WriteString( strings.Replace(IncludeTemplate, "PATH", imp.Path.Value[1 : len(imp.Path.Value)-1], -1) )
					c.Includes.WriteString("\n")
				}
			}
		/*
		case token.CONST:
			fset := token.NewFileSet()
			for _, spec := range g.Specs {
				const_spec := spec.(*ast.ValueSpec)
				if typ := type_info.TypeOf(const_spec.Type); typ != nil {
					type_name := typ.String()
					//is_array_type := strings.ContainsRune(type_name, rune('['))
					type_name = strings.TrimFunc(type_name, func(r rune) bool {
						return !unicode.IsLetter(r)
					})
					temp := strings.Replace(VarTemplate, "<TYPE>", type_name, -1)
					for i:=0; i<len(const_spec.Names); i++ {
						temp += strings.Replace(temp, "<NAME>", type_name, -1)
						
						var buf bytes.Buffer
						format.Node(&buf, fset, const_spec.Values[i])
						temp += strings.Replace(" = <EXPR>", "<EXPR>", buf.String(), -1)
						if i+1 != len(const_spec.Names) {
							temp += ", "
						}
					}
					c.Globals.WriteString(temp + ";\n\n")
				}
				//Names   []*Ident      // value names (len(Names) > 0)
				//Type    Expr          // value type; or nil
				//Values  []Expr
			}
			**/
	}
}