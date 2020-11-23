/**
 * sourcemod/floats.go
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

type float = float64


const FLOAT_PI = 3.1415926535897932384626433832795


func FloatFraction(value float) float
func RoundToZero(value float) int
func RoundToCeil(value float) int
func RoundToFloor(value float) int
func RoundToNearest(value float) int
func FloatCompare(fOne, fTwo float) int

func SquareRoot(val float) float
func Pow(val, exp float) float
func Exponential(val float) float
func Logarithm(val float, base float) float
func Sine(val float) float
func Cosine(val float) float
func Tangent(val float) float
func FloatAbs(val float) float
func ArcTangent(val float) float
func ArcCosine(val float) float
func ArcSine(val float) float
func ArcTangent2(x,y float) float

func RoundFloat(val float) int

func DegToRad(angle float) float
func RadToDeg(angle float) float

func GetURandomInt() int
func GetURandomFloat() float
func SetURandomSeed(seeds []int, numSeeds int)

func SetURandomSeedSimple(seed int)