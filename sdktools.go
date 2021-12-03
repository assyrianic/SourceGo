/**
 * sdktools.go
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
)


/** sdktools_engine */
const MAX_LIGHTSTYLES = 64

func SetClientViewEntity(client, entity Entity)
func SetLightStyle(style int, value string)
func GetClientEyePosition(client Entity, pos *Vec3)
/********************/


/** sdktools_functions */
func RemovePlayerItem(client, item Entity) bool
func GivePlayerItem(client Entity, item string, iSubType int) Entity
func GetPlayerWeaponSlot(client, slot int) Entity
func IgniteEntity(entity Entity, time float, npc bool, size float, level bool)
func ExtinguishEntity(entity Entity)
func TeleportEntity(entity Entity, origin, angles, velocity Vec3)
func ForcePlayerSuicide(client Entity)
func SlapPlayer(client, health int, sound bool)
func FindEntityByClassname(startEnt Entity, classname string) Entity
func GetClientEyeAngles(client Entity, ang *Vec3) bool
func CreateEntityByName(classname string, ForceEdictIndex int) int
func DispatchSpawn(entity Entity) bool
func DispatchKeyValue(entity Entity, keyName, value string) bool
func DispatchKeyValueFloat(entity Entity, keyName string, value float) bool
func DispatchKeyValueVector(entity Entity, keyName string, vec Vec3) bool
func GetClientAimTarget(client Entity, only_clients bool)
func GetTeamCount() int
func GetTeamName(index int, name []char, maxlength int)
func GetTeamScore(index int) int
func SetTeamScore(index, value int)
func GetTeamClientCount(index int) int
func GetTeamEntity(teamIndex int) int
func SetEntityModel(entity Entity, model string)
func GetPlayerDecalFile(client Entity, hex []char, maxlength int) bool
func GetPlayerJingleFile(client Entity, hex []char, maxlength int) bool
func GetServerNetStats(inAmount, outAmout *float)
func EquipPlayerWeapon(client, weapon Entity)
func ActivateEntity(entity Entity)
func SetClientInfo(client Entity, key, value string)
func SetClientName(client Entity, name string)
func GivePlayerAmmo(client, amount, ammotype int, suppressSound bool) int
/***********************/


