package main

const (
	a, b = "A", MAXPLAYERS
	c = a
	d string = "D"
	e = "e1"
	f = 1.00
	
)

var (
	myself = Plugin{
		name: "SrcGo Plugin",
		author: "Nergal",
		description: "Plugin made into SP from SrcGo.",
		version: "1.0a",
		url: "https://github.com/assyrianic/Go2SourcePawn",
	}
	t = [...]string{
		"kek",
		"foo",
		"bar",
		"bazz",
	}
	
	a int
	b Handle = nil
	c = 50.0
)

func NBC() Entity

func main() {
	var pf func() Entity = NBC
	client := pf()
}

type (
	Point struct{ x, y float }
	Points struct { p [3*99]Point }
	polar [3*5]Points
	QAngle Vec3
	Name [64]char
	Color = [4]int
)

type PlayerInfo struct {
	K polar
	P QAngle
	M [3][10]int
}

type Kektus func(i, x Vec3, b string, blocks *Name, KC *int) Handle
type EventFunc func(event Event, name string, dontBroadcast bool) Action

func GetNil() Handle {
	return nil
}

func OnlyScoutsLeft(team int) bool {
	for i:=MaxClients; i>0; i-- {
		if !IsValidClient(i) || !IsPlayerAlive(i) {
			continue;
		}
		if GetClientTeam(i) == team && TF2_GetPlayerClass(i) != TFClass_Scout {
			return false
		}
	}
	return true
}
func SpawnSmallHealthPackAt(client, ownerteam int) {
	if !IsValidClient(client) || !IsPlayerAlive(client) {
		return
	}
	var healthpack int = CreateEntityByName("item_healthkit_small")
	if IsValidEntity(healthpack) {
		var pos Vec3
		GetClientAbsOrigin(client, pos)
		pos[2] += 20.0
		
		/// for safety, though it normally doesn't respawn
		DispatchKeyValue(healthpack, "OnPlayerTouch", "!self,Kill,,0,-1")
		DispatchSpawn(healthpack)
		SetEntProp(healthpack, Prop_Send, "m_iTeamNum", ownerteam, 4)
		SetEntityMoveType(healthpack, MOVETYPE_VPHYSICS)
		var vel Vec3
		vel[0], vel[1], vel[2] = float(GetRandomInt(-10, 10)), float(GetRandomInt(-10, 10)), 50.0
		TeleportEntity(healthpack, pos, NULL_VECTOR, vel)
	}
}

func FindPlayerBack(client Entity, indices []int, len int) int {
	if len <= 0 {
		return -1
	}
	var numwearables int = TF2_GetNumWearables(client)
	for i:=0; i<numwearables; i++ {
		var wearable int = TF2_GetWearable(client, i)
		if wearable > 0 && !GetEntProp(wearable, Prop_Send, "m_bDisguiseWearable") {
			var idx int = GetItemIndex(wearable)
			for u:=0; u<len; u++ {
				if idx==indices[u] {
					return wearable
				}
			}
		}
	}
	return -1
}

func IsNearSpencer(client Entity) bool {
	var healers int = GetEntProp(client, Prop_Send, "m_nNumHealers")
	medics := 0
	if healers > 0 {
		for i:=MaxClients; i > 0; i-- {
			if IsValidClient(i) && GetHealingTarget(i) == client {
				medics++
			}
		}
	}
	return (healers > medics)
}

