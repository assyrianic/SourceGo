/**
 * sourcemod/adt_trie.go
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


type (
	StringMap struct {
		Size int
	}
	
	StringMapSnapshot struct {
		Length int
	}
)

func CreateTrie() StringMap
func (StringMap) SetValue(key string, value any) bool
func (StringMap) SetArray(key string, array []any, num_items int) bool
func (StringMap) SetString(key, value string) bool
func (StringMap) GetValue(key string, value *any) bool
func (StringMap) GetArray(key string, array []any, max_size int, size *int) bool
func (StringMap) GetString(key string, value []char, max_size int, size *int) bool
func (StringMap) Remove(key string)
func (StringMap) Clear()
func (StringMap) Snapshot() StringMapSnapshot

func (StringMapSnapshot) KeyBufferSize(index int) int
func (StringMapSnapshot) GetKey(index int, buffer []char, maxlength int) int