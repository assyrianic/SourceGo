/**
 * sourcemod/functions.go
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

const SP_PARAMFLAG_BYREF = 1

type Function __function__

type ParamType int
const (
	Param_Any           = ParamType(0)                            /**< Any data type can be pushed */
	Param_Cell          = ParamType(1<<1)                       /**< Only basic cells can be pushed */
	Param_Float         = ParamType(2<<1)                       /**< Only floats can be pushed */
	Param_String        = ParamType((3<<1)|SP_PARAMFLAG_BYREF)    /**< Only strings can be pushed */
	Param_Array         = ParamType((4<<1)|SP_PARAMFLAG_BYREF)    /**< Only arrays can be pushed */
	Param_VarArgs       = ParamType(5<<1)                       /**< Same as "..." in plugins, anything can be pushed, but it will always be byref */
	Param_CellByRef     = ParamType((1<<1)|SP_PARAMFLAG_BYREF)    /**< Only a cell by reference can be pushed */
	Param_FloatByRef    = ParamType((2<<1)|SP_PARAMFLAG_BYREF)     /**< Only a float by reference can be pushed */
)

type ExecType int
const (
	ET_Ignore   = 0    /**< Ignore all return values, return 0 */
	ET_Single    /**< Only return the last exec, ignore all others */
	ET_Event    /**< Acts as an event with the Actions defined in core.inc, no mid-Stops allowed, returns highest */
	ET_Hook     /**< Acts as a hook with the Actions defined in core.inc, mid-Stops allowed, returns highest */
)

const (
	SM_PARAM_COPYBACK      = (1<<0)      /**< Copy an array/reference back after call */
	SM_PARAM_STRING_UTF8   = (1<<0)      /**< String should be UTF-8 handled */
	SM_PARAM_STRING_COPY   = (1<<1)      /**< String should be copied into the plugin */
	SM_PARAM_STRING_BINARY = (1<<2)      /**< Treat the string as a binary string */
)

const (
	SP_ERROR_NONE                =   0   /**< No error occurred */
	SP_ERROR_FILE_FORMAT         =   1   /**< File format unrecognized */
	SP_ERROR_DECOMPRESSOR        =   2   /**< A decompressor was not found */
	SP_ERROR_HEAPLOW             =   3   /**< Not enough space left on the heap */
	SP_ERROR_PARAM               =   4   /**< Invalid parameter or parameter type */
	SP_ERROR_INVALID_ADDRESS     =   5   /**< A memory address was not valid */
	SP_ERROR_NOT_FOUND           =   6   /**< The object in question was not found */
	SP_ERROR_INDEX               =   7   /**< Invalid index parameter */
	SP_ERROR_STACKLOW            =   8   /**< Not enough space left on the stack */
	SP_ERROR_NOTDEBUGGING        =   9   /**< Debug mode was not on or debug section not found */
	SP_ERROR_INVALID_INSTRUCTION =   10  /**< Invalid instruction was encountered */
	SP_ERROR_MEMACCESS           =   11  /**< Invalid memory access */
	SP_ERROR_STACKMIN            =   12  /**< Stack went beyond its minimum value */
	SP_ERROR_HEAPMIN             =   13  /**< Heap went beyond its minimum value */
	SP_ERROR_DIVIDE_BY_ZERO      =   14  /**< Division by zero */
	SP_ERROR_ARRAY_BOUNDS        =   15  /**< Array index is out of bounds */
	SP_ERROR_INSTRUCTION_PARAM   =   16  /**< Instruction had an invalid parameter */
	SP_ERROR_STACKLEAK           =   17  /**< A native leaked an item on the stack */
	SP_ERROR_HEAPLEAK            =   18  /**< A native leaked an item on the heap */
	SP_ERROR_ARRAY_TOO_BIG       =   19  /**< A dynamic array is too big */
	SP_ERROR_TRACKER_BOUNDS      =   20  /**< Tracker stack is out of bounds */
	SP_ERROR_INVALID_NATIVE      =   21  /**< Native was pending or invalid */
	SP_ERROR_PARAMS_MAX          =   22  /**< Maximum number of parameters reached */
	SP_ERROR_NATIVE              =   23  /**< Error originates from a native */
	SP_ERROR_NOT_RUNNABLE        =   24  /**< Function or plugin is not runnable */
	SP_ERROR_ABORTED             =   25  /**< Function call was aborted */
)


type Forward interface{}


type GlobalForward struct {
	FuncCount int
}

/// make(GlobalForward, name, exectype, param_types...)
//func GlobalForward(name string, exectype ExecType, param_types ...ParamType) GlobalForward
func CreateGlobalForward(name string, exectype ExecType, param_types ...ParamType) GlobalForward


type PrivateForward struct {
	FuncCount int
}
/// make(PrivateForward, exectype, param_types...)
//func PrivateForward(exectype ExecType, param_types ...ParamType) PrivateForward
func CreateForward(exectype ExecType, param_types ...ParamType) PrivateForward

func (PrivateForward) AddFunction(plugin Handle, fn Function) bool
func (PrivateForward) RemoveFunction(plugin Handle, fn Function) bool
func (PrivateForward) RemoveAllFunctions(plugin Handle) int

func GetForwardFunctionCount(fwd Forward) int
func AddToForward(fwd PrivateForward, plugin Handle, fn Function) bool
func RemoveFromForward(fwd PrivateForward, plugin Handle, fn Function) bool
func RemoveAllFromForward(fwd PrivateForward, plugin Handle) int

func Call_StartForward(fwd Forward)
func Call_StartFunction(plugin Handle, fn Function)
func Call_PushCell(value any)
func Call_PushCellRef(value *any)
func Call_PushFloat(value float)
func Call_PushFloatRef(value *float)
func Call_PushArray(value []any, size int)
func Call_PushArrayEx(value []any, size, cpflags int)
func Call_PushNullVector()
func Call_PushString(value string)
func Call_PushStringEx(value string, length, szflags, cpflags int)
func Call_PushNullString()
func Call_Finish(result *any) int
func Call_Cancel()


type NativeCall func(plugin Handle, numParams int) any

func GetFunctionByName(plugin Handle, name string) Function
func CreateNative(name string, fn NativeCall)
func ThrowNativeError(err int, fmt string, args ...any) int
func GetNativeStringLength(param int, length *int) int
func GetNativeString(param int, buffer string, maxlength int, bytes *int) int
func SetNativeString(param int, source string, maxlength int, utf8 bool, bytes *int) int
func GetNativeCell(param int) any
func GetNativeFunction(param int) Function
func GetNativeCellRef(param int) any
func SetNativeCellRef(param int, value any)
func GetNativeArray(param int, local []any, size int) int
func SetNativeArray(param int, local []any, size int) int
func IsNativeParamNullVector(param int) bool
func IsNativeParamNullString(param int) bool
func FormatNativeString(out_param, fmt_param, vararg_param, out_len int, written *int, out_string, fmt_string string) int


type RequestFrameCallback func(data any)
func RequestFrame(framecb RequestFrameCallback, data any)