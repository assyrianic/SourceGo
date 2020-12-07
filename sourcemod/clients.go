/**
 * sourcemod/clients.go
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


type NetFlow int
const (
	NetFlow_Outgoing = NetFlow(0)   /**< Outgoing traffic */
	NetFlow_Incoming       /**< Incoming traffic */
	NetFlow_Both            /**< Both values added together */
)

type AuthIdType int
const (
	AuthId_Engine = AuthIdType(0)     /**< The game-specific auth string as returned from the engine */
	
	// The following are only available on games that support Steam authentication.
	AuthId_Steam2         /**< Steam2 rendered format, ex "STEAM_1:1:4153990" */
	AuthId_Steam3         /**< Steam3 rendered format, ex "[U:1:8307981]" */
	AuthId_SteamID64       /**< A SteamID64 (uint64) as a String, ex "76561197968573709" */
)


const (
	MAXPLAYERS = 65
	PLAYERS_SIZE = MAXPLAYERS + 1
	MAX_NAME_LENGTH = 128
	MAX_TF_PLAYERS = 36
)

var MaxClients int


func GetMaxHumanPlayers() int
func GetClientCount(inGameOnly bool) int
func GetClientName(client Entity, name []char, maxlen int) bool
func GetClientIP(client Entity, ip []char, maxlen int, remport bool) bool
func GetClientAuthId(client Entity, authtype AuthIdType, auth []char, maxlen int, validate bool) bool
func GetSteamAccountID(client Entity, validate bool) int
func GetClientUserId(client Entity) int
func IsClientConnected(client Entity) bool
func IsClientInGame(client Entity) bool
func IsClientInKickQueue(client Entity) bool
func IsClientAuthorized(client Entity) bool
func IsFakeClient(client Entity) bool
func IsClientSourceTV(client Entity) bool
func IsClientReplay(client Entity) bool
func IsClientObserver(client Entity) bool
func IsPlayerAlive(client Entity) bool
func GetClientInfo(client Entity, key string, value []char, maxlen int) bool
func GetClientTeam(client Entity) int
func SetUserAdmin(client Entity, id AdminId, temp bool)
func GetUserAdmin(client Entity) AdminId
func AddUserFlags(client Entity, flags ...AdminFlag)
func RemoveUserFlags(client Entity, flags ...AdminFlag)
func SetUserFlagBits(client, flags int)
func GetUserFlagBits(client Entity) int
func CanUserTarget(client, target int) bool
func RunAdminCacheChecks(client Entity) bool
func NotifyPostAdminCheck(client Entity)
func CreateFakeClient(name string) int
func SetFakeClientConVar(client Entity, cvar, value string)
func GetClientHealth(client Entity) int
func GetClientModel(client Entity, model []char, maxlen int)
func GetClientWeapon(client Entity, wep []char, maxlen int)
func GetClientMaxs(client Entity, vec *Vec3)
func GetClientMins(client Entity, vec *Vec3)
func GetClientAbsAngles(client Entity, vec *Vec3)
func GetClientAbsOrigin(client Entity, vec *Vec3)
func GetClientArmor(client Entity) int
func GetClientDeaths(client Entity) int
func GetClientFrags(client Entity) int
func GetClientDataRate(client Entity) int
func IsClientTimingOut(client Entity) bool
func GetClientTime(client Entity) float
func GetClientLatency(client Entity, flow NetFlow) float
func GetClientAvgLatency(client Entity, flow NetFlow) float
func GetClientAvgLoss(client Entity, flow NetFlow) float
func GetClientAvgChoke(client Entity, flow NetFlow) float
func GetClientAvgData(client Entity, flow NetFlow) float
func GetClientAvgPackets(client Entity, flow NetFlow) float
func GetClientOfUserId(userid int) int
func KickClient(client Entity, format string, args ...any)
func KickClientEx(client Entity, format string, args ...any)
func ChangeClientTeam(client, team int)
func GetClientSerial(client Entity) int
func GetClientFromSerial(cl_serial int) int