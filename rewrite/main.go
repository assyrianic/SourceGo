/*
 * main.go
 *
 * Copyright 2022 Nirari Technologies.
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
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/scanner"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	///"path/filepath"
	"os"
	"os/exec"
	"runtime"
	///"strings"
	
	"./srcgo"
)



func RunTypeCheck(fset *token.FileSet, file *ast.File) *types.Info {
	curr_dir, _ := os.Getwd()
	var errs []error
	conf := types.Config{
		DisableUnusedImportCheck: true,
		Importer: importer.ForCompiler(fset, "source", nil), // YOU NEED THIS.
		Error: func(err error) {
			errs = append(errs, err)
		},
	}
	
	ti := &types.Info{
		Types:     make(map[ast.Expr]types.TypeAndValue),
		Defs:      make(map[*ast.Ident]types.Object),
		Uses:      make(map[*ast.Ident]types.Object),
		Implicits: make(map[ast.Node]types.Object),
		Scopes:     make(map[ast.Node]*types.Scope),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	_, err := conf.Check(curr_dir, fset, []*ast.File{file}, ti)
	if err != nil {
		for _, type_err := range errs {
			fmt.Printf("SrcGo Error: **** '%s' ****\n", type_err)
		}
		return nil
	}
	return ti
}

func main() {
	for _, arg_str := range os.Args[1:] {
		switch arg_str {
		default:
			fset := token.NewFileSet()
			code, _ := LoadFile(arg_str)
			file_ast, parse_err := parser.ParseFile(fset, arg_str, code, parser.AllErrors)
			if parse_err != nil {
				for _, e := range parse_err.(scanner.ErrorList) {
					fmt.Println(e)
				}
				return
			}
			
			ti := RunTypeCheck(fset, file_ast)
			if ti==nil {
				return
			}
			
			transmitter := SrcGo.MakeAstTransmitter(fset, ti, nil, true)
			if !transmitter.AnalyzeIllegalCode(file_ast) {
				transmitter.PrintErrs()
				return
			}
			
			// this one doesn't need transmitter data.
			SrcGo.NameAnonFuncs(file_ast)
			
			// do another type check just in case.
			ti = RunTypeCheck(fset, file_ast)
			if ti==nil {
				return
			}
			transmitter.TypeInfo = ti
			
			SrcGo.MergeRetTypes(file_ast)
			SrcGo.MutateAndNotExpr(file_ast)
			
			transmitter.MutateRetExprs(file_ast)
			transmitter.MutateAssignDecls(file_ast)
			transmitter.MutateAssigns(file_ast)
			
			SrcGo.ChangeReceiverNames(file_ast)
			transmitter.MutateRanges(file_ast)
			transmitter.MutateNoRetCalls(file_ast)
			
			ti = RunTypeCheck(fset, file_ast)
			if ti==nil {
				return
			}
			
			fmt.Printf("'%s'\nEverything Ay-Ok!\n", SrcGo.PrettyPrintAST(file_ast))
		}
	}
}

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadFile(filename string) (string, string) {
	if text, read_err := ioutil.ReadFile(filename); read_err==nil {
		return string(text), "none"
	} else {
		return "", read_err.Error()
	}
}

func WriteToFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	if _, err = io.WriteString(file, data); err != nil {
		return err
	}
	return file.Sync()
}

func InvokeSPComp(file string) {
	cmd_str := "compile"
	if runtime.GOOS == "windows" {
		cmd_str += ".bat"
	} else {
		cmd_str = fmt.Sprintf("./%s.sh", cmd_str)
	}

	if msg, err := exec.Command(cmd_str).Output(); err == nil {
		fmt.Printf("SourceGo::SPComp Invoked:: %s\n", string(msg))
	}
}