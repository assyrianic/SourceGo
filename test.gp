package main


import (
	//"sourcemod" /// automatically inserted.
	"tf2_stocks"
	"sdkhooks"
	"arraylist"
	"fishy_is_rusty"
	".Impacts_suggestion"
)

func OnPluginStart() {
	for i:=MaxClients; i; i-- {
		if !IsValidClient(i) {
			continue
		}
		OnClientPutInServer(i)
	}
}

func IsValidClient(client int, replaycheck bool) bool {
	if !IsClientValid(client) || !IsClientInGame(client) || GetEntProp(client, Prop_Send, "m_bIsCoaching") {
		return false
	} else if replaycheck && (IsClientSourceTV(client) || IsClientReplay(client)) {
		return false
	} else {
		return true
	}
}

func GetOwner(ent int) int {
	if IsValidEntity(ent) {
		return GetEntPropEnt(ent, Prop_Send, "m_hOwnerEntity")
	} else {
		return -1
	}
}

func Func1(x, n int) int
func (n int) Receiver(c float)