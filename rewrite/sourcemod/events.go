/*
 * sourcemod/events.go
 * 
 * Copyright 2022 Nirari Technologies, Alliedmodders LLC.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package sm


type EventHookMode int
const (
	EventHookMode_Pre        = EventHookMode(0) /* Hook callback fired before event is fired */
	EventHookMode_Post       = EventHookMode(1) /* Hook callback fired after event is fired */
	EventHookMode_PostNoCopy = EventHookMode(2) /* Hook callback fired after event is fired, but event data won't be copied */
)


type (
	EventHook func(event Event, name []byte, dontBroadcast bool) Action
	Event struct {
		BroadcastDisabled bool
	}
)

func (*Event) Fire(dontBroadcast bool)
func (*Event) FireToClient(client Entity)
func (*Event) Cancel()
func (*Event) GetBool(key []byte, defValue bool) bool
func (*Event) SetBool(key []byte, value bool)
func (*Event) GetInt(key []byte, defValue int) int
func (*Event) SetInt(key []byte, value int)
func (*Event) GetFloat(key []byte, defValue float32) float32
func (*Event) SetFloat(key []byte, value float32)
func (*Event) GetString(key []byte, value []byte, maxlength int, defvalue []byte)
func (*Event) SetString(key, value []byte)
func (*Event) GetName(name []byte, maxlength int)

func HookEvent(name []byte, callback EventHook, mode EventHookMode)
func HookEventEx(name []byte, callback EventHook, mode EventHookMode) bool
func UnhookEvent(name []byte, callback EventHook, mode EventHookMode)
func CreateEvent(name []byte, force bool) Event
