/**
 * cstrike.go
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
	CSRoundEndReason int
	CSWeaponID int
)

const (
	CS_TEAM_NONE =        0   /**< No team yet. */
	CS_TEAM_SPECTATOR =   1   /**< Spectators. */
	CS_TEAM_T =           2   /**< Terrorists. */
	CS_TEAM_CT =          3   /**< Counter-Terrorists. */
	
	CS_SLOT_PRIMARY =     0   /**< Primary weapon slot. */
	CS_SLOT_SECONDARY =   1   /**< Secondary weapon slot. */
	CS_SLOT_KNIFE =       2   /**< Knife slot. */
	CS_SLOT_GRENADE =     3   /**< Grenade slot (will only return one grenade). */
	CS_SLOT_C4 =          4   /**< C4 slot. */
	
	CS_DMG_HEADSHOT =     (1 << 30)    /**< Headshot */
	
	
	CSRoundEnd_TargetBombed = CSRoundEndReason(0)           /**< Target Successfully Bombed! */
	CSRoundEnd_VIPEscaped                 /**< The VIP has escaped! - Doesn't exist on CS:GO */
	CSRoundEnd_VIPKilled                  /**< VIP has been assassinated! - Doesn't exist on CS:GO */
	CSRoundEnd_TerroristsEscaped          /**< The terrorists have escaped! */
	CSRoundEnd_CTStoppedEscape            /**< The CTs have prevented most of the terrorists from escaping! */
	CSRoundEnd_TerroristsStopped          /**< Escaping terrorists have all been neutralized! */
	CSRoundEnd_BombDefused                /**< The bomb has been defused! */
	CSRoundEnd_CTWin                      /**< Counter-Terrorists Win! */
	CSRoundEnd_TerroristWin               /**< Terrorists Win! */
	CSRoundEnd_Draw                       /**< Round Draw! */
	CSRoundEnd_HostagesRescued            /**< All Hostages have been rescued! */
	CSRoundEnd_TargetSaved                /**< Target has been saved! */
	CSRoundEnd_HostagesNotRescued         /**< Hostages have not been rescued! */
	CSRoundEnd_TerroristsNotEscaped       /**< Terrorists have not escaped! */
	CSRoundEnd_VIPNotEscaped              /**< VIP has not escaped! - Doesn't exist on CS:GO */
	CSRoundEnd_GameStart                  /**< Game Commencing! */
	
	// The below only exist on CS:GO
	CSRoundEnd_TerroristsSurrender        /**< Terrorists Surrender */
	CSRoundEnd_CTSurrender                /**< CTs Surrender */
	CSRoundEnd_TerroristsPlanted          /**< Terrorists Planted the bomb */
	CSRoundEnd_CTsReachedHostage           /**< CTs Reached the hostage */
	
	
	CSWeapon_NONE = CSWeaponID(0)
	CSWeapon_P228
	CSWeapon_GLOCK
	CSWeapon_SCOUT
	CSWeapon_HEGRENADE
	CSWeapon_XM1014
	CSWeapon_C4
	CSWeapon_MAC10
	CSWeapon_AUG
	CSWeapon_SMOKEGRENADE
	CSWeapon_ELITE
	CSWeapon_FIVESEVEN
	CSWeapon_UMP45
	CSWeapon_SG550
	CSWeapon_GALIL
	CSWeapon_FAMAS
	CSWeapon_USP
	CSWeapon_AWP
	CSWeapon_MP5NAVY
	CSWeapon_M249
	CSWeapon_M3
	CSWeapon_M4A1
	CSWeapon_TMP
	CSWeapon_G3SG1
	CSWeapon_FLASHBANG
	CSWeapon_DEAGLE
	CSWeapon_SG552
	CSWeapon_AK47
	CSWeapon_KNIFE
	CSWeapon_P90
	CSWeapon_SHIELD
	CSWeapon_KEVLAR
	CSWeapon_ASSAULTSUIT
	CSWeapon_NIGHTVISION //Anything below is CS:GO ONLY
	CSWeapon_GALILAR
	CSWeapon_BIZON
	CSWeapon_MAG7
	CSWeapon_NEGEV
	CSWeapon_SAWEDOFF
	CSWeapon_TEC9
	CSWeapon_TASER
	CSWeapon_HKP2000
	CSWeapon_MP7
	CSWeapon_MP9
	CSWeapon_NOVA
	CSWeapon_P250
	CSWeapon_SCAR17
	CSWeapon_SCAR20
	CSWeapon_SG556
	CSWeapon_SSG08
	CSWeapon_KNIFE_GG
	CSWeapon_MOLOTOV
	CSWeapon_DECOY
	CSWeapon_INCGRENADE
	CSWeapon_DEFUSER
	CSWeapon_HEAVYASSAULTSUIT
	//The rest are actual item definition indexes for CS:GO
	CSWeapon_CUTTERS = 56
	CSWeapon_HEALTHSHOT = 57
	CSWeapon_KNIFE_T = 59
	CSWeapon_M4A1_SILENCER = 60
	CSWeapon_USP_SILENCER = 61
	CSWeapon_CZ75A = 63
	CSWeapon_REVOLVER = 64
	CSWeapon_TAGGRENADE = 68
	CSWeapon_FISTS = 69
	CSWeapon_BREACHCHARGE = 70
	CSWeapon_TABLET = 72
	CSWeapon_MELEE = 74
	CSWeapon_AXE = 75
	CSWeapon_HAMMER = 76
	CSWeapon_SPANNER = 78
	CSWeapon_KNIFE_GHOST = 80
	CSWeapon_FIREBOMB = 81
	CSWeapon_DIVERSION = 82
	CSWeapon_FRAGGRENADE = 83
	CSWeapon_SNOWBALL = 84
	CSWeapon_BUMPMINE = 85
	CSWeapon_MAX_WEAPONS_NO_KNIFES // Max without the knife item defs useful when treating all knives as a regular knife.
	CSWeapon_BAYONET = 500
	CSWeapon_KNIFE_FLIP = 505
	CSWeapon_KNIFE_GUT = 506
	CSWeapon_KNIFE_KARAMBIT = 507
	CSWeapon_KNIFE_M9_BAYONET = 508
	CSWeapon_KNIFE_TATICAL = 509
	CSWeapon_KNIFE_FALCHION = 512
	CSWeapon_KNIFE_SURVIVAL_BOWIE = 514
	CSWeapon_KNIFE_BUTTERFLY = 515
	CSWeapon_KNIFE_PUSH = 516
	CSWeapon_KNIFE_URSUS = 519
	CSWeapon_KNIFE_GYPSY_JACKKNIFE = 520
	CSWeapon_KNIFE_STILETTO = 522
	CSWeapon_KNIFE_WIDOWMAKER = 523
	CSWeapon_MAX_WEAPONS // THIS MUST BE LAST EASY WAY TO CREATE LOOPS. When looping do CS_IsValidWeaponID(i) to check.
)