/** sdktools_sound */
const (
	/**
	 * Sound should be from the target client.
	 */
	SOUND_FROM_PLAYER =       -2
	
	/**
	 * Sound should be from the listen server player.
	 */
	SOUND_FROM_LOCAL_PLAYER = -1
	
	/**
	 * Sound is from the world.
	 */
	SOUND_FROM_WORLD =        0
	
	SNDCHAN_REPLACE = -1       /**< Unknown */
	SNDCHAN_AUTO = 0           /**< Auto */
	SNDCHAN_WEAPON = 1         /**< Weapons */
	SNDCHAN_VOICE = 2          /**< Voices */
	SNDCHAN_ITEM = 3           /**< Items */
	SNDCHAN_BODY = 4           /**< Player? */
	SNDCHAN_STREAM = 5         /**< "Stream channel from the static or dynamic area" */
	SNDCHAN_STATIC = 6         /**< "Stream channel from the static area" */
	SNDCHAN_VOICE_BASE = 7     /**< "Channel for network voice data" */
	SNDCHAN_USER_BASE = 135     /**< Anything >= this is allocated to game code */
	
	SND_NOFLAGS= 0             /**< Nothing */
	SND_CHANGEVOL = 1          /**< Change sound volume */
	SND_CHANGEPITCH = 2        /**< Change sound pitch */
	SND_STOP = 3               /**< Stop the sound */
	SND_SPAWNING = 4           /**< Used in some cases for ambients */
	SND_DELAY = 5              /**< Sound has an initial delay */
	SND_STOPLOOPING = 6        /**< Stop looping all sounds on the entity */
	SND_SPEAKER = 7            /**< Being played by a mic through a speaker */
	SND_SHOULDPAUSE = 8         /**< Pause if game is paused */
	
	
	SNDLEVEL_NONE = 0          /**< None */
	SNDLEVEL_RUSTLE = 20       /**< Rustling leaves */
	SNDLEVEL_WHISPER = 25      /**< Whispering */
	SNDLEVEL_LIBRARY = 30      /**< In a library */
	SNDLEVEL_FRIDGE = 45       /**< Refrigerator */
	SNDLEVEL_HOME = 50         /**< Average home (3.9 attn) */
	SNDLEVEL_CONVO = 60        /**< Normal conversation (2.0 attn) */
	SNDLEVEL_DRYER = 60        /**< Clothes dryer */
	SNDLEVEL_DISHWASHER = 65   /**< Dishwasher/washing machine (1.5 attn) */
	SNDLEVEL_CAR = 70          /**< Car or vacuum cleaner (1.0 attn) */
	SNDLEVEL_NORMAL = 75       /**< Normal sound level */
	SNDLEVEL_TRAFFIC = 75      /**< Busy traffic (0.8 attn) */
	SNDLEVEL_MINIBIKE = 80     /**< Mini-bike alarm clock (0.7 attn) */
	SNDLEVEL_SCREAMING = 90    /**< Screaming child (0.5 attn) */
	SNDLEVEL_TRAIN = 100       /**< Subway train pneumatic drill (0.4 attn) */
	SNDLEVEL_HELICOPTER = 105  /**< Helicopter */
	SNDLEVEL_SNOWMOBILE = 110  /**< Snow mobile */
	SNDLEVEL_AIRCRAFT = 120    /**< Auto horn aircraft */
	SNDLEVEL_RAIDSIREN = 130   /**< Air raid siren */
	SNDLEVEL_GUNFIRE = 140     /**< Gunshot jet engine (0.27 attn) */
	SNDLEVEL_ROCKET = 180       /**< Rocket launching (0.2 attn) */
	
	SNDVOL_NORMAL =       1.0     /**< Normal volume */
	SNDPITCH_NORMAL =     100     /**< Normal pitch */
	SNDPITCH_LOW =        95      /**< A low pitch */
	SNDPITCH_HIGH =       120     /**< A high pitch */
	SNDATTN_NONE =        0.0     /**< No attenuation */
	SNDATTN_NORMAL =      0.8     /**< Normal attenuation */
	SNDATTN_STATIC =      1.25    /**< Static attenuation? */
	SNDATTN_RICOCHET =    1.5     /**< Ricochet effect */
	SNDATTN_IDLE =        2.0     /**< Idle attenuation? */
)

func PrefetchSound(name string)
func EmitAmbientSound(name string, pos Vec3, entity, level, flags int, vol float, pitch int, delay float)
func FadeClientVolume(client Entity, percent, outtime, holdtime, intime float)
func StopSound(entity, channel int, name string)
func EmitSound(clients []int, numClients int, sample string, entity, channel, level, flags int, volume float, pitch, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float, origins ...Vec3)
func EmitSoundEntry(clients []int, numClients int, soundEntry, sample string, entity, channel, level, flags int, volume float, pitch, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float, origins ...Vec3)
func EmitSentence(clients []int, numClients, sentence, entity, channel, level, flags int, volume float, pitch, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float, origins ...Vec3)
func GetDistGainFromSoundLevel(soundlevel int, distance float) float


type (
	AmbientSHook func(sample PathStr, entity *int, volume *float, level, pitch *int, pos *Vec3, flags *int, delay *float) Action
	NormalSHook func(clients [MAXPLAYERS]int, numClients *int, sample PathStr, entity, channel *int, volume *float, level, pitch, flags *int, soundEntry PathStr, seed *int) Action
)

func AddAmbientSoundHook(hook AmbientSHook)
func RemoveAmbientSoundHook(hook AmbientSHook)
func AddNormalSoundHook(hook NormalSHook)
func RemoveNormalSoundHook(hook NormalSHook)
func EmitSoundToClient(client int, sample string, entity, channel, level, flags int, volume float, pitch, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float)
func EmitSoundToAll(sample string, entity, channel, level, flags int, volume float, pitch, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float)
func ATTN_TO_SNDLEVEL(attn float) int
func GetGameSoundParams(gameSound string, channel, soundLevel *int, volume *float, pitch *int, sample []char, maxlength, entity int) bool
func EmitGameSound(clients []int, numClients int, gameSound string, entity, flags, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float) bool
func EmitAmbientGameSound(gameSound string, pos Vec3, entity, flags int, delay float) bool
func EmitGameSoundToClient(client int, gameSound string, entity, flags, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float) bool
func EmitGameSoundToAll(gameSound string, entity, flags, speakerentity int, origin, dir Vec3, updatePos bool, soundtime float) bool
func PrecacheScriptSound(soundname string) bool
/*******************/


