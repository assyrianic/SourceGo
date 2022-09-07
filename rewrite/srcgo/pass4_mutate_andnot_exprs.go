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
	"go/token"
	"go/ast"
)


/*
 * Pass #4 - Mutate AND-NOT expressions.
 * So we basically turn `a &^ b` into `a & ^(b)`
 * which then becomes `a & ~(b)` in SourcePawn.
 *
 * Ditto for the `a &^= b` as well
 * `a &= ^(b)` becoming `a &= ~(b)`
 */
func MutateAndNotExpr(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch x := n.(type) {
			case *ast.BinaryExpr:
				if x.Op==token.AND_NOT {
					x.Op = token.AND
					x.Y = MakeBitNotExpr(MakeParenExpr(x.Y))
				}
			case *ast.AssignStmt:
				if x.Tok==token.AND_NOT_ASSIGN {
					x.Tok = token.AND_ASSIGN
					for i := range x.Rhs {
						x.Rhs[i] = MakeBitNotExpr(MakeParenExpr(x.Rhs[i]))
					}
				}
			}
		}
		return true
	})
}