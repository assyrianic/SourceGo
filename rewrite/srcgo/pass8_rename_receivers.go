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
	///"fmt"
	///"go/types"
	"go/ast"
	///"go/token"
)

func ChangeReceiverNames(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch f := n.(type) {
			case *ast.FuncDecl:
				if f.Recv != nil && f.Recv.List[0].Names != nil && len(f.Recv.List[0].Names) > 0 {
					recvr := f.Recv.List[0].Names[0].Name
					ast.Inspect(f.Body, func(n ast.Node) bool {
						if n != nil {
							switch i := n.(type) {
							case *ast.Ident:
								if recvr==i.Name {
									i.Name = "this"
								}
							}
						}
						return true
					})
				}
			}
		}
		return true
	})
}