func SpawnRandomHealth() {
	iEnt, spawned := MaxClients+1, 0
	var vPos, vAng Vec3
	var maxlim, minlim int = g_vsh2.m_hCvars.HealthKitLimitMax.IntValue, g_vsh2.m_hCvars.HealthKitLimitMin.IntValue
	for iEnt = FindEntityByClassname(iEnt, "info_player_teamspawn"); iEnt != -1; iEnt = FindEntityByClassname(iEnt, "info_player_teamspawn") {
		if spawned >= minlim {
			if GetRandomInt(0, 3) {
				continue;
			}
		}
		if spawned >= maxlim {
			break;
		}
		GetEntPropVector(iEnt, Prop_Send, "m_vecOrigin", vPos)
		GetEntPropVector(iEnt, Prop_Send, "m_angRotation", vAng)
		var healthkit int = CreateEntityByName("item_healthkit_small")
		TeleportEntity(healthkit, vPos, vAng, NULL_VECTOR)
		DispatchSpawn(healthkit)
		if g_vsh2.m_hCvars.Enabled.BoolValue {
			SetEntProp(healthkit, Prop_Send, "m_iTeamNum", VSH2Team_Red, 4)
		} else {
			SetEntProp(healthkit, Prop_Send, "m_iTeamNum", VSH2Team_Neutral, 4)
		}
		spawned++
	}
}


func OnPluginStart() {
	for i := MaxClients; i > 0; i-- {
		if IsClientInGame(i) {
			OnClientPutInServer(i)
		}
	}
}

func OnClientPutInServer(client Entity) {
	SDKHook(client, SDKHooks_OnTakeDamage, OnTakeDamage)
}

func OnTakeDamage(victim Entity, attacker, inflictor, damagetype, weapon *Entity, damage *float, damageForce, damagePos *Vec3) Action {
	if IsValidEntity(*weapon) {
		if GetEntProp(*weapon, SM.Prop_Send, "m_iItemDefinitionIndex")==444 {
			damage *= 5.0
			return Plugin_Changed
		}
	}
	return Plugin_Continue
}

func BackPackReload(weapref int, flHolsterTime *float, flSecondsDelay float, SingleReload bool, ReloadInterval int) {
	var weapon int = EntRefToEntIndex(weapref)
	if weapon <= 0 {
		return
	}

	var client int = GetOwner(weapon)
	if client <= 0 || flSecondsDelay >= 100.0 {
		return
	} else if !IsPlayerAlive(client) {
		return
	}
	if bChatSpam[client] {
		PrintToConsole(client, "got past weapon and client checks")
	}
	if (GetGameTime()-*flHolsterTime) > flSecondsDelay {
		*flHolsterTime = GetGameTime()
		if GetWeaponClip(weapon) < ClipTable[weapon] {
			if SingleReload {
				if ClipTable[weapon]-GetWeaponClip(weapon) < ReloadInterval {
					ReloadInterval = ClipTable[weapon]-GetWeaponClip(weapon)
				}
				if GetWeaponAmmo(weapon) < ReloadInterval {
					ReloadInterval = GetWeaponAmmo(weapon)
				}
				if ReloadInterval < 1 {
					return
				}
				SetWeaponClip(weapon, GetWeaponClip(weapon)+ReloadInterval)
				SetWeaponAmmo(weapon, GetWeaponAmmo(weapon)-ReloadInterval)
				if bChatSpam[client] {
					PrintToConsole(client, "did single reload")
				}
				EmitSoundToClient(client, strSound)
				EmitSoundToClient(client, strSound)
			} else {
				ReloadInterval = ClipTable[weapon]-GetWeaponClip(weapon)
				if GetWeaponAmmo(weapon) < ReloadInterval {
					ReloadInterval = GetWeaponAmmo(weapon)
				}
				if ReloadInterval < 1 {
					return
				}

				SetWeaponClip(weapon, GetWeaponClip(weapon)+ReloadInterval)
				SetWeaponAmmo(weapon, GetWeaponAmmo(weapon)-ReloadInterval)
				if bChatSpam[client] {
					PrintToConsole(client, "did mag reload")
				}
				EmitSoundToClient(client, strSound)
				EmitSoundToClient(client, strSound)
			}
		} else if GetWeaponClip(weapon) > ClipTable[weapon] {
			ClipTable[weapon] = GetWeaponClip(weapon)
		}
	}
}

func CmdGoTest(client Entity, args int) Action {
    var cmd[64] char
    s := len(cmd)
    GetCmdArgString(cmd, s)
    PrintToServer("%s", cmd)
}