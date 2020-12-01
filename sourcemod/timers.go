/**
 * sourcemod/timers.go
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


const (
	TIMER_REPEAT            = (1<<0)      /**< Timer will repeat until it returns Plugin_Stop */
	TIMER_FLAG_NO_MAPCHANGE = (1<<1)      /**< Timer will not carry over mapchanges */
	TIMER_HNDL_CLOSE        = (1<<9)      /**< Deprecated define, replaced by below */
	TIMER_DATA_HNDL_CLOSE   = (1<<9)      /**< Timer will automatically call CloseHandle() on its data when finished */
)


type (
	Timer = Handle
	TimerFunc func(timer Timer, data any) Action
)

func CreateTimer(interval float, tfn TimerFunc, data any, flags int) Timer
func CreateDataTimer(interval float, tfn TimerFunc, datapack *any, flags int) Timer

func KillTimer(timer Timer, autoClose bool)
func TriggerTimer(timer Timer, reset bool)
func GetTickedTime() float
func GetMapTimeLeft(timeleft *int) bool
func GetMapTimeLimit(time *int) bool
func ExtendMapTimeLimit(time int) bool
func GetTickInterval() float
func IsServerProcessing() bool