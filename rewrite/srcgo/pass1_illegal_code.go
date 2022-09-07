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

/*
 * Pass #1 - Analyze Illegal Code
 * checks for ANY Go constructs that cannot be faithfully recreated in SourcePawn.
 */
func (a *AstTransmitter) AnalyzeIllegalCode(file *ast.File) bool {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Recv != nil && len(x.Recv.List) > 1 {
					a.PrintErr(x.Pos(), "Multiple Receivers are not allowed.")
				}
				if x.Type.Results != nil {
					for _, ret := range x.Type.Results.List {
						if ptr, is_ptr := ret.Type.(*ast.StarExpr); is_ptr {
							a.PrintErr(ptr.Pos(), "Returning Pointers isn't Allowed." + fmt.Sprintf(" Param %v is of pointer type", ret.Names))
						}
					}
				}
				
			case *ast.FuncType:
				if x.Results != nil {
					for _, ret := range x.Results.List {
						if ptr, is_ptr := ret.Type.(*ast.StarExpr); is_ptr {
							a.PrintErr(ptr.Pos(), "Returning Pointers isn't Allowed." + fmt.Sprintf(" Param %v is of pointer type", ret.Names))
						}
					}
				}
			
			case *ast.StructType:
				for _, f := range x.Fields.List {
					switch t := f.Type.(type) {
					case *ast.StarExpr:
						a.PrintErr(t.Pos(), "Pointers are not allowed in Structs.")
					case *ast.ArrayType:
						if t.Len==nil {
							a.PrintErr(t.Pos(), "Arrays of unknown size are not allowed in Structs.")
						}
					}
				}
			case *ast.BranchStmt:
				if x.Tok==token.GOTO || x.Tok==token.FALLTHROUGH {
					a.PrintErr(x.Pos(), fmt.Sprintf("'%s' is Illegal.", x.Tok.String()))
				} else if x.Label != nil {
					a.PrintErr(x.Pos(), "Branched Labels are Illegal.")
				}
			
			case *ast.CommClause: // case chan<-var
				a.PrintErr(x.Pos(), "Comm Select Cases are Illegal.")
			case *ast.DeferStmt:  // defer func()
				a.PrintErr(x.Pos(), "Defer Statements are Illegal.")
			case *ast.TypeSwitchStmt: // switch a.(type) {}
				a.PrintErr(x.Pos(), "Type-Switches are Illegal.")
			case *ast.LabeledStmt:
				a.PrintErr(x.Pos(), "Labels are Illegal.")
			case *ast.GoStmt: // go func()
				a.PrintErr(x.Pos(), "Goroutines are Illegal.")
			case *ast.SelectStmt: // select { case chan<-var: }
				a.PrintErr(x.Pos(), "Select Statements are Illegal.")
			case *ast.SendStmt: // chan <- var, var = <-chan
				a.PrintErr(x.Pos(), "Send Statements are Illegal.")
			case *ast.ChanType:
				a.PrintErr(x.Pos(), "Channel Types are Illegal.")
			
			case *ast.BasicLit:
				if x.Kind==token.IMAG {
					a.PrintErr(x.Pos(), "Imaginary Numbers are Illegal.")
				}
			case *ast.TypeAssertExpr: // b,c := a.(type)
				a.PrintErr(x.Pos(), "Type Assertions are Illegal.")
			case *ast.SliceExpr: // a[ low : high : max ]
				a.PrintErr(x.Pos(), "Slice Expressions are Illegal.")
			case *ast.MapType: // map[K]V
				/// check if the key isn't 'string', only string keys are allowed.
				if typ, is_ident := x.Key.(*ast.Ident); !is_ident || typ.Name != "string" {
					a.PrintErr(x.Pos(), "Non-string Maps are Illegal.")
				}
			case *ast.UnaryExpr:
				if x.Op==token.ARROW {
					a.PrintErr(x.Pos(), "Channel Expressions are Illegal.")
				}
			}
		}
		return true
	})
	return len(a.Errors)==0
}