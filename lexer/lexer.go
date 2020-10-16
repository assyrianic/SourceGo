/**
 * lexer.go
 * 
 * Copyright 2020 Nirari Technologies.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package lexer

import (
	//"strings"
	//"unicode"
	"go/scanner"
	"go/token"
)

const (
	TokenNative = token.VAR + 1
	TokenTrue
	TokenFalse
	TokenLen
	
	/// types
	TokenInt        /// int    => int
	TokenFloat      /// float  => float
	TokenBool       /// bool   => float
	TokenVec        /// vec3   => float[3]
	TokenMap        /// map    => StringMap
	TokenArray      /// array  => ArrayList
	TokenHandle     /// obj    => Handle
)

var (
	fixed_tokens = map[string]token.Token {
		"native": TokenNative,
		"true": TokenTrue,
		"false": TokenFalse,
		"len": TokenLen,
		
		"int": TokenInt,
		"float": TokenFloat,
		"bool": TokenBool,
		"vec3": TokenVec,
		"map": TokenMap,
		"array": TokenArray,
		"obj": TokenHandle,
	}
)

type GPToken struct {
	l string
	t token.Token
	p token.Pos
}

/**
 * lexer.Tokenize
 * Takes a (preprocessed) source code and transforms it into an array of tokens
 */
func Tokenize(src []byte) []GPToken {
	tokens := make([]GPToken, 0)
	
	/// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet() /// positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) /// register input "file"
	s.Init(file, src, nil /** no error handler */, scanner.ScanComments)
	
	/// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		tokens = append(tokens, GPToken{ l: lit, t: tok, p: pos })
	}
	return tokens
}