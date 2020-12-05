/**
 * sourcemod/halflife.go
 * 
 * Copyright 2020 Nirari Technologies Alliedmodders LLC.
 * 
 * Permission is hereby granted free of charge to any person obtaining a copy of this software and associated documentation files (the "Software") to deal in the Software without restriction including without limitation the rights to use copy modify merge publish distribute sublicense and/or sell copies of the Software and to permit persons to whom the Software is furnished to do so subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND EXPRESS OR IMPLIED INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM DAMAGES OR OTHER LIABILITY WHETHER IN AN ACTION OF CONTRACT TORT OR OTHERWISE ARISING FROM OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package main


const (
	SOURCE_SDK_UNKNOWN =           0      /**< Could not determine the engine version */
	SOURCE_SDK_ORIGINAL =         10      /**< Original Source engine (still used by "The Ship") */
	SOURCE_SDK_DARKMESSIAH =      15      /**< Modified version of original engine used by Dark Messiah (no SDK) */
	SOURCE_SDK_EPISODE1 =         20      /**< SDK+Engine released after Episode 1 */
	SOURCE_SDK_EPISODE2 =         30      /**< SDK+Engine released after Episode 2/Orange Box */
	SOURCE_SDK_BLOODYGOODTIME =   32      /**< Modified version of ep2 engine used by Bloody Good Time (no SDK) */
	SOURCE_SDK_EYE =              33      /**< Modified version of ep2 engine used by E.Y.E Divine Cybermancy (no SDK) */
	SOURCE_SDK_CSS =              34      /**< Sometime-older version of Source 2009 SDK+Engine used for Counter-Strike: Source */
	SOURCE_SDK_EPISODE2VALVE =    35      /**< SDK+Engine released after Episode 2/Orange Box "Source 2009" or "Source MP" */
	SOURCE_SDK_LEFT4DEAD =        40      /**< Engine released after Left 4 Dead (no SDK yet) */
	SOURCE_SDK_LEFT4DEAD2 =       50      /**< Engine released after Left 4 Dead 2 (no SDK yet) */
	SOURCE_SDK_ALIENSWARM =       60      /**< SDK+Engine released after Alien Swarm */
	SOURCE_SDK_CSGO =             80      /**< Engine released after CS:GO (no SDK yet) */
	SOURCE_SDK_DOTA =             90      /**< Engine released after Dota 2 (no SDK) */

	MOTDPANEL_TYPE_TEXT =          0      /**< Treat msg as plain text */
	MOTDPANEL_TYPE_INDEX =         1      /**< Msg is auto determined by the engine */
	MOTDPANEL_TYPE_URL =           2      /**< Treat msg as an URL link */
	MOTDPANEL_TYPE_FILE =          3      /**< Treat msg as a filename to be opened */
)


type DialogType int
const (
	DialogType_Msg = DialogType(0)     /**< just an on screen message */
	DialogType_Menu        /**< an options menu */
	DialogType_Text        /**< a richtext dialog */
	DialogType_Entry       /**< an entry box */
	DialogType_AskConnect   /**< ask the client to connect to a specified IP */
)


type EngineVersion int
const (
	Engine_Unknown = EngineVersion(0)             /**< Could not determine the engine version */
	Engine_Original            /**< Original Source Engine (used by The Ship) */
	Engine_SourceSDK2006       /**< Episode 1 Source Engine (second major SDK) */
	Engine_SourceSDK2007       /**< Orange Box Source Engine (third major SDK) */
	Engine_Left4Dead           /**< Left 4 Dead */
	Engine_DarkMessiah         /**< Dark Messiah Multiplayer (based on original engine) */
	Engine_Left4Dead2 = EngineVersion(7)      /**< Left 4 Dead 2 */
	Engine_AlienSwarm          /**< Alien Swarm (and Alien Swarm SDK) */
	Engine_BloodyGoodTime      /**< Bloody Good Time */
	Engine_EYE                 /**< E.Y.E Divine Cybermancy */
	Engine_Portal2             /**< Portal 2 */
	Engine_CSGO                /**< Counter-Strike: Global Offensive */
	Engine_CSS                 /**< Counter-Strike: Source */
	Engine_DOTA                /**< Dota 2 */
	Engine_HL2DM               /**< Half-Life 2 Deathmatch */
	Engine_DODS                /**< Day of Defeat: Source */
	Engine_TF2                 /**< Team Fortress 2 */
	Engine_NuclearDawn         /**< Nuclear Dawn */
	Engine_SDK2013             /**< Source SDK 2013 */
	Engine_Blade               /**< Blade Symphony */
	Engine_Insurgency          /**< Insurgency (2013 Retail version)*/
	Engine_Contagion           /**< Contagion */
	Engine_BlackMesa           /**< Black Mesa Multiplayer */
	Engine_DOI                  /**< Day of Infamy */
)