/** sdktools_stringtables */
const (
	INVALID_STRING_TABLE = -1
	INVALID_STRING_INDEX = -1
)

func FindStringTable(name string) int
func GetNumStringTables() int
func GetStringTableNumStrings(tableidx int) int
func GetStringTableMaxStrings(tableidx int) int
func GetStringTableName(tableidx int, name []char, maxlength int) int
func FindStringIndex(tableidx int, str string) int
func ReadStringTable(tableidx, stringidx int, buffer []char, maxlength int) int
func GetStringTableDataLength(tableidx, stringidx int) int
func GetStringTableData(tableidx, stringidx int, buffer []char, maxlength int) int
func SetStringTableData(tableidx, stringidx int, buffer []char, maxlength int) int
func AddToStringTable(tableidx int, str, userdata string, length int)
func LockStringTables(lock bool) bool
func AddFileToDownloadsTable(filename string)
/**************************/


/** sdktools_trace */
const (
	CONTENTS_EMPTY =                   0           /**< No contents. */
	CONTENTS_SOLID =                   0x1         /**< an eye is never valid in a solid . */
	CONTENTS_WINDOW =                  0x2         /**< translucent, but not watery (glass). */
	CONTENTS_AUX =                     0x4
	CONTENTS_GRATE =                   0x8         /**< alpha-tested "grate" textures.  Bullets/sight pass through, but solids don't. */
	CONTENTS_SLIME =                   0x10
	CONTENTS_WATER =                   0x20
	CONTENTS_MIST =                    0x40
	CONTENTS_OPAQUE =                  0x80        /**< things that cannot be seen through (may be non-solid though). */
	LAST_VISIBLE_CONTENTS =            0x80
	ALL_VISIBLE_CONTENTS =             (LAST_VISIBLE_CONTENTS | (LAST_VISIBLE_CONTENTS-1))
	CONTENTS_TESTFOGVOLUME =           0x100
	CONTENTS_UNUSED5 =                 0x200
	CONTENTS_UNUSED6 =                 0x4000
	CONTENTS_TEAM1 =                   0x800       /**< per team contents used to differentiate collisions. */
	CONTENTS_TEAM2 =                   0x1000      /**< between players and objects on different teams. */
	CONTENTS_IGNORE_NODRAW_OPAQUE =    0x2000      /**< ignore CONTENTS_OPAQUE on surfaces that have SURF_NODRAW. */
	CONTENTS_MOVEABLE =                0x4000      /**< hits entities which are MOVETYPE_PUSH (doors, plats, etc) */
	CONTENTS_AREAPORTAL =              0x8000      /**< remaining contents are non-visible, and don't eat brushes. */
	CONTENTS_PLAYERCLIP =              0x10000
	CONTENTS_MONSTERCLIP =             0x20000
	
	CONTENTS_CURRENT_0 =      0x40000
	CONTENTS_CURRENT_90 =     0x80000
	CONTENTS_CURRENT_180 =    0x100000
	CONTENTS_CURRENT_270 =    0x200000
	CONTENTS_CURRENT_UP =     0x400000
	CONTENTS_CURRENT_DOWN =   0x800000
	
	CONTENTS_ORIGIN =         0x1000000   /**< removed before bsp-ing an entity. */
	CONTENTS_MONSTER =        0x2000000   /**< should never be on a brush, only in game. */
	CONTENTS_DEBRIS =         0x4000000
	CONTENTS_DETAIL =         0x8000000   /**< brushes to be added after vis leafs. */
	CONTENTS_TRANSLUCENT =    0x10000000  /**< auto set if any surface has trans. */
	CONTENTS_LADDER =         0x20000000
	CONTENTS_HITBOX =         0x40000000  /**< use accurate hitboxes on trace. */
	
	MASK_ALL =                    (0xFFFFFFFF)
	MASK_SOLID =                  (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_WINDOW|CONTENTS_MONSTER|CONTENTS_GRATE)                      /**< everything that is normally solid */
	MASK_PLAYERSOLID =            (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_PLAYERCLIP|CONTENTS_WINDOW|CONTENTS_MONSTER|CONTENTS_GRATE)  /**< everything that blocks player movement */
	MASK_NPCSOLID =               (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_MONSTERCLIP|CONTENTS_WINDOW|CONTENTS_MONSTER|CONTENTS_GRATE) /**< blocks npc movement */
	MASK_WATER =                  (CONTENTS_WATER|CONTENTS_MOVEABLE|CONTENTS_SLIME)                                                       /**< water physics in these contents */
	MASK_OPAQUE =                 (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_OPAQUE)                                                      /**< everything that blocks line of sight for AI, lighting, etc */
	MASK_OPAQUE_AND_NPCS =        (MASK_OPAQUE|CONTENTS_MONSTER)                                                                          /**< everything that blocks line of sight for AI, lighting, etc, but with monsters added. */
	MASK_VISIBLE =                (MASK_OPAQUE|CONTENTS_IGNORE_NODRAW_OPAQUE)                                                             /**< everything that blocks line of sight for players */
	MASK_VISIBLE_AND_NPCS =       (MASK_OPAQUE_AND_NPCS|CONTENTS_IGNORE_NODRAW_OPAQUE)                                                    /**< everything that blocks line of sight for players, but with monsters added. */
	MASK_SHOT =                   (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_MONSTER|CONTENTS_WINDOW|CONTENTS_DEBRIS|CONTENTS_HITBOX)     /**< bullets see these as solid */
	MASK_SHOT_HULL =              (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_MONSTER|CONTENTS_WINDOW|CONTENTS_DEBRIS|CONTENTS_GRATE)      /**< non-raycasted weapons see this as solid (includes grates) */
	MASK_SHOT_PORTAL =            (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_WINDOW)                                                      /**< hits solids (not grates) and passes through everything else */
	MASK_SOLID_BRUSHONLY =        (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_WINDOW|CONTENTS_GRATE)                                       /**< everything normally solid, except monsters (world+brush only) */
	MASK_PLAYERSOLID_BRUSHONLY =  (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_WINDOW|CONTENTS_PLAYERCLIP|CONTENTS_GRATE)                   /**< everything normally solid for player movement, except monsters (world+brush only) */
	MASK_NPCSOLID_BRUSHONLY =     (CONTENTS_SOLID|CONTENTS_MOVEABLE|CONTENTS_WINDOW|CONTENTS_MONSTERCLIP|CONTENTS_GRATE)                  /**< everything normally solid for npc movement, except monsters (world+brush only) */
	MASK_NPCWORLDSTATIC =         (CONTENTS_SOLID|CONTENTS_WINDOW|CONTENTS_MONSTERCLIP|CONTENTS_GRATE)                                    /**< just the world, used for route rebuilding */
	MASK_SPLITAREAPORTAL =        (CONTENTS_WATER|CONTENTS_SLIME) 
	
	SURF_LIGHT =       0x0001      /**< value will hold the light strength */
	SURF_SKY2D =       0x0002      /**< don't draw, indicates we should skylight + draw 2d sky but not draw the 3D skybox */
	SURF_SKY =         0x0004      /**< don't draw, but add to skybox */
	SURF_WARP =        0x0008      /**< turbulent water warp */
	SURF_TRANS =       0x0010
	SURF_NOPORTAL =    0x0020      /**< the surface can not have a portal placed on it */
	SURF_TRIGGER =     0x0040      /**< This is an xbox hack to work around elimination of trigger surfaces, which breaks occluders */
	SURF_NODRAW =      0x0080      /**< don't bother referencing the texture */

	SURF_HINT =        0x0100      /**< make a primary bsp splitter */

	SURF_SKIP =        0x0200      /**< completely ignore, allowing non-closed brushes */
	SURF_NOLIGHT =     0x0400      /**< Don't calculate light */
	SURF_BUMPLIGHT =   0x0800      /**< calculate three lightmaps for the surface for bumpmapping */
	SURF_NOSHADOWS =   0x1000      /**< Don't receive shadows */
	SURF_NODECALS =    0x2000      /**< Don't receive decals */
	SURF_NOCHOP =      0x4000      /**< Don't subdivide patches on this surface */
	SURF_HITBOX =      0x8000      /**< surface is part of a hitbox */ 
	
	PARTITION_SOLID_EDICTS =        (1 << 1) /**< every edict_t that isn't SOLID_TRIGGER or SOLID_NOT (and static props) */
	PARTITION_TRIGGER_EDICTS =      (1 << 2) /**< every edict_t that IS SOLID_TRIGGER */
	PARTITION_NON_STATIC_EDICTS =   (1 << 5) /**< everything in solid & trigger except the static props, includes SOLID_NOTs */
	PARTITION_STATIC_PROPS =        (1 << 7)
	
	DISPSURF_FLAG_SURFACE =         (1<<0)
	DISPSURF_FLAG_WALKABLE =        (1<<1)
	DISPSURF_FLAG_BUILDABLE =       (1<<2)
	DISPSURF_FLAG_SURFPROP1 =       (1<<3)
	DISPSURF_FLAG_SURFPROP2 =       (1<<4)
)


