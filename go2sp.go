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
	"./ast_to_sp"
)

func main() {
	files := os.Args[1:]
	for _, file := range files {
		fset := token.NewFileSet()
		code, err1 := ioutil.ReadFile("./" + file)
		CheckErr(err1)
		f, err2 := parser.ParseFile(fset, file, code, parser.AllErrors /*| parser.ParseComments*/)
		if err2 != nil {
			for _, e := range err2.(scanner.ErrorList) {
				fmt.Println(e)
			}
		} else {
			AST2SP.AddSourceGoTypes(f)
			var typeErrors []error
			conf := types.Config{Importer: importer.Default(), Error: func(err error) {
				typeErrors = append(typeErrors, err)
			}}
			info := &types.Info{
				Types:      make(map[ast.Expr]types.TypeAndValue), 
				Defs:       make(map[*ast.Ident]types.Object),
			}
			if _, err := conf.Check("", fset, []*ast.File{f}, info); err != nil {
				for _, e := range typeErrors {
					fmt.Println(e) /// type error
				}
			}
			
			fmt.Println(fmt.Sprintf("SourceGo: '%s' transpiled successfully as '%s.sp'", file, file))
			AST2SP.AnalyzeFile(f, info)
			//if err := WriteToFile(file + ".sp", sp_gen.Finalize()); err != nil {
			//	fmt.Println(fmt.Sprintf("SourceGo: unable to generate file '%s'.sp, %s", file), err)
			//}
			AST2SP.PrintAST(f)
			AST2SP.PrettyPrintAST(f)
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