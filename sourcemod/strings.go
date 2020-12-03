/**
 * sourcemod/strings.go
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

func strlen(str string) int
func StrContains(str, substr string, caseSensitive bool) int
func strcmp(str1, str2 string, caseSensitive bool) int
func strncmp(str1, str2 string, chars int, caseSensitive bool) int
func StrCompare(str1, str2 string, caseSensitive bool) int
func StrEqual(str1, str2 string, caseSensitive bool) bool
func strcopy(dest []char, destLen int, source string) int
func StrCopy(dest []char, destLen int, source string) int
func Format(buffer []char, maxlength int, format string, args ...any) int
func FormatEx(buffer []char, maxlength int, format string, args ...any) int
func VFormat(buffer []char, maxlength int, format string, varpos int) int
func StringToInt(str string, nBase int) int
func StringToIntEx(str string, result *int, nBase int) int
func IntToString(num int, str string, maxlength int) int
func StringToFloat(str string) float
func StringToFloatEx(str string, result *float) int
func FloatToString(num float, str []char, maxlength int) int
func BreakString(source []char, arg string, argLen int) int
func StrBreak(source []char, arg string, argLen int) int
func TrimString(str []char) int
func SplitString(source []char, split, part string, partLen int) int
func ReplaceString(text []char, maxlength int, search, replace string, caseSensitive bool) int
func ReplaceStringEx(text []char, maxlength int, search, replace string, searchLen, replaceLen int, caseSensitive bool) int

func GetCharBytes(source string) int
func IsCharAlpha(chr int) bool
func IsCharNumeric(chr int) bool
func IsCharSpace(chr int) bool
func IsCharMB(chr int) int
func IsCharUpper(chr int) bool
func IsCharLower(chr int) bool
func StripQuotes(text string) bool
func CharToUpper(chr int) char
func CharToLower(chr int) char
func FindCharInString(str string, c char, reverse bool) int
func StrCat(buffer string, maxlength int, source string) int
func ExplodeString(text []char, split string, buffers []string, maxStrings, maxStringLength int, copyRemainder bool) int
func ImplodeStrings(strings []string, numStrings int, join string, buffer []char, maxLength int) int