/**
 * sourcemod/banning.go
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
	BANFLAG_AUTO =        (1<<0)  /**< Auto-detects whether to ban by steamid or IP */
	BANFLAG_IP =          (1<<1)  /**< Always ban by IP address */
	BANFLAG_AUTHID =      (1<<2)  /**< Always ban by authstring (for BanIdentity) if possible */
	BANFLAG_NOKICK =      (1<<3)  /**< Does not kick the client */
)


func BanClient(client, time, flags int, reason, kick_message, command string, source any) bool
func BanIdentity(identity string, time, flags int, reason, command string, source any) bool
func RemoveBan(identity string, flags int, command string, source any) bool