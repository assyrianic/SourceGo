/**
 * clientprefs.go
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

import "sourcemod"

type (
	CookieAccess int
	CookieMenu int
	CookieMenuAction int
	
	CookieMenuHandler func(client int, action CookieMenuAction, info any, buffer []char, maxlen int)
	
	Cookie struct {
		AccessLevel CookieAccess
	}
)

const (
	CookieAccess_Public            /**< Visible and Changeable by users */
	CookieAccess_Protected         /**< Read only to users */
	CookieAccess_Private           /**< Completely hidden cookie */
	
	CookieMenu_YesNo           /**< Yes/No menu with "yes"/"no" results saved into the cookie */
	CookieMenu_YesNo_Int       /**< Yes/No menu with 1/0 saved into the cookie */
	CookieMenu_OnOff           /**< On/Off menu with "on"/"off" results saved into the cookie */
	CookieMenu_OnOff_Int       /**< On/Off menu with 1/0 saved into the cookie */
	
	CookieMenuAction_DisplayOption = 0
	CookieMenuAction_SelectOption  = 1
)


func RegClientCookie(name, description string, access CookieAccess) Cookie
func FindClientCookie(name string) Cookie

func (Cookie) Set(client int, value string)
func (Cookie) Get(client int, buffer []char, maxlen int)
func (Cookie) SetByAuthId(authID, value string)
func (Cookie) SetPrefabMenu(menutype CookieMenu, display string, handler CookieMenuHandler, info any)
func (Cookie) GetClientTime(client int) int

func AreClientCookiesCached(client int) bool
func ShowCookieMenu(client int)