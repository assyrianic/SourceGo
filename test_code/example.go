package main

import (
	"sourcemod"
	"tf2_stocks"
	"sdkhooks"
)

const (
	PLUGIN_VERSION = "1.2"
	PLYR           = MAXPLAYERS+1
)

var (
	registrar = Plugin {
		name:        "example gopawn plugin.",
		author:      "Nergal/Assyrian/Ashurian",
		description: "kektus",
		version:     PLUGIN_VERSION,
		url:         "Alliedmodders",
	}
	
	bDetectProjs [PLYR]bool
	
	PluginEnabled, DetectGrenades, DetectStickies, DetectGrenadeRadius, DetectStickyRadius, DetectFriendly ConVar
	
	names map[string]SM.Entity
)

func OnPluginStart() {
	PluginEnabled = CreateConVar("projindic_enabled", "1", "Enable Projectile Indicator plugin", FCVAR_NONE, true, 0.0, true, 1.0)
	
	DetectGrenades = CreateConVar("projindic_grenades", "1", "Enable the Projectile Indicator plugin to detect pipe grenades", FCVAR_NONE, true, 0.0, true, 1.0) // THIS INCLUDES CANNONBALLS
	
	DetectStickies = CreateConVar("projindic_stickies", "1", "Enable the Projectile Indicator plugin to detect stickybombs", FCVAR_NONE, true, 0.0, true, 1.0)
	
	DetectGrenadeRadius = CreateConVar("projindic_grenaderadius", "300.0", "Detection radius for pipe grenades in Hammer Units", FCVAR_NONE, true, 0.0, true, 99999.0)
	
	DetectStickyRadius = CreateConVar("projindic_stickyradius", "300.0", "Detection radius for stickybombs in Hammer Units", FCVAR_NONE, true, 0.0, true, 99999.0)
	
	DetectFriendly = CreateConVar("projindic_detectfriendly", "1", "Detect friendly projectiles", FCVAR_NONE, true, 0.0, true, 1.0)
	
	RegConsoleCmd("sm_detect", ToggleIndicator)
	RegConsoleCmd("sm_indic", ToggleIndicator)
	RegConsoleCmd("sm_forcedetect", ForceDetection)
	RegConsoleCmd("sm_forcedetection", ForceDetection)
	
	AutoExecConfig(true, "Projectile-Indicator")
	
	for i:=MaxClients; i; i-- {
		if !IsValidClient(i) {
			continue
		}
		OnClientPutInServer(i)
	}
}

func OnClientPutInServer(client int) {
	bDetectProjs[client] = true
	SDKHook(client, SDKHook_PostThinkPost, IndicatorThink)
}

func IndicatorThink(client int) {
	if !PluginEnabled.BoolValue || client <= 0 || !bDetectProjs[client] || !IsPlayerAlive(client) || IsClientObserver(client) {
		return
	}
	
	var screenx, screeny       float
	var GrenDelta, StickyDelta Vec3
	
	if DetectGrenades.BoolValue {
		iEntity := FindEntityByClassname(iEntity, "tf_projectile_pipe")
		for iEntity != -1 {
			refentity := EntIndexToEntRef(iEntity)
			if GetDistFromProj(client, refentity) > DetectGrenadeRadius.FloatValue {
				continue
			}
			
			thrower := GetThrower(iEntity)
			if thrower != -1 && GetClientTeam(thrower) == GetClientTeam(client) && !DetectFriendly.BoolValue {
				continue
			}
			GrenDelta = GetDeltaVector(client, refentity)
			NormalizeVector(GrenDelta, GrenDelta)
			GetProjPosToScreen(client, GrenDelta, screenx, screeny)
			DrawIndicator(client, screenx, screeny, "O")
			iEntity = FindEntityByClassname(iEntity, "tf_projectile_pipe")
		}
	}
	
	if DetectStickies.BoolValue {
		iEntity := FindEntityByClassname(iEntity, "tf_projectile_pipe_remote")
		for iEntity != -1 {
			refentity := EntIndexToEntRef(iEntity)
			if GetDistFromProj(client, refentity) > DetectStickyRadius.FloatValue {
				continue
			}
			
			thrower := GetThrower(iEntity)
			if thrower != -1 && GetClientTeam(thrower) == GetClientTeam(client) && !DetectFriendly.BoolValue {
				continue
			}
			StickyDelta = GetDeltaVector(client, refentity)
			NormalizeVector(StickyDelta, StickyDelta)
			GetProjPosToScreen(client, StickyDelta, screenx, screeny)
			DrawIndicator(client, screenx, screeny, "X")
			iEntity = FindEntityByClassname(iEntity, "tf_projectile_pipe")
		}
	}
}

func ToggleIndicator(client, args int) Action {
	if !PluginEnabled.BoolValue {
		return Plugin_Continue
	}
	bDetectProjs[client] = true
	ReplyToCommand(client, "Projectile Indicator on")
	return Plugin_Handled
}

