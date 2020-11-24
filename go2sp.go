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
	"strings"
	"./srcgo/ast_transform"
	"./srcgo/ast_to_sp"
)


const (
	Flag_Debug = (iota + 1) << 1
	Flag_Force
)


func DoImports(dir string, file *ast.File, fset *token.FileSet, pkgs map[string]*ast.File) []*ast.File {
	var ast_files []*ast.File
	ast_files = append(ast_files, file)
	for _, imp := range file.Imports {
		file_to_import := dir + "/" + strings.Replace(imp.Path.Value, "\"", "", -1) + ".go"
		if _, ok := pkgs[file_to_import]; ok {
			/// prevent multiple importing.
			continue
		}
		
		imp_ast, imp_err := parser.ParseFile(fset, file_to_import, nil, parser.DeclarationErrors)
		if imp_err != nil {
			switch err_type := imp_err.(type) {
				case *os.PathError:
					fmt.Println(err_type, imp.Path.Value)
				case scanner.ErrorList:
					for _, e := range err_type {
						fmt.Println(e, imp.Path.Value)
					}
			}
			return nil
		} else {
			pkgs[file_to_import] = imp_ast
			for _, more := range DoImports(dir, imp_ast, fset, pkgs) {
				ast_files = append(ast_files, more)
			}
		}
	}
	return ast_files
}

func main() {
	srcgo_args := os.Args[1:]
	ASTMod.AddSrcGoTypes()
	var opts int
	for _, argStr := range srcgo_args {
		var bad_compile bool
		switch argStr {
			case "--debug", "-dbg":
				opts |= Flag_Debug
			case "-file_ast", "--force", "--force-gen":
				opts |= Flag_Force
			case "--help", "-h":
				fmt.Println("SourceGo Usage: " + os.Args[0] + " [options] files... | options: [--debug, --force, --help, --version]")
			case "--version", "-v":
				fmt.Println("SourceGo version: v0.36a")
			default:
				fset := token.NewFileSet()
				code, err1 := ioutil.ReadFile(argStr)
				CheckErr(err1)
				/// parse the file and get a File AST Node.
				file_ast, err2 := parser.ParseFile(fset, argStr, code, parser.AllErrors /*| parser.ParseComments*/)
				if err2 != nil {
					for _, e := range err2.(scanner.ErrorList) {
						fmt.Println(e)
					}
					bad_compile = true
				} else {
					dir, _ := os.Getwd()
					pkgs := make(map[string]*ast.File)
					ast_files := DoImports(dir, file_ast, fset, pkgs)
					
					/*defer func() {
						if err := recover(); err != nil {
							fmt.Printf("RECOVERY: %T - %+v\n", err, err)
							//debug.PrintStack()
							//ast.Print(fset, file_ast)
						}
					}()*/
					var typeErrs, transpileErrs []error
					conf := types.Config{
						Importer: importer.Default(),
						DisableUnusedImportCheck: true,
						Error: func(err error) {
							if strings.Contains(err.Error(), "could not import") {
							} else if strings.Contains(err.Error(), "declared but not used") {
								fmt.Printf("%s %20s\n", err, "[warning]")
							} else {
								typeErrs = append(typeErrs, err)
								bad_compile = true
							}
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
					
					/// Do initial type-check of the File AST Node so we can get type information.
					if _, err := conf.Check(``, fset, ast_files, info); err != nil {
						for _, e := range typeErrs {
							fmt.Printf("%s %20s\n", e, "[error]")
						}
					}
					
					/// initialize our transpiler.
					ASTMod.SetUpSrcGo(fset, info, func(err error) {
						transpileErrs = append(transpileErrs, err)
						bad_compile = true
					})
					
					/// first step: Analyze for illegal golang constructs.
					ASTMod.AnalyzeIllegalCode(file_ast)
					
					ASTMod.MergeRetVals(file_ast)
					
					ASTMod.ChangeRecvrNames(file_ast)
					
					ASTMod.MutateAndNotExpr(file_ast)
					
					ASTMod.MutateRets(file_ast)
					
					ASTMod.MutateAssigns(file_ast)
					
					ASTMod.MutateRanges(file_ast)
					
					ASTMod.MutateNoRetCalls(file_ast)
					
					/// TODO: for for loop inits that have multiple vars.
					//ASTMod.MutateForInits(file_ast)
					
					for _, e := range transpileErrs {
						fmt.Printf("%s %20s\n", e, "[error]")
					}
					
					conf.Check(``, fset, ast_files, info)
					if opts & Flag_Debug > 0 {
						WriteToFile(fmt.Sprintf("%s_AST.txt",   argStr), ASTMod.PrintAST(file_ast))
						WriteToFile(fmt.Sprintf("%s_output.go", argStr), ASTMod.PrettyPrintAST(file_ast))
					}
				}
				new_file_name := fmt.Sprintf("%s.sp", argStr)
				if bad_compile && opts & Flag_Force==0 {
					fmt.Println(fmt.Sprintf("SourceGo: file '%s' generation FAILED.", new_file_name))
				} else {
					final_code := GoToSPGen.GeneratePluginFile(file_ast) //GoToSPGen.GenSPFile(file_ast)
					WriteToFile(argStr + ".sp", final_code)
					if bad_compile {
						fmt.Println("SourceGo: transpiled " + new_file_name + " but might need correction.")
					} else {
						fmt.Println("SourceGo: successfully transpiled " + new_file_name)
					}
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