type FindMapResult int
const (
	// A direct match for this name was found
	FindMap_Found = FindMapResult(0)
	// No match for this map name could be found.
	FindMap_NotFound
	// A fuzzy match for this map name was found.
	// Ex: cp_dust -> cp_dustbowl c1m1 -> c1m1_hotel
	// Only supported for maps that the engine knows about. (This excludes workshop maps on Orangebox).
	FindMap_FuzzyMatch
	// A non-canonical match for this map name was found.
	// Ex: workshop/1234 -> workshop/cp_qualified_name.ugc1234
	// Only supported on "Orangebox" games with workshop support.
	FindMap_NonCanonical
	// No currently available match for this map name could be found but it may be possible to load
	// Only supported on "Orangebox" games with workshop support.
	FindMap_PossiblyAvailable
)

const INVALID_ENT_REFERENCE = 0xFFFFFFFF


func LogToGame(format string, args ...any)
func SetRandomSeed(seed int)
func GetRandomFloat(fMin, fMax float) float
func GetRandomInt(nmin, nmax int) int
func IsMapValid(map_name string) bool
func FindMap(map_name string, foundmap []char, maxlen int) FindMapResult
func GetMapDisplayName(map_name string, displayName []char, maxlen int) bool
func IsDedicatedServer() bool
func GetEngineTime() float
func GetGameTime() float
func GetGameTickCount() int
func GetGameFrameTime() float
func GetGameDescription(buffer []char, maxlength int, original bool) int
func GetGameFolderName(buffer []char, maxlength int) int
func GetCurrentMap(buffer []char, maxlength int) int
func PrecacheModel(model string, preload bool) int
func PrecacheSentenceFile(file string, preload bool) int
func PrecacheDecal(decal string, preload bool) int
func PrecacheGeneric(generic string, preload bool) int
func IsModelPrecached(model string) bool
func IsDecalPrecached(decal string) bool
func IsGenericPrecached(generic string) bool
func PrecacheSound(sound string, preload bool) bool
func IsSoundPrecached(sound string) bool
func CreateDialog(client int, kv KeyValues, dialog_type DialogType)
func GetEngineVersion() EngineVersion

func PrintToChat(client Entity, format string, args ...any)
func PrintToChatAll(format string, args ...any)

func PrintCenterText(client Entity, format string, args ...any)
func PrintCenterTextAll(format string, args ...any)

func PrintHintText(client Entity, format string, args ...any)
func PrintHintTextToAll(format string, args ...any)

func ShowVGUIPanel(client Entity, name string, Kv KeyValues, show bool)
func CreateHudSynchronizer() Handle
func SetHudTextParams(x, y, holdTime float, r, g, b, a, effect int, fxTime, fadeIn, fadeOut float)
func SetHudTextParamsEx(x, y, holdTime float, color1, color2 [4]int, effect int, fxTime, fadeIn, fadeOut float)
func ShowSyncHudText(client Entity, sync Handle, format string, args ...any) int
func ClearSyncHud(client Entity, sync Handle)
func ShowHudText(client, channel int, format string, args ...any) int
func ShowMOTDPanel(client Entity, title, msg string, motd_type int)
func DisplayAskConnectBox(client Entity, time float, ip, password string)
func EntIndexToEntRef(entity Entity) int
func EntRefToEntIndex(ref int) Entity
func MakeCompatEntRef(ref int) int


type ClientRangeType int
const (
	RangeType_Visibility = ClientRangeType(0)
	RangeType_Audibility
)

func GetClientsInRange(origin Vec3, rangeType ClientRangeType, clients []int, size int) int
func GetServerAuthId(authType AuthIdType, auth []char, maxlen int)
func GetServerSteamAccountId() int