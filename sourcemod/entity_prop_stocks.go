/**
 * sourcemod/entity_prop_stocks.go
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


type MoveType int
const (
	MOVETYPE_NONE = MoveType(0)          /**< never moves */
	MOVETYPE_ISOMETRIC         /**< For players */
	MOVETYPE_WALK              /**< Player only - moving on the ground */
	MOVETYPE_STEP              /**< gravity special edge handling -- monsters use this */
	MOVETYPE_FLY               /**< No gravity but still collides with stuff */
	MOVETYPE_FLYGRAVITY        /**< flies through the air + is affected by gravity */
	MOVETYPE_VPHYSICS          /**< uses VPHYSICS for simulation */
	MOVETYPE_PUSH              /**< no clip to world push and crush */
	MOVETYPE_NOCLIP            /**< No gravity no collisions still do velocity/avelocity */
	MOVETYPE_LADDER            /**< Used by players only when going onto a ladder */
	MOVETYPE_OBSERVER          /**< Observer movement depends on player's observer mode */
	MOVETYPE_CUSTOM             /**< Allows the entity to describe its own physics */
)

type RenderMode int
const (
	RENDER_NORMAL = RenderMode(0)              /**< src */
	RENDER_TRANSCOLOR          /**< c*a+dest*(1-a) */
	RENDER_TRANSTEXTURE        /**< src*a+dest*(1-a) */
	RENDER_GLOW                /**< src*a+dest -- No Z buffer checks -- Fixed size in screen space */
	RENDER_TRANSALPHA          /**< src*srca+dest*(1-srca) */
	RENDER_TRANSADD            /**< src*a+dest */
	RENDER_ENVIRONMENTAL       /**< not drawn used for environmental effects */
	RENDER_TRANSADDFRAMEBLEND  /**< use a fractional frame value to blend between animation frames */
	RENDER_TRANSALPHAADD       /**< src + dest*(1-a) */
	RENDER_WORLDGLOW           /**< Same as kRenderGlow but not fixed size in screen space */
	RENDER_NONE                 /**< Don't render. */
)

type RenderFx int
const (
	RENDERFX_NONE = RenderFx(0)
	RENDERFX_PULSE_SLOW
	RENDERFX_PULSE_FAST
	RENDERFX_PULSE_SLOW_WIDE
	RENDERFX_PULSE_FAST_WIDE
	RENDERFX_FADE_SLOW
	RENDERFX_FADE_FAST
	RENDERFX_SOLID_SLOW
	RENDERFX_SOLID_FAST
	RENDERFX_STROBE_SLOW
	RENDERFX_STROBE_FAST
	RENDERFX_STROBE_FASTER
	RENDERFX_FLICKER_SLOW
	RENDERFX_FLICKER_FAST
	RENDERFX_NO_DISSIPATION
	RENDERFX_DISTORT           /**< Distort/scale/translate flicker */
	RENDERFX_HOLOGRAM          /**< kRenderFxDistort + distance fade */
	RENDERFX_EXPLODE           /**< Scale up really big! */
	RENDERFX_GLOWSHELL         /**< Glowing Shell */
	RENDERFX_CLAMP_MIN_SCALE   /**< Keep this sprite from getting very small (SPRITES only!) */
	RENDERFX_ENV_RAIN          /**< for environmental rendermode make rain */
	RENDERFX_ENV_SNOW          /**<  "        "            "     make snow */
	RENDERFX_SPOTLIGHT         /**< TEST CODE for experimental spotlight */
	RENDERFX_RAGDOLL           /**< HACKHACK: TEST CODE for signalling death of a ragdoll character */
	RENDERFX_PULSE_FAST_WIDER
	RENDERFX_MAX
)

