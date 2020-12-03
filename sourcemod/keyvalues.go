/**
 * sourcemod/keyvalues.go
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
	KvDataTypes int
	KeyValues struct {
		ExportLength int
	}
)

const (
	KvData_None = KvDataTypes(0)    /**< Type could not be identified, or no type */
	KvData_String      /**< String value */
	KvData_Int         /**< Integer value */
	KvData_Float       /**< Floating point value */
	KvData_Ptr         /**< Pointer value (sometimes called "long") */
	KvData_WString     /**< Wide string value */
	KvData_Color       /**< Color value */
	KvData_UInt64      /**< Large integer value */
	/* --- */
	KvData_NUMTYPES
)

func CreateKeyValues(name, firstKey, firstValue string) KeyValues

func (KeyValues) ExportToFile(file string) bool
func (KeyValues) ExportToString(buffer []char, maxlength int) int
func (KeyValues) ImportFromFile(file string) bool
func (KeyValues) ImportFromString(buffer, resourceName string) bool
func (KeyValues) Import(other KeyValues)

func (KeyValues) SetString(key, value string) string
func (KeyValues) SetNum(key string, value int)
func (KeyValues) SetUInt64(key string, value [2]int)
func (KeyValues) SetFloat(key string, value float)
func (KeyValues) SetColor(key string, r, g, b, a int)
func (KeyValues) SetColor4(key string, color [4]int)
func (KeyValues) SetVector(key string, vec Vec3)

func (KeyValues) GetString(key string, value []char, maxlength int, defvalue string)
func (KeyValues) GetNum(key string, defvalue int) int
func (KeyValues) GetFloat(key string, defvalue float) float
func (KeyValues) GetColor(key string, r, g, b, a *int)
func (KeyValues) GetColor4(key string, color *[4]int)
func (KeyValues) GetUInt64(key string, value [2]int, defvalue [2]int)
func (KeyValues) GetVector(key string, vec *Vec3, defvalue Vec3)

func (KeyValues) JumpToKey(key string, create bool) bool
func (KeyValues) JumpToKeySymbol(id int) bool
func (KeyValues) GotoFirstSubKey(keyOnly bool) bool
func (KeyValues) GotoNextKey(keyOnly bool) bool
func (KeyValues) SavePosition()
func (KeyValues) GoBack() bool
func (KeyValues) DeleteKey(key string) bool
func (KeyValues) DeleteThis() int
func (KeyValues) Rewind()

func (KeyValues) GetSectionName(section []char, maxlength int) bool
func (KeyValues) SetSectionName(section string)
func (KeyValues) GetDataType(key string) KvDataTypes
func (KeyValues) SetEscapeSequences(useEscapes bool)
func (KeyValues) NodesInStack() int
func (KeyValues) FindKeyById(id int, name []char, maxlength int) bool
func (KeyValues) GetNameSymbol(key string, id *int) bool
func (KeyValues) GetSectionSymbol(id *int) bool

func KvSetString(kv KeyValues, key, value string)
func KvSetNum(kv KeyValues, key string, value int)
func KvSetUInt64(kv KeyValues, key string, value [2]int)
func KvSetFloat(kv KeyValues, key string, value float)
func KvSetColor(kv KeyValues, key string, r, g, b, a int)
func KvSetVector(kv KeyValues, key string, vec Vec3)

func KvGetString(kv KeyValues, key string, value []char, maxlength int, defvalue string)
func KvGetNum(kv KeyValues, key string, defvalue int) int
func KvGetFloat(kv KeyValues, key string, defvalue float) float
func KvGetColor(kv KeyValues, key string, r, g, b, a *int)
func KvGetUInt64(kv KeyValues, key string, value [2]int, defvalue [2]int)
func KvGetVector(kv KeyValues, key string, vec *Vec3, defvalue Vec3)

func KvJumpToKey(kv KeyValues, key string, create bool) bool
func KvJumpToKeySymbol(kv KeyValues, id int) bool
func KvGotoFirstSubKey(kv KeyValues, keyOnly bool) bool
func KvGotoNextKey(kv KeyValues, keyOnly bool) bool
func KvSavePosition(kv KeyValues)
func KvDeleteKey(kv KeyValues, key string) bool
func KvDeleteThis(kv KeyValues) int
func KvGoBack(kv KeyValues) bool
func KvRewind(kv KeyValues)

func KvGetSectionName(kv KeyValues, section []char, maxlength int) bool
func KvSetSectionName(kv KeyValues, section string)
func KvGetDataType(kv KeyValues, key string) KvDataTypes
func KeyValuesToFile(kv KeyValues, file string) bool
func FileToKeyValues(kv KeyValues, file string) bool
func StringToKeyValues(kv KeyValues, buffer, resourceName string) bool
func KvSetEscapeSequences(kv KeyValues, useEscapes bool)
func KvNodesInStack(kv KeyValues) int
func KvCopySubkeys(origin, dest KeyValues)
func KvFindKeyById(kv KeyValues, id int, name []char, maxlength int) bool
func KvGetNameSymbol(kv KeyValues, key string, id *int) bool
func KvGetSectionSymbol(kv KeyValues, id *int) bool