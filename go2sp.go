/**
 * go2sp.go
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

package main


import (
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"go/token"
	"go/scanner"
	"go/parser"
	"go/importer"
	"go/types"
	"go/ast"
	"./srcgo/ast_transform"
	"./srcgo/ast_to_sp"
)


const (
	Flag_Debug = (iota + 1) << 1
	Flag_Force
)

func main() {
	files := os.Args[1:]
	SrcGo_ASTMod.AddSrcGoTypes()
	var opts int
	for _, file := range files {
		var bad_compile bool
		switch file {
			case "--debug", "-dbg":
				opts |= Flag_Debug
			case "-f", "--force", "--force-gen":
				opts |= Flag_Force
			case "--help", "-h":
				fmt.Println("SourceGo Usage: " + os.Args[0] + " [options] files... | options: [--debug, --force, --help, --version]")
			case "--version", "-v":
				fmt.Println("SourceGo version: v0.19a")
			default:
				fset := token.NewFileSet()
				code, err1 := ioutil.ReadFile(file)
				CheckErr(err1)
				f, err2 := parser.ParseFile(fset, file, code, parser.AllErrors /*| parser.ParseComments*/)
				if err2 != nil {
					for _, e := range err2.(scanner.ErrorList) {
						fmt.Println(e)
					}
					bad_compile = true
				} else {
					var typeErrs, transpileErrs []error
					conf := types.Config{
						Importer: importer.Default(),
						DisableUnusedImportCheck: true,
						Error: func(err error) {
							typeErrs = append(typeErrs, err)
						},
					}
					info := &types.Info{
						Types:      make(map[ast.Expr]types.TypeAndValue), 
						Defs:       make(map[*ast.Ident]types.Object),
						Uses:       make(map[*ast.Ident]types.Object),
						Implicits:  make(map[ast.Node]types.Object),
						Scopes:     make(map[ast.Node]*types.Scope),
						Selections: make(map[*ast.SelectorExpr]*types.Selection),
					}
					
					if _, err := conf.Check("", fset, []*ast.File{f}, info); err != nil {
						for _, e := range typeErrs {
							fmt.Println(e) /// type error
						}
						bad_compile = true
					}
					SrcGo_ASTMod.ASTCtxt.FSet = fset
					SrcGo_ASTMod.AnalyzeFile(f, info, func(err error) {
						transpileErrs = append(transpileErrs, err)
					})
					for _, e := range transpileErrs {
						fmt.Println(e)
					}
					
					/// do second type check.
					conf.Check("", fset, []*ast.File{f}, info)
					if (opts & Flag_Debug) > 0 {
						WriteToFile(fmt.Sprintf("%s_AST.txt", file), SrcGo_ASTMod.PrintAST(f))
						WriteToFile(fmt.Sprintf("%s_PrettyPrintAST.txt", file), SrcGo_ASTMod.PrettyPrintAST(f))
					}
				}
				if bad_compile && (opts & Flag_Force)==0 {
					fmt.Println(fmt.Sprintf("SourceGo: file '%s'.sp generation FAILED.", file))
				} else {
					final_code := SrcGoSPGen.GenSPFile(f)
					WriteToFile(file + ".sp", final_code)
					fmt.Println(fmt.Sprintf("SourceGo: successfully transpiled '%s.sp'.", file))
				}
		}
	}
}

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteToFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}