const (
	IN_ATTACK =                (1 << 0)
	IN_JUMP =                  (1 << 1)
	IN_DUCK =                  (1 << 2)
	IN_FORWARD =               (1 << 3)
	IN_BACK =                  (1 << 4)
	IN_USE =                   (1 << 5)
	IN_CANCEL =                (1 << 6)
	IN_LEFT =                  (1 << 7)
	IN_RIGHT =                 (1 << 8)
	IN_MOVELEFT =              (1 << 9)
	IN_MOVERIGHT =             (1 << 10)
	IN_ATTACK2 =               (1 << 11)
	IN_RUN =                   (1 << 12)
	IN_RELOAD =                (1 << 13)
	IN_ALT1 =                  (1 << 14)
	IN_ALT2 =                  (1 << 15)
	IN_SCORE =                 (1 << 16)   /**< Used by client.dll for when scoreboard is held down */
	IN_SPEED =                 (1 << 17)   /**< Player is holding the speed key */
	IN_WALK =                  (1 << 18)   /**< Player holding walk key */
	IN_ZOOM =                  (1 << 19)   /**< Zoom key for HUD zoom */
	IN_WEAPON1 =               (1 << 20)   /**< weapon defines these bits */
	IN_WEAPON2 =               (1 << 21)   /**< weapon defines these bits */
	IN_BULLRUSH =              (1 << 22)
	IN_GRENADE1 =              (1 << 23)   /**< grenade 1 */
	IN_GRENADE2 =              (1 << 24)   /**< grenade 2 */
	IN_ATTACK3 =               (1 << 25)


// Note: these are only for use with GetEntityFlags and SetEntityFlags
//       and may not match the game's actual internal m_fFlags values.
// PLAYER SPECIFIC FLAGS FIRST BECAUSE WE USE ONLY A FEW BITS OF NETWORK PRECISION
	FL_ONGROUND =              (1 << 0)   /**< At rest / on the ground */
	FL_DUCKING =               (1 << 1)   /**< Player flag -- Player is fully crouched */
	FL_WATERJUMP =             (1 << 2)   /**< player jumping out of water */
	FL_ONTRAIN =               (1 << 3)   /**< Player is _controlling_ a train so movement commands should be ignored on client during prediction. */
	FL_INRAIN =                (1 << 4)   /**< Indicates the entity is standing in rain */
	FL_FROZEN =                (1 << 5)   /**< Player is frozen for 3rd person camera */
	FL_ATCONTROLS =            (1 << 6)   /**< Player can't move but keeps key inputs for controlling another entity */
	FL_CLIENT =                (1 << 7)   /**< Is a player */
	FL_FAKECLIENT =            (1 << 8)   /**< Fake client simulated server side; don't send network messages to them */
// NOTE if you move things up make sure to change this value
	PLAYER_FLAG_BITS =          9
// NON-PLAYER SPECIFIC (i.e. not used by GameMovement or the client .dll ) -- Can still be applied to players though
	FL_INWATER =               (1 << 9)   /**< In water */
	FL_FLY =                   (1 << 10)  /**< Changes the SV_Movestep() behavior to not need to be on ground */
	FL_SWIM =                  (1 << 11)  /**< Changes the SV_Movestep() behavior to not need to be on ground (but stay in water) */
	FL_CONVEYOR =              (1 << 12)
	FL_NPC =                   (1 << 13)
	FL_GODMODE =               (1 << 14)
	FL_NOTARGET =              (1 << 15)
	FL_AIMTARGET =             (1 << 16)  /**< set if the crosshair needs to aim onto the entity */
	FL_PARTIALGROUND =         (1 << 17)  /**< not all corners are valid */
	FL_STATICPROP =            (1 << 18)  /**< Eetsa static prop!  */
	FL_GRAPHED =               (1 << 19)  /**< worldgraph has this ent listed as something that blocks a connection */
	FL_GRENADE =               (1 << 20)
	FL_STEPMOVEMENT =          (1 << 21)  /**< Changes the SV_Movestep() behavior to not do any processing */
	FL_DONTTOUCH =             (1 << 22)  /**< Doesn't generate touch functions generates Untouch() for anything it was touching when this flag was set */
	FL_BASEVELOCITY =          (1 << 23)  /**< Base velocity has been applied this frame (used to convert base velocity into momentum) */
	FL_WORLDBRUSH =            (1 << 24)  /**< Not moveable/removeable brush entity (really part of the world but represented as an entity for transparency or something) */
	FL_OBJECT =                (1 << 25)  /**< Terrible name. This is an object that NPCs should see. Missiles for example. */
	FL_KILLME =                (1 << 26)  /**< This entity is marked for death -- will be freed by game DLL */
	FL_ONFIRE =                (1 << 27)  /**< You know... */
	FL_DISSOLVING =            (1 << 28)  /**< We're dissolving! */
	FL_TRANSRAGDOLL =          (1 << 29)  /**< In the process of turning into a client side ragdoll. */
	FL_UNBLOCKABLE_BY_PLAYER = (1 << 30)  /**< pusher that can't be blocked by the player */
	FL_FREEZING =              (1 << 31)  /**< We're becoming frozen! */
	FL_EP2V_UNKNOWN1 =         (1 << 31)  /**< Unknown */
)

func GetEntityFlags(entity Entity) int
func SetEntityFlags(entity, flags int)
func GetEntityMoveType(entity Entity) MoveType
func SetEntityMoveType(entity Entity, mt MoveType)
func GetEntityRenderMode(entity Entity) RenderMode
func SetEntityRenderMode(entity Entity, mode RenderMode)
func GetEntityRenderFx(entity Entity) RenderFx
func SetEntityRenderFx(entity Entity, fx RenderFx)
func GetEntityRenderColor(entity Entity, r, g, b, a *int)
func SetEntityRenderColor(entity, r, g, b, a int)
func GetEntityGravity(entity Entity) float
func SetEntityGravity(entity Entity, amount float)
func SetEntityHealth(entity, amount int)
func GetClientButtons(client Entity) int