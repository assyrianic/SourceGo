/**
 * sourcemod/vector.go
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

import (
	"sourcemod/types"
)

type Vec3 = [3]float

func GetVectorLength(vec Vec3, squared bool) float
func GetVectorDistance(vec1, vec2 Vec3, squared bool) float
func GetVectorDotProduct(vec1, vec2 Vec3) float
func GetVectorCrossProduct(vec1, vec2 Vec3, result *Vec3)
func NormalizeVector(vec Vec3, result *Vec3) float
func GetAngleVectors(angle Vec3, fwd, right, up *Vec3)
func GetVectorAngles(vec Vec3, angle *Vec3)
func GetVectorVectors(vec Vec3, right, up *Vec3)
func AddVectors(vec1, vec2 Vec3, result *Vec3)
func SubtractVectors(vec1, vec2 Vec3, result *Vec3)
func ScaleVector(vec *Vec3, scale float)
func NegateVector(vec *Vec3)
func MakeVectorFromPoints(pt1, pt2 Vec3, output *Vec3)