func ForceDetection(client, args int) Action {
	if !PluginEnabled.BoolValue {
		return Plugin_Handled
	}
	
	if args < 1 {
		ReplyToCommand(client, "[Projectile Indicator] Usage: sm_forcedetect <player/target>")
		return Plugin_Handled
	}
	var name, target_name [PLATFORM_MAX_PATH]int8
	GetCmdArg(1, name, sizeof(name)) /// rework into `var name [MAX_PATH]char = arg[0]`?

	var target_list [PLYR]int
	var tn_is_ml bool
	if target_count := ProcessTargetString(name, client, target_list, MAXPLAYERS, COMMAND_FILTER_NO_BOTS, target_name, sizeof(target_name), tn_is_ml); target_count <= 0 {
		/** This function replies to the admin with a failure message */
		ReplyToTargetError(client, target_count)
		return Plugin_Handled
	}
	
	for i:=0; i<target_count; i++ {
		if IsValidClient(target_list[i]) {
			bDetectProjs[target_list[i]] = true
		}
	}
	ReplyToCommand(client, "Forcing Projectile Indicators")
	return Plugin_Handled
}

func isValidClient(client int, replaycheck bool) bool {
	if !IsClientValid(client) || !IsClientInGame(client) || GetEntProp(client, Prop_Send, "m_bIsCoaching") {
		return false
	} else if replaycheck && (IsClientSourceTV(client) || IsClientReplay(client)) {
		return false
	}
	return true
}


/// stock functions are represented as snake_case or camelCase
func getOwner(ent int) int {
	if IsValidEntity(ent) {
		return GetEntPropEnt(ent, Prop_Send, "m_hOwnerEntity")
	} else { return -1 }
}

func getThrower(ent int) int {
	if IsValidEntity(ent) {
		return GetEntPropEnt(ent, Prop_Send, "m_hThrower")
	} else { return -1 }
}

/**
 * Gets the position of the projectile from the player's position
 * and converts the data to screen numbers
 *
 * @param vecDelta	delta vector to work from
 * @param xpos		x position of the screen
 * @param ypos		y position of the screen
 * @noreturn
 * @note		set xpos and ypos as references so we can "return" both of them.
 * @props		Code by Valve from their Source Engine hud_damageindicator.cpp
 */
func GetProjPosToScreen(client int, vecDelta Vec3) (xpos, ypos float) {
	var playerAngles Vec3
	GetClientEyeAngles(client, playerAngles)
	
	var vecforward, right Vec3
	var up Vec3 = Vec3{0.0, 0.0, 1.0}
	GetAngleVectors(playerAngles, vecforward, NULL_VECTOR, NULL_VECTOR)
	vecforward[2] = 0.0
	
	NormalizeVector(vecforward, vecforward)
	GetVectorCrossProduct(up, vecforward, right)

	front := GetVectorDotProduct(vecDelta, vecforward)
	side := GetVectorDotProduct(vecDelta, right)
	
	xpos = 360.0 * -front
	ypos = 360.0 * -side
	
	flRotation := (ArcTangent2(xpos, ypos) + FLOAT_PI) * (57.29577951)
	
	yawRadians := -flRotation * 0.017453293

	// Rotate it around the circle
	xpos = ( 500 + (360.0 * Cosine(yawRadians)) ) / 1000.0
	ypos = ( 500 - (360.0 * Sine(yawRadians)) ) / 1000.0
}

/**
 * gets the delta vector between player and projectile!
 *
 * @param entref	serial reference of the entity
 * @param vecBuffer	float buffer to store vector result
 * @return		delta vector from vecBuffer
 */
func GetDeltaVector(client, entref int) Vec3 {
	var vec Vec3
	proj := EntRefToEntIndex(entref)
	if proj <= 0 || !IsValidEntity(proj) {
		return vec
	}
	var vecPlayer, vecPos Vec3
	GetClientAbsOrigin(client, vecPlayer)
	GetEntPropVector(proj, Prop_Data, "m_vecAbsOrigin", vecPos)
	SubtractVectors(vecPlayer, vecPos, vec)
	return vec
}

/**
 * gets the distance between player and projectile!
 *
 * @param entref	serial reference of the entity
 * @return		distance between player and proj
 */
func GetDistFromProj(client, entref int) float {
	proj := EntRefToEntIndex(entref)
	if proj <= 0 || !IsValidEntity(proj) {
		return -1.0
	}
	var vecProjpos, vecClientpos Vec3
	GetEntPropVector(proj, Prop_Data, "m_vecAbsOrigin", vecProjpos)
	GetClientAbsOrigin(client, vecClientpos)
	return GetVectorDistance(vecClientpos, vecProjpos)
}

/**
 * Displays the Projectile indicator to alert player of nearly projectiles
 *
 * @param xpos		x position of the screen
 * @param ypos		y position of the screen
 * @noreturn
 */
/** This is how SetHudTextParams sets the x and y pos
	|----------------------------------------------------------------|
	|                           y 0.0                                |
	|                           |                                    |
	|                           |                                    |
	|                           |                                    |
	|                           |                                    |
	|                           |                                    |
	|                           |                                    |
	|x 0.0 ---------------------|------------------------------> 1.0 |
	|                           |                                    |
	|                           |                                    |
	|                           |                                    |
	|                           |                                    |
	|                           V                                    |
	|                           1.0                                  |
	|----------------------------------------------------------------|
*/
func DrawIndicator(client int, xpos, ypos float, textc string) {
	SetHudTextParams(xpos, ypos, 0.1, 255, 100, 0, 255, 0, 0.35, 0.0, 0.1)
	ShowHudText(client, -1, textc)
}
