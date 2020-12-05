/**
 * sourcemod/sorting.go
 * 
 * Copyright 2020 Nirari Technologies, Alliedmodders LLC.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package main


type SortOrder int
const (
	Sort_Ascending = SortOrder(0)     /**< Ascending order */
	Sort_Descending    /**< Descending order */
	Sort_Random         /**< Random order */
)


type SortType int
const (
	Sort_Integer = SortType(0)
	Sort_Float
	Sort_String
)


func SortIntegers(array []int, array_size int, order SortOrder)
func SortFloats(array []float, array_size int, order SortOrder)
func SortStrings(array [][]char, array_size int, order SortOrder)


type SortFunc1D func(elem1, elem2 int, array []int, hndl Handle) int
func SortCustom1D(array []int, array_size int, sorter SortFunc1D, hndl Handle)

type SortFunc2D func(elem1, elem2 []any, array [][]any, hndl Handle)
func SortCustom2D(array [][]any, array_size int, sorter SortFunc2D, hndl Handle)
func SortADTArray(array Handle, order SortOrder, sorting SortType)

type SortFuncADTArray func(index1, index2 int, array, hndl Handle)
func SortADTArrayCustom(array Handle, sorter SortFuncADTArray, hndl Handle)