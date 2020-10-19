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
	"./ast_to_sp"
)


/// use '.gp' file ext and '.ginc'.
func main() {
	files := os.Args[1:]
	for _, file := range files {
		code, err1 := ioutil.ReadFile("./" + file)
		CheckErr(err1)
		fset := token.NewFileSet()
		f, err2 := parser.ParseFile(fset, file, code, parser.AllErrors)
		if err2 != nil {
			for _, e := range err2.(scanner.ErrorList) {
				fmt.Println(e)
			}
		} else {
			fmt.Println(fmt.Sprintf("SourceGo: '%s' transpiled successfully as '%s.sp'", file, file))
			sp_gen := ASTtoSP.SPGen{ SrcGoAST: f }
			sp_gen.PrintAST()
			sp_gen.AnalyzeFile()
			if err := WriteToFile(file + ".sp", sp_gen.Finalize()); err != nil {
				fmt.Println(fmt.Sprintf("SourceGo: unable to generate file '%s'.sp, %s", file), err)
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