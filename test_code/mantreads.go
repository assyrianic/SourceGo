package main

import (
	"sourcemod"
	"sdkhooks"
)

func main() {
	for i := 1; i<=MaxClients; i++ {
		if IsClientInGame(i) {
			OnClientPutInServer(i)
		}
	}
}

func OnClientPutInServer(client Entity) {
	SDKHook(client, SDKHook_OnTakeDamage, func (victim int, attacker, inflictor *int, damage *float, damagetype, weapon *int, damageForce, damagePosition *Vec3, damagecustom int) Action {
		if IsValidEntity(*weapon) && GetEntProp(*weapon, Prop_Send, "m_iItemDefinitionIndex")==444 {
			*damage *= 5.0
			return Plugin_Changed
		}
		return Plugin_Continue
	})
}