type RayType int
const (
	RayType_EndPoint = RayType(0)   /**< The trace ray will go from the start position to the end position. */
	RayType_Infinite    /**< The trace ray will go from the start position to infinity using a direction vector. */
)

type (
	TraceEntityFilter func(entity, contentsMask int, data any) bool
	TraceEntityEnumerator func(entity int, data any) bool
)

func TR_GetPointContents(pos Vec3, entindex *int) int
func TR_GetPointContentsEnt(entindex int, pos Vec3) int
func TR_TraceRay(pos, vec Vec3, flags int, rtype RayType)
func TR_TraceHull(pos, vec, mins, maxs Vec3, flags int)
func TR_EnumerateEntities(pos, vec Vec3, mask int, rtype RayType, enum_fn TraceEntityEnumerator, data any)
func TR_EnumerateEntitiesHull(pos, vec, mins, maxs Vec3, mask int, enum_fn TraceEntityEnumerator, data any)
func TR_TraceRayFilter(pos, vec Vec3, flags int, rtype RayType, filter TraceEntityFilter, data any)
func TR_TraceHullFilter(pos, vec, mins, maxs Vec3, flags int, filter TraceEntityFilter, data any)
func TR_ClipRayToEntity(pos, vec Vec3, flags int, rtype RayType, entity Entity)
func TR_ClipRayHullToEntity(pos, vec, mins, maxs Vec3, flags, entity int)
func TR_ClipCurrentRayToEntity(flags, entity int)
func TR_TraceRayEx(pos, vec Vec3, flags int, rtype RayType) Handle
func TR_TraceHullEx(pos, vec, mins, maxs Vec3, flags int) Handle
func TR_TraceRayFilterEx(pos, vec Vec3, flags int, rtype RayType, filter TraceEntityFilter, data any) Handle
func TR_TraceHullFilterEx(pos, vec, mins, maxs Vec3, flags int, filter TraceEntityFilter, data any) Handle
func TR_ClipRayToEntityEx(pos, vec Vec3, flags int, rtype RayType, entity Entity) Handle
func TR_ClipRayHullToEntityEx(pos, vec, mins, maxs Vec3, flags, entity int) Handle
func TR_ClipCurrentRayToEntityEx(flags, entity int) Handle
func TR_GetFraction(hndl Handle) float
func TR_GetFractionLeftSolid(hndl Handle) float
func TR_GetStartPosition(hndl Handle, pos *Vec3)
func TR_GetEndPosition(pos *Vec3, hndl Handle)
func TR_GetEntityIndex(hndl Handle) int
func TR_GetDisplacementFlags(hndl Handle) int
func TR_GetSurfaceName(hndl Handle, buffer []char, maxlen int)
func TR_GetSurfaceProps(hndl Handle) int
func TR_GetSurfaceFlags(hndl Handle) int
func TR_GetPhysicsBone(hndl Handle) int
func TR_AllSolid(hndl Handle) bool
func TR_StartSolid(hndl Handle) bool
func TR_DidHit(hndl Handle) bool
func TR_GetHitGroup(hndl Handle) int
func TR_GetHitBoxIndex(hndl Handle) int
func TR_GetPlaneNormal(hndl Handle, normal *Vec3)
func TR_PointOutsideWorld(pos *Vec3) bool
/*******************/


