/**
 * sourcemod/adt_array.go
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


func ByteCountToCells(size int) int

/// new ArrayList(int blocksize=1, int startsize=0);
type ArrayList struct {
	Length, BlockSize int
}

func CreateArray(blocksize, startsize int) ArrayList
func (ArrayList) Clear()
func (ArrayList) Clone() ArrayList
func (ArrayList) Resize(newsize int)
func (ArrayList) Push(value any) int
func (ArrayList) PushString(value string) int
func (ArrayList) PushArray(values []any, size int) int
func (ArrayList) Get(index, block int, asChar bool) any
func (ArrayList) GetString(index int, buffer []any, maxlength int) int
func (ArrayList) GetArray(index int, buffer *[]any, size int) int
func (ArrayList) Set(index int, value any, block int, asChar bool)
func (ArrayList) SetString(index int, value string)
func (ArrayList) SetArray(index int, values []any, size int)
func (ArrayList) ShiftUp(index int)
func (ArrayList) Erase(index int)
func (ArrayList) SwapAt(index1, index2 int)
func (ArrayList) FindString(item string) int
func (ArrayList) FindValue(item any, block int)
func (ArrayList) Sort(order SortOrder, sort SortType)
func (ArrayList) SortCustom(sorter SortFuncADTArray, hndl Handle)
