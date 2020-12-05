/**
 * sourcemod/usermessages.go
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


type UserMsg int
const INVALID_MESSAGE_ID = UserMsg(-1)

type UserMessageType int
const (
	UM_BitBuf = UserMessageType(0)
	UM_Protobuf
)

const (
	USERMSG_RELIABLE =        (1<<2)    /**< Message will be set on the reliable stream */
	USERMSG_INITMSG =         (1<<3)    /**< Message will be considered to be an initmsg */
	USERMSG_BLOCKHOOKS =      (1<<7)    /**< Prevents the message from triggering SourceMod and Metamod hooks */
)


func GetUserMessageType() UserMessageType
func UserMessageToProtobuf(msg Handle) Protobuf
func UserMessageToBfWrite(msg Handle) BfWrite
func UserMessageToBfRead(msg Handle) BfRead

func GetUserMessageId(msg string) UserMsg
func GetUserMessageName(msg_id UserMsg, msg []char, maxlength int) bool
func StartMessage(msgname string, clients []int, numClients, flags int) Handle
func StartMessageEx(msg UserMsg, clients []int, numClients, flags int) Handle
func EndMessage()


type (
	MsgHook func(msg_id UserMsg, msg Handle, players []int, playersNum int, reliable, init bool) Action
	MsgPostHook func(msg_id UserMsg, sent bool)
)

func HookUserMessage(msg_id UserMsg, hook MsgHook, intercept bool, post MsgPostHook)
func UnhookUserMessage(msg_id UserMsg, hook MsgHook, intercept bool)
func StartMessageAll(msgname string, flags int) Handle
func StartMessageOne(msgname string, client, flags int) Handle