/** sdktools_tempents */
type TEHook func(te_name string, players []int, numClients int, delay float) Action

func AddTempEntHook(te_name string, hook TEHook)
func RemoveTempEntHook(te_name string, hook TEHook)
func TE_Start(te_name string)
func TE_IsValidProp(prop string) bool
func TE_WriteNum(prop string, value int)
func TE_ReadNum(prop string) int
func TE_WriteFloat(prop string, value float)
func TE_ReadFloat(prop string) float
func TE_WriteVector(prop string, vec Vec3)
func TE_ReadVector(prop string, vec *Vec3)
func TE_WriteAngles(prop string, angles Vec3)
func TE_WriteFloatArray(prop string, array []float, arraySize int)
func TE_Send(clients []int, numClients int, delay float)
func TE_WriteEncodedEnt(prop string, value int)
func TE_SendToAll(delay float)
func TE_SendToClient(client int, delay float)
func TE_SendToAllInRange(origin Vec3, rangeType ClientRangeType, delay float)
/**********************/


/** sdktools_tempents_stocks */
const (
	TE_EXPLFLAG_NONE =            0x0   /**< all flags clear makes default Half-Life explosion */
	TE_EXPLFLAG_NOADDITIVE =      0x1   /**< sprite will be drawn opaque (ensure that the sprite you send is a non-additive sprite) */
	TE_EXPLFLAG_NODLIGHTS =       0x2   /**< do not render dynamic lights */
	TE_EXPLFLAG_NOSOUND =         0x4   /**< do not play client explosion sound */
	TE_EXPLFLAG_NOPARTICLES =     0x8   /**< do not draw particles */
	TE_EXPLFLAG_DRAWALPHA =       0x10  /**< sprite will be drawn alpha */
	TE_EXPLFLAG_ROTATE =          0x20  /**< rotate the sprite randomly */
	TE_EXPLFLAG_NOFIREBALL =      0x40  /**< do not draw a fireball */
	TE_EXPLFLAG_NOFIREBALLSMOKE = 0x80  /**< do not draw smoke with the fireball */
	
	
	FBEAM_STARTENTITY =   0x00000001
	FBEAM_ENDENTITY =     0x00000002
	FBEAM_FADEIN =        0x00000004
	FBEAM_FADEOUT =       0x00000008
	FBEAM_SINENOISE =     0x00000010
	FBEAM_SOLID =         0x00000020
	FBEAM_SHADEIN =       0x00000040
	FBEAM_SHADEOUT =      0x00000080
	FBEAM_ONLYNOISEONCE = 0x00000100  /**< Only calculate our noise once */
	FBEAM_NOTILE =        0x00000200
	FBEAM_USE_HITBOXES =  0x00000400  /**< Attachment indices represent hitbox indices instead when this is set. */
	FBEAM_STARTVISIBLE =  0x00000800  /**< Has this client actually seen this beam's start entity yet? */
	FBEAM_ENDVISIBLE =    0x00001000  /**< Has this client actually seen this beam's end entity yet? */
	FBEAM_ISACTIVE =      0x00002000
	FBEAM_FOREVER =       0x00004000
	FBEAM_HALOBEAM =      0x00008000  /**< When drawing a beam with a halo, don't ignore the segments and endwidth */
)