func CS_RespawnPlayer(client Entity)
func CS_SwitchTeam(client, team int)
func CS_DropWeapon(client, weaponIndex Entity, toss, blockhook bool)
func CS_TerminateRound(delay float, reason CSRoundEndReason, blockhook bool)
func CS_GetTranslatedWeaponAlias(alias string, weapon []char, size int)
func CS_GetWeaponPrice(client Entity, id CSWeaponID, defaultprice bool) int
func CS_GetClientClanTag(client Entity, buffer []char, size int) int
func CS_SetClientClanTag(client Entity, tag string)
func CS_GetTeamScore(team int) int
func CS_SetTeamScore(team, value int)
func CS_GetMVPCount(client Entity) int
func CS_SetMVPCount(client, value int)
func CS_GetClientContributionScore(client Entity) int
func CS_SetClientContributionScore(client, value int)
func CS_GetClientAssists(client int) int
func CS_SetClientAssists(client, value int)
func CS_AliasToWeaponID(alias string) CSWeaponID
func CS_WeaponIDToAlias(id CSWeaponID, destination []char, maxlen int) int
func CS_IsValidWeaponID(id CSWeaponID) bool
func CS_UpdateClientModel(client int)
func CS_ItemDefIndexToID(defindex int) CSWeaponID
func CS_WeaponIDToItemDefIndex(id CSWeaponID) int