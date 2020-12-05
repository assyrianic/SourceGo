/**
 * sourcemod/protobuf.go
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


const PB_FIELD_NOT_REPEATED = -1

type Protobuf Handle
func (Protobuf) ReadInt(field string, index int) int
func (Protobuf) ReadInt64(field string, value *[2]int, index int)
func (Protobuf) ReadFloat(field string, index int) float
func (Protobuf) ReadBool(field string, index int) bool
func (Protobuf) ReadString(field string, buffer []char, maxlen int, index int)
func (Protobuf) ReadColor(field string, buffer *[4]int, index int)
func (Protobuf) ReadAngle(field string, buffer *Vec3, index int)
func (Protobuf) ReadVector(field string, buffer *Vec3, index int)
func (Protobuf) ReadVector2D(field string, buffer *[2]float, index int)
func (Protobuf) GetRepeatedFieldCount(field string) int
func (Protobuf) HasField(field string) bool
func (Protobuf) SetInt(field string, value int, index int)
func (Protobuf) SetInt64(field string, value [2]int, index int)
func (Protobuf) SetFloat(field string, value float, index int)
func (Protobuf) SetBool(field string, value bool, index int)
func (Protobuf) SetString(field, value string, index int)
func (Protobuf) SetColor(field string, color [4]int, index int)
func (Protobuf) SetAngle(field string, vec Vec3, index int)
func (Protobuf) SetVector(field string, vec Vec3, index int)
func (Protobuf) SetVector2D(field string, vec [2]float, index int)
func (Protobuf) AddInt(field string, value int)
func (Protobuf) AddInt64(field string, value [2]int)
func (Protobuf) AddFloat(field string, value float)
func (Protobuf) AddBool(field string, value bool)
func (Protobuf) AddString(field, value string)
func (Protobuf) AddColor(field string, color [4]int)
func (Protobuf) AddAngle(field string, vec Vec3)
func (Protobuf) AddVector(field string, vec Vec3)
func (Protobuf) AddVector2D(field string, vec [2]float)
func (Protobuf) RemoveRepeatedFieldValue(field string, index int)
func (Protobuf) ReadMessage(field string) Protobuf
func (Protobuf) ReadRepeatedMessage(field string, index int) Protobuf
func (Protobuf) AddMessage(field string) Protobuf


func PbReadInt(pb Protobuf, field string, index int) int
func PbReadFloat(pb Protobuf, field string, index int) float
func PbReadBool(pb Protobuf, field string, index int) bool
func PbReadString(pb Protobuf, field string, buffer []char, maxlen int, index int)
func PbReadColor(pb Protobuf, field string, buffer *[4]int, index int)
func PbReadAngle(pb Protobuf, field string, buffer *Vec3, index int)
func PbReadVector(pb Protobuf, field string, buffer *Vec3, index int)
func PbReadVector2D(pb Protobuf, field string, buffer *[2]float, index int)
func PbGetRepeatedFieldCount(pb Protobuf, field string) int
func PbHasField(pb Protobuf, field string) bool
func PbSetInt(pb Protobuf, field string, value int, index int)
func PbSetFloat(pb Protobuf, field string, value float, index int)
func PbSetBool(pb Protobuf, field string, value bool, index int)
func PbSetString(field, value string, index int)
func PbSetColor(pb Protobuf, field string, color [4]int, index int)
func PbSetAngle(pb Protobuf, field string, vec Vec3, index int)
func PbSetVector(pb Protobuf, field string, vec Vec3, index int)
func PbSetVector2D(pb Protobuf, field string, vec [2]float, index int)
func PbAddInt(pb Protobuf, field string, value int)
func PbAddFloat(pb Protobuf, field string, value float)
func PbAddBool(pb Protobuf, field string, value bool)
func PbAddString(field, value string)
func PbAddColor(pb Protobuf, field string, color [4]int)
func PbAddAngle(pb Protobuf, field string, vec Vec3)
func PbAddVector(pb Protobuf, field string, vec Vec3)
func PbAddVector2D(pb Protobuf, field string, vec [2]float)
func PbRemoveRepeatedFieldValue(pb Protobuf, field string, index int)
func PbReadMessage(pb Protobuf, field string) Protobuf
func PbReadRepeatedMessage(pb Protobuf, field string, index int) Protobuf
func PbAddMessage(pb Protobuf, field string) Protobuf