func TE_SetupSparks(pos, dir Vec3, Magnitude, TrailLength int)
func TE_SetupSmoke(pos Vec3, Model int, Scale float, FrameRate int)
func TE_SetupDust(pos, dir Vec3, Size, Speed float)
func TE_SetupMuzzleFlash(pos, angles Vec3, Scale float, Type int)
func TE_SetupMetalSparks(pos, dir Vec3)
func TE_SetupEnergySplash(pos, dir Vec3, Explosive bool)
func TE_SetupArmorRicochet(pos, dir Vec3)
func TE_SetupGlowSprite(pos Vec3, Model int, Life, Size float, Brightness int)

func TE_SetupExplosion(pos Vec3, Model int, Scale float, Framerate, Flags, Radius, Magnitude int, normal Vec3, MaterialType int)

func TE_SetupBloodSprite(pos, dir Vec3, color [4]int, Size, SprayModel, BloodDropModel int)

func TE_SetupBeamRingPoint(center Vec3, Start_Radius, End_Radius float, ModelIndex, HaloIndex, StartFrame, FrameRate int, Life, Width, Amplitude float, Color [4]int, Speed, Flags int)

func TE_SetupBeamPoints(start, end Vec3, ModelIndex, HaloIndex, StartFrame, FrameRate int, Life, Width, EndWidth float, FadeLength int, Amplitude float, Color [4]int, Speed int)

