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
	"go/ast"
)


/*
 * Pass #3 - Merge Return Types.
 */
func MergeRetTypes(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			switch f := n.(type) {
			case *ast.FuncDecl:
				f.Type.Params.List = MutateRetTypes(&f.Type.Results, f.Type.Params, f.Name.Name)
			case *ast.TypeSpec:
				if t, is_func_type := f.Type.(*ast.FuncType); is_func_type {
					t.Params.List = MutateRetTypes(&t.Results, t.Params, f.Name.Name)
				}
			}
		}
		return true
	})
}

/*
 * Modifies the return values of a function by mutating them into references and moving them to the parameters.
 * Example Go code: func f() (int, float) {}
 * Result  Go code: func f(f_param1 *float) int {}
 */
func MutateRetTypes(retvals **ast.FieldList, curr_params *ast.FieldList, obj_name string) []*ast.Field {
	if *retvals==nil || (*retvals).List==nil {
		return curr_params.List
	}
	
	new_params := make([]*ast.Field, 0)
	for _, param := range curr_params.List {
		new_params = append(new_params, param)
	}
	
	results := len((*retvals).List)
	
	// multiple different return values.
	if results > 1 {
		for i := 1; i<results; i++ {
			ret := (*retvals).List[i]
			// if they're named, treat as reference types.
			if ret.Names != nil && len(ret.Names) > 1 {
				ret.Type = PtrizeExpr(ret.Type)
				new_params = append(new_params, ret)
			} else {
				///param_num := len(new_params)
				ret.Names = append(ret.Names, ast.NewIdent(fmt.Sprintf("%s_param%d", obj_name, i)))
				ret.Type = PtrizeExpr(ret.Type)
				new_params = append(new_params, ret)
			}
		}
		(*retvals).List = (*retvals).List[:1]
	} else if results==1 && (*retvals).List[0].Names != nil && len((*retvals).List[0].Names) > 1 {
		// This condition can happen if there's multiple return values of the same type but they're named!
		(*retvals).List[0].Type = PtrizeExpr((*retvals).List[0].Type)
		new_params = append(new_params, (*retvals).List[0])
		*retvals = nil
	}
	return new_params
}