/**
 * sourcemod.go
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
	"sourcemod/core"
	"sourcemod/floats"
	"sourcemod/vector"
	"sourcemod/strings"
	"sourcemod/handles"
	"sourcemod/functions"
	"sourcemod/files"
	"sourcemod/logging"
	"sourcemod/timers"
	
	/** TODO: Finish making SourceGo iface for these includes:
		#include <admin>
		#include <keyvalues>
		#include <dbi>
		#include <lang>
		#include <sorting>
		#include <textparse>
		#include <clients>
		#include <console>
		#include <convars>
		#include <events>
		#include <bitbuffer>
		#include <protobuf>
		#include <usermessages>
		#include <menus>
		#include <halflife>
		#include <adt>
		#include <banning>
		#include <commandfilters>
		#include <nextmap>
		#include <commandline>
		#include <helpers>
		#include <entity>
		#include <entity_prop_stocks>
	*/
)


var MaxClients int

const (
	MAXPLAYERS = 65
	MAXENTS = 2048
)


func IsClientInGame(client int) bool
func PrintToServer(fmt string, args ...any)
func IsValidEntity(ent Entity) bool