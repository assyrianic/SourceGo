/**
 * sourcemod/convars.go
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


type ConVarBounds int
const (
	ConVarBound_Upper = ConVarBounds(0)
	ConVarBound_Lower
)


type ConVarQueryResult int
const (
	ConVarQuery_Okay = ConVarQueryResult(0)               //< Retrieval of client convar value was successful. */
	ConVarQuery_NotFound               //< Client convar was not found. */
	ConVarQuery_NotValid               //< A console command with the same name was found, but there is no convar. */
	ConVarQuery_Protected               //< Client convar was found, but it is protected. The server cannot retrieve its value. */
)


type ConVarChanged func(convar ConVar, oldValue, newValue string)

func CreateConVar(name, defaultValue, description string, flags int, hasMin bool, min float, hasMax bool, max float) ConVar
func FindConVar(name string) ConVar

type ConVar struct {
	BoolValue bool
	IntValue, Flags int
	FloatValue float
}
func (ConVar) SetBool(value, replicate, notify bool)
func (ConVar) SetInt(value int, replicate, notify bool)
func (ConVar) SetFloat(value float, replicate, notify bool)
func (ConVar) GetString(value []char, maxlength int)
func (ConVar) SetString(value string, replicate, notify bool)
func (ConVar) RestoreDefault(replicate, notify bool)
func (ConVar) GetDefault(value []char, maxlength int) int
func (ConVar) GetBounds(bounds_type ConVarBounds, value *float) bool
func (ConVar) SetBounds(bounds_type ConVarBounds, set bool, value float)
func (ConVar) GetName(name []char, maxlength int)
func (ConVar) ReplicateToClient(client Entity, value string) bool
func (ConVar) AddChangeHook(callback ConVarChanged)
func (ConVar) RemoveChangeHook(callback ConVarChanged)

func HookConVarChange(cvar ConVar, callback ConVarChanged)
func UnhookConVarChange(cvar ConVar, callback ConVarChanged)

func GetConVarBool(cvar ConVar) bool
func GetConVarInt(cvar ConVar) int
func GetConVarFloat(cvar ConVar) float
func GetConVarString(cvar ConVar, value []char, maxlength int)

func SetConVarBool(cvar ConVar, value, replicate, notify bool)
func SetConVarInt(cvar ConVar, value int, replicate, notify bool)
func SetConVarFloat(cvar ConVar, value float, replicate, notify bool)
func SetConVarString(cvar ConVar, value string, replicate, notify bool)

func ResetConVar(cvar ConVar, replicate, notify bool)
func GetConVarDefault(cvar ConVar, value []char, maxlength int) int
func GetConVarFlags(cvar ConVar) int
func SetConVarFlags(cvar ConVar, flags int)
func GetConVarBounds(cvar ConVar, bounds_type ConVarBounds, value *float) bool
func SetConVarBounds(cvar ConVar, bounds_type ConVarBounds, set bool, value float)
func GetConVarName(cvar ConVar, name []char, maxlength int)
func SendConVarValue(client Entity, cvar ConVar, value string) bool

type ConVarQueryFinished func(cookie QueryCookie, client Entity, result ConVarQueryResult, cvarName, cvarValue string, value any)
func QueryClientConVar(client Entity, cvarName string, callback ConVarQueryFinished, value any) QueryCookie
func IsValidConVarChar(c int) bool