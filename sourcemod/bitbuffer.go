/**
 * sourcemod/bitbuffer.go
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
	BfWrite Handle
	BfRead struct {
		BytesLeft int
	}
)
func (BfWrite) WriteBool(bit bool)
func (BfWrite) WriteByte(byt int)
func (BfWrite) WriteChar(chr int)
func (BfWrite) WriteShort(shrt int)
func (BfWrite) WriteWord(w int)
func (BfWrite) WriteNum(num int)
func (BfWrite) WriteFloat(num float)
func (BfWrite) WriteString(s string)
func (BfWrite) WriteEntity(ent Entity)
func (BfWrite) WriteAngle(angle float, numBits int)
func (BfWrite) WriteCoord(coord float)
func (BfWrite) WriteVecCoord(vec Vec3)
func (BfWrite) WriteVecNormal(vec Vec3)
func (BfWrite) WriteAngles(vec Vec3)

func (BfRead) ReadBool() bool
func (BfRead) ReadByte() int
func (BfRead) ReadChar() int
func (BfRead) ReadShort() int
func (BfRead) ReadWord() int
func (BfRead) ReadNum() int
func (BfRead) ReadFloat() float
func (BfRead) ReadString(buffer []char, maxlength int, line bool) int
func (BfRead) ReadEntity() Entity
func (BfRead) ReadAngle(numBits int) float
func (BfRead) ReadCoord() float
func (BfRead) ReadVecCoord(vec *Vec3)
func (BfRead) ReadVecNormal(vec *Vec3)
func (BfRead) ReadAngles(vec *Vec3)


func BfWriteBool(bf BfWrite, bit bool)
func BfWriteByte(bf BfWrite, byt int)
func BfWriteChar(bf BfWrite, chr int)
func BfWriteShort(bf BfWrite, num int)
func BfWriteWord(bf BfWrite, num int)
func BfWriteNum(bf BfWrite, num int)
func BfWriteFloat(bf BfWrite, num float)
func BfWriteString(bf BfWrite, str string)
func BfWriteEntity(bf BfWrite, ent Entity)
func BfWriteAngle(bf BfWrite, angle float, numBits int)
func BfWriteCoord(bf BfWrite, coord float)
func BfWriteVecCoord(bf BfWrite, vec Vec3)
func BfWriteVecNormal(bf BfWrite, vec Vec3)
func BfWriteAngles(bf BfWrite, vec Vec3)

func BfReadBool(bf BfRead) bool
func BfReadByte(bf BfRead) int
func BfReadChar(bf BfRead) int
func BfReadShort(bf BfRead) int
func BfReadWord(bf BfRead) int
func BfReadNum(bf BfRead) int
func BfReadFloat(bf BfRead) float
func BfReadString(bf BfRead, buffer []char, maxlength int, line bool) int
func BfReadEntity(bf BfRead) int
func BfReadAngle(bf BfRead, numBits int) float
func BfReadCoord(bf BfRead) float
func BfReadVecCoord(bf BfRead, vec *Vec3)
func BfReadVecNormal(bf BfRead, vec *Vec3)
func BfReadAngles(bf BfRead, vec *Vec3)
func BfGetNumBytesLeft(bf BfRead) int