func TE_SetupBeamLaser(StartEntity, EndEntity, ModelIndex, HaloIndex, StartFrame, FrameRate int, Life, Width, EndWidth float, FadeLength int, Amplitude float, Color [4]int, Speed int)

func TE_SetupBeamRing(StartEntity, EndEntity, ModelIndex, HaloIndex, StartFrame, FrameRate int, Life, Width, Amplitude float, Color [4]int, Speed, Flags int)

func TE_SetupBeamFollow(EntIndex, ModelIndex, HaloIndex int, Life, Width, EndWidth float, FadeLength int, Color [4]int)
/*****************************/


/** sdktools_voice */
const (
	VOICE_NORMAL =        0   /**< Allow the client to listen and speak normally. */
	VOICE_MUTED =         1   /**< Mutes the client from speaking to everyone. */
	VOICE_SPEAKALL =      2   /**< Allow the client to speak to everyone. */
	VOICE_LISTENALL =     4   /**< Allow the client to listen to everyone. */
	VOICE_TEAM =          8   /**< Allow the client to always speak to team, even when dead. */
	VOICE_LISTENTEAM =    16  /**< Allow the client to always hear teammates, including dead ones. */
)

type ListenOverride int
const (
	Listen_Default = ListenOverride(0) /**< Leave it up to the game */
	Listen_No          /**< Can't hear */
	Listen_Yes          /**< Can hear */
)

func SetClientListeningFlags(client, flags int)
func GetClientListeningFlags(client Entity) int
func SetListenOverride(Receiver, Sender int, override ListenOverride) bool
func GetListenOverride(Receiver, Sender int) ListenOverride
func IsClientMuted(Muter, Mutee int) bool
/*******************/


/** sdktools_variant_t */
func SetVariantBool(val bool)
func SetVariantString(val string)
func SetVariantInt(val int)
func SetVariantFloat(val float)
func SetVariantVector3D(val Vec3)
func SetVariantPosVector3D(val Vec3)
func SetVariantColor(val [4]int)
func SetVariantEntity(val Entity)
/***********************/


/** sdktools_entinput */
func AcceptEntityInput(dest Entity, input string, activator, caller Entity, outputid int) bool
/**********************/


/** sdktools_entoutput */
type EntityOutput func(output string, caller, activator Entity, delay float) Action

func HookEntityOutput(classname, output string, callback EntityOutput)
func UnhookEntityOutput(classname, output string, callback EntityOutput)
func HookSingleEntityOutput(entity Entity, output string, callback EntityOutput, once bool)
func UnhookSingleEntityOutput(entity Entity, output string, callback EntityOutput) bool
func FireEntityOutput(caller Entity, output string, activator Entity, delay float)
/***********************/


/** sdktools_hooks */
const FEATURECAP_PLAYERRUNCMD_11PARAMS string = "SDKTools PlayerRunCmd 11Params"
/*******************/


/** sdktools_gamerules */
type RoundState int
const (
	// initialize the game, create teams
	RoundState_Init = RoundState(0)
	
	// Before players have joined the game. Periodically checks to see if enough players are ready
	// to start a game. Also reverts to this when there are no active players
	RoundState_Pregame
	
	// The game is about to start, wait a bit and spawn everyone
	RoundState_StartGame
	
	// All players are respawned, frozen in place
	RoundState_Preround
	
	// Round is on, playing normally
	RoundState_RoundRunning
	
	// Someone has won the round
	RoundState_TeamWin
	
	// Noone has won, manually restart the game, reset scores
	RoundState_Restart
	
	// Noone has won, restart the game
	RoundState_Stalemate
	
	// Game is over, showing the scoreboard etc
	RoundState_GameOver
	
	// Game is over, doing bonus round stuff
	RoundState_Bonus
	
	// Between rounds
	RoundState_BetweenRounds
)

func GameRules_GetProp(prop string, size, element int) int
func GameRules_SetProp(prop string, value any, size, element int, changeState bool)
func GameRules_GetPropFloat(prop string, element int) float
func GameRules_SetPropFloat(prop string, value float, element int, changeState bool)
func GameRules_GetPropEnt(prop string, element int) Entity
func GameRules_SetPropEnt(prop string, other Entity, element int, changeState bool)
func GameRules_GetPropVector(prop string, vec *Vec3, element int)
func GameRules_SetPropVector(prop string, vec Vec3, element int, changeState bool)
func GameRules_GetPropString(prop string, buffer []char, maxlen int) int
func GameRules_SetPropString(prop, buffer string, changeState bool) int
func GameRules_GetRoundState() RoundState
/***********************/


/** sdktools_client */
func InactivateClient(client Entity)
func ReconnectClient(client Entity)
/********************/


/** sdktools_stocks */
func FindTeamByName(name string) int
/********************/



type SDKCallType int
const (
	SDKCall_Static = SDKCallType(0)         /**< Static call */
	SDKCall_Entity         /**< CBaseEntity call */
	SDKCall_Player         /**< CBasePlayer call */
	SDKCall_GameRules      /**< CGameRules call */
	SDKCall_EntityList     /**< CGlobalEntityList call */
	SDKCall_Raw            /**< |this| pointer with an arbitrary address */
)


type SDKLibrary int
const (
	SDKLibrary_Server = SDKLibrary(0)      /**< server.dll/server_i486.so */
	SDKLibrary_Engine       /**< engine.dll/engine_*.so */
)


type SDKFuncConfSource int
const (
	SDKConf_Virtual = SDKFuncConfSource(0)    /**< Read a virtual index from the Offsets section */
	SDKConf_Signature  /**< Read a signature from the Signatures section */
	SDKConf_Address    /**< Read an address from the Addresses section */
)

type SDKType int
const (
	SDKType_CBaseEntity = SDKType(0)    /**< CBaseEntity (always as pointer) */
	SDKType_CBasePlayer    /**< CBasePlayer (always as pointer) */
	SDKType_Vector         /**< Vector (pointer, byval, or byref) */
	SDKType_QAngle         /**< QAngles (pointer, byval, or byref) */
	SDKType_PlainOldData   /**< Integer/generic data <=32bit (any) */
	SDKType_Float          /**< Float (any) */
	SDKType_Edict          /**< edict_t (always as pointer) */
	SDKType_String         /**< NULL-terminated string (always as pointer) */
	SDKType_Bool           /**< Boolean (any) */
)

type SDKPassMethod int
const (
	SDKPass_Pointer = SDKPassMethod(0)        /**< Pass as a pointer */
	SDKPass_Plain          /**< Pass as plain data */
	SDKPass_ByValue        /**< Pass an object by value */
	SDKPass_ByRef          /**< Pass an object by reference */
)


const (
	VDECODE_FLAG_ALLOWNULL =      (1<<0)    /**< Allow NULL for pointers */
	VDECODE_FLAG_ALLOWNOTINGAME = (1<<1)    /**< Allow players not in game */
	VDECODE_FLAG_ALLOWWORLD =     (1<<2)    /**< Allow World entity */
	VDECODE_FLAG_BYREF =          (1<<3)    /**< Floats/ints by reference */

	VENCODE_FLAG_COPYBACK =       (1<<0)    /**< Copy back data once done */
)


func StartPrepSDKCall(calltype SDKCallType)
func PrepSDKCall_SetVirtual(vtblidx int)
func PrepSDKCall_SetSignature(lib SDKLibrary, signature string, bytes int) bool
func PrepSDKCall_SetAddress(addr Address) bool
func PrepSDKCall_SetFromConf(gameconf Handle, source SDKFuncConfSource, name string) bool
func PrepSDKCall_SetReturnInfo(sdktype SDKType, pass SDKPassMethod, decflags, encflags int)
func PrepSDKCall_AddParameter(sdktype SDKType, pass SDKPassMethod, decflags, encflags int)
func EndPrepSDKCall() Handle
func SDKCall(call Handle, params ...any) any
func GetPlayerResourceEntity() int