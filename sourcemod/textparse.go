/**
 * sourcemod/textparse.go
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


type SMCResult int
const (
	SMCParse_Continue = SMCResult(0)          /**< Continue parsing */
	SMCParse_Halt              /**< Stop parsing here */
	SMCParse_HaltFail           /**< Stop parsing and return failure */
)


type SMCError int
const (
	SMCError_Okay = SMCError(0)          /**< No error */
	SMCError_StreamOpen        /**< Stream failed to open */
	SMCError_StreamError       /**< The stream died... somehow */
	SMCError_Custom            /**< A custom handler threw an error */
	SMCError_InvalidSection1   /**< A section was declared without quotes, and had extra tokens */
	SMCError_InvalidSection2   /**< A section was declared without any header */
	SMCError_InvalidSection3   /**< A section ending was declared with too many unknown tokens */
	SMCError_InvalidSection4   /**< A section ending has no matching beginning */
	SMCError_InvalidSection5   /**< A section beginning has no matching ending */
	SMCError_InvalidTokens     /**< There were too many unidentifiable strings on one line */
	SMCError_TokenOverflow     /**< The token buffer overflowed */
	SMCError_InvalidProperty1   /**< A property was declared outside of any section */
)


type (
	SMC_ParseStart func(smc SMCParser)
	SMC_NewSection func(smc SMCParser, name string, opt_quotes bool) SMCResult
	SMC_KeyValue func (smc SMCParser, key, value string, key_quotes, value_quotes bool) SMCResult
	SMC_EndSection func(smc SMCParser) SMCResult
	SMC_ParseEnd func(smc SMCParser, halted, failed bool)
	SMC_RawLine func(smc SMCParser, line string, lineno int) SMCResult
	
	
	/// new SMCParser();
	SMCParser struct {
		OnStart SMC_ParseStart
		OnEnd SMC_ParseEnd
		OnEnterSection SMC_NewSection
		OnLeaveSection SMC_EndSection
		OnKeyValue SMC_KeyValue
		OnRawLine SMC_RawLine
	}
)


func (SMCParser) ParseFile(file string, line, col *int) SMCError
func (SMCParser) GetErrorString(err SMCError, buffer []char, buf_max int)


func SMC_CreateParser() SMCParser
func SMC_ParseFile(smc SMCParser, file string, line, col *int) SMCError
func SMC_GetErrorString(smc SMCParser, err SMCError, buffer []char, buf_max int)
func SMC_SetParseStart(smc SMCParser, fn SMC_ParseStart)
func SMC_SetParseEnd(smc SMCParser, fn SMC_ParseEnd)
func SMC_SetReaders(smc SMCParser, ns SMC_NewSection, kv SMC_KeyValue, es SMC_EndSection)
func SMC_SetRawLine(smc SMCParser, fn SMC_RawLine)