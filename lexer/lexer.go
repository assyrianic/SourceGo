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
)

type TokenVal    int32

const (
	TokenInvalid    = TokenVal(iota)
	TokenIdent      /// id
	TokenIntLit     /// 0x123
	TokenFloatLit   /// 1.0
	TokenStringLit  /// "string"
	TokenCharLit    /// 'char'
	
	/// keywords.
	TokenBreak
	TokenContinue
	TokenIf
	TokenElse
	TokenFor
	TokenReturn
	TokenVar
	TokenNative
	TokenFunc
	TokenConst
	TokenSwitch
	TokenCase
	TokenDefault
	TokenTrue
	TokenFalse
	TokenLen
	
	/// types
	TokenInt        /// int    => int
	TokenEntity     /// Entity => int 
	TokenFloat      /// float  => float
	TokenBool       /// bool   => float
	TokenVec        /// vec3   => float[3]
	TokenMap        /// map    => StringMap
	TokenArray      /// array  => ArrayList
	TokenHandle     /// obj    => Handle
	
	/// operators and delims.
	TokenLeftParen  /// (
	TokenRiteParen  /// )
	TokenLeftSq     /// [
	TokenRiteSq     /// ]
	TokenLeftCurly  /// {
	TokenRiteCurly  /// }
)

func is_keyword(t TokenVal) bool {
	return TokenBreak <= t && t <= TokenHandle
}

var (
	fixed_tokens = map[string]TokenVal {
		"break": TokenBreak,
		"cont": TokenContinue,
		"if": TokenIf,
		"else": TokenElse,
		"for": TokenFor,
		"return": TokenReturn,
		"var": TokenVar,
		"native": TokenNative,
		"func": TokenFunc,
		"const": TokenConst,
		"switch": TokenSwitch,
		"case": TokenCase,
		"default": TokenDefault,
		"true": TokenTrue,
		"false": TokenFalse,
		"len": TokenLen,
		
		"int": TokenInt,
		"Entity": TokenEntity,
		"float": TokenFloat,
		"bool": TokenBool,
		"vec3": TokenVec,
		"map": TokenMap,
		"array": TokenArray,
		"obj": TokenHandle,
		
		"(": TokenLeftParen,
		")": TokenRiteParen,
		"[": TokenLeftSq,
		"]": TokenRiteSq,
		"{": TokenLeftCurly,
		"}": TokenRiteCurly,
	}
)

type Token struct {
	lexeme string
	tag    TokenVal
}


func is_whitespace(c byte) bool {
	return( c==' ' || c=='\t' || c=='\r' || c=='\v' || c=='\f' || c=='\n' )
}

func is_alphabetic(c byte) bool {
	return( (c>='a' && c<='z') || (c>='A' && c<='Z') || c=='_' );
}

func is_possible_id(c byte) bool {
	return( is_alphabetic(c) || (c>='0' && c<='9') );
}

func is_decimal(c byte) bool {
	return( c>='0' && c<='9' );
}

/**
 * lexer.Preprocess
 * Takes the source code as a byte array, removes comments, and returns the processed array.
 */
func Preprocess(src []byte) []byte {
	i := 0
	for i<len(src) {
		c := src[i]
		if is_whitespace(c) {
			i++
			continue
		} else if string(src[i:i+2])=="/*" {
			/// multi-line comment.
			for string(src[i:i+2]) != "*/" {
				src[i] = ' '
				i++
			}
			src[i] = ' '
			src[i+1] = ' '
		} else if string(src[i:i+2])=="//" {
			/// single-line comment.
			for src[i] != '\n' {
				src[i] = ' '
				i++
			}
		} else if src[i]=='"' || src[i]=='\'' || src[i]=='`' {
			/// carefully ignore strings.
			q := src[i]
			i++
			for src[i] != q {
				if src[i]=='\\' {
					i += 2
				} else {
					i++
				}
			}
		}
		i++
	}
	return src
}

/**
 * lexer.Tokenize
 * Takes a (preprocessed) source code and transforms it into an array of tokens
 */
func Tokenize(src []byte) []Token {
	tokens := make([]Token, 0)
	i := 0
	for i<len(src) {
		c := src[i]
		if is_alphabetic(c) {
			token := Token{}
			for is_possible_id(c) {
				token.lexeme += string(c)
				i++
				c = src[i]
			}
			if val, ok := fixed_tokens[token.lexeme]; ok {
				token.tag = val
			} else {
				token.tag = TokenIdent
			}
			tokens = append(tokens, token)
		} else if is_decimal(c) || (c=='.' && is_decimal(src[i+1])) {
			
		}
		i++
	}
	return tokens
}

func lex_number(src []byte) (bool, is_float bool) {
	
}