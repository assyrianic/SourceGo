/**
 * sourcemod/events.go
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


type EventHookMode int
const (
	EventHookMode_Pre = EventHookMode(0)       //< Hook callback fired before event is fired */
	EventHookMode_Post                 //< Hook callback fired after event is fired */
	EventHookMode_PostNoCopy            //< Hook callback fired after event is fired, but event data won't be copied */
)


type (
	EventHook func(event Event, name string, dontBroadcast bool) Action
	Event struct {
		BroadcastDisabled bool
	}
)

func (Event) Fire(dontBroadcast bool)
func (Event) FireToClient(client Entity)
func (Event) Cancel()
func (Event) GetBool(key string, defValue bool) bool
func (Event) SetBool(key string, value bool)
func (Event) GetInt(key string, defValue int) int
func (Event) SetInt(key string, value int)
func (Event) GetFloat(key string, defValue float) float
func (Event) SetFloat(key string, value float)
func (Event) GetString(key string, value []char, maxlength int, defvalue string)
func (Event) SetString(key, value string)
func (Event) GetName(name []char, maxlength int)

func HookEvent(name string, callback EventHook, mode EventHookMode)
func HookEventEx(name string, callback EventHook, mode EventHookMode) bool
func UnhookEvent(name string, callback EventHook, mode EventHookMode)
func CreateEvent(name string, force bool) Event
func FireEvent(event Event, dontBroadcast bool)
func CancelCreatedEvent(event Event)
func GetEventBool(event Event, key string, defValue bool) bool
func SetEventBool(event Event, key string, value bool)
func GetEventInt(event Event, key string, defValue int) int
func SetEventInt(event Event, key string, value int)
func GetEventFloat(event Event, key string, defValue float) float
func SetEventFloat(event Event, key string, value float)
func GetEventString(event Event, key string, value []char, maxlength int, defvalue string)
func SetEventString(event Event, key, value string)
func GetEventName(event Event, name []char, maxlength int)
func SetEventBroadcast(event Event, dontBroadcast bool)