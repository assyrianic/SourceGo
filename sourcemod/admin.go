/**
 * sourcemod/admin.go
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


type AdminFlag int
const (
	Admin_Reservation = AdminFlag(0)  /**< Reserved slot */
	Admin_Generic          /**< Generic admin abilities */
	Admin_Kick             /**< Kick another user */
	Admin_Ban              /**< Ban another user */
	Admin_Unban            /**< Unban another user */
	Admin_Slay             /**< Slay/kill/damage another user */
	Admin_Changemap        /**< Change the map */
	Admin_Convars          /**< Change basic convars */
	Admin_Config           /**< Change configuration */
	Admin_Chat             /**< Special chat privileges */
	Admin_Vote             /**< Special vote privileges */
	Admin_Password         /**< Set a server password */
	Admin_RCON             /**< Use RCON */
	Admin_Cheats           /**< Change sv_cheats and use its commands */
	Admin_Root             /**< All access by default */
	Admin_Custom1          /**< First custom flag type */
	Admin_Custom2          /**< Second custom flag type */
	Admin_Custom3          /**< Third custom flag type */
	Admin_Custom4          /**< Fourth custom flag type */
	Admin_Custom5          /**< Fifth custom flag type */
	Admin_Custom6          /**< Sixth custom flag type */
	AdminFlags_TOTAL
	
	ADMFLAG_RESERVATION        = (1<<0)      /**< Convenience macro for Admin_Reservation as a FlagBit */
	ADMFLAG_GENERIC            = (1<<1)      /**< Convenience macro for Admin_Generic as a FlagBit */
	ADMFLAG_KICK               = (1<<2)      /**< Convenience macro for Admin_Kick as a FlagBit */
	ADMFLAG_BAN                = (1<<3)      /**< Convenience macro for Admin_Ban as a FlagBit */
	ADMFLAG_UNBAN              = (1<<4)      /**< Convenience macro for Admin_Unban as a FlagBit */
	ADMFLAG_SLAY               = (1<<5)      /**< Convenience macro for Admin_Slay as a FlagBit */
	ADMFLAG_CHANGEMAP          = (1<<6)      /**< Convenience macro for Admin_Changemap as a FlagBit */
	ADMFLAG_CONVARS            = (1<<7)      /**< Convenience macro for Admin_Convars as a FlagBit */
	ADMFLAG_CONFIG             = (1<<8)      /**< Convenience macro for Admin_Config as a FlagBit */
	ADMFLAG_CHAT               = (1<<9)      /**< Convenience macro for Admin_Chat as a FlagBit */
	ADMFLAG_VOTE               = (1<<10)     /**< Convenience macro for Admin_Vote as a FlagBit */
	ADMFLAG_PASSWORD           = (1<<11)     /**< Convenience macro for Admin_Password as a FlagBit */
	ADMFLAG_RCON               = (1<<12)     /**< Convenience macro for Admin_RCON as a FlagBit */
	ADMFLAG_CHEATS             = (1<<13)     /**< Convenience macro for Admin_Cheats as a FlagBit */
	ADMFLAG_ROOT               = (1<<14)     /**< Convenience macro for Admin_Root as a FlagBit */
	ADMFLAG_CUSTOM1            = (1<<15)     /**< Convenience macro for Admin_Custom1 as a FlagBit */
	ADMFLAG_CUSTOM2            = (1<<16)     /**< Convenience macro for Admin_Custom2 as a FlagBit */
	ADMFLAG_CUSTOM3            = (1<<17)     /**< Convenience macro for Admin_Custom3 as a FlagBit */
	ADMFLAG_CUSTOM4            = (1<<18)     /**< Convenience macro for Admin_Custom4 as a FlagBit */
	ADMFLAG_CUSTOM5            = (1<<19)     /**< Convenience macro for Admin_Custom5 as a FlagBit */
	ADMFLAG_CUSTOM6            = (1<<20)     /**< Convenience macro for Admin_Custom6 as a FlagBit */
	
	AUTHMETHOD_STEAM          = "steam"     /**< SteamID based authentication */
	AUTHMETHOD_IP             = "ip"        /**< IP based authentication */
	AUTHMETHOD_NAME           = "name"      /**< Name based authentication */
)


type OverrideType int
const (
	Override_Command = OverrideType(1)   /**< Command */
	Override_CommandGroup   /**< Command group */
)

type OverrideRule int
const (
	Command_Deny = OverrideRule(0)
	Command_Allow = OverrideRule(1)
)

type AdmAccessMode int
const (
	Access_Real = AdmAccessMode(0)        /**< Access the user has inherently */
	Access_Effective    /**< Access the user has from their groups */
)


type AdminCachePart int
const (
	AdminCache_Overrides = AdminCachePart(0)       /**< Global overrides */
	AdminCache_Groups = AdminCachePart(1)          /**< All groups (automatically invalidates admins too) */
	AdminCache_Admins = AdminCachePart(2)          /**< All admins */
)

type ImmunityType int
const (
	Immunity_Default = ImmunityType(1)   /**< Deprecated. */
	Immunity_Global         /**< Deprecated. */
)

type (
	AdminId struct {
		GroupCount, ImmunityLevel int
	}
	GroupId struct {
		GroupImmunitiesCount, ImmunityLevel int
	}
)

var (
	INVALID_ADMIN_ID AdminId
	INVALID_GROUP_ID GroupId
)

func (AdminId) GetUsername(name []char, maxlength int)
func (AdminId) BindIdentity(authMethod, ident string) bool
func (AdminId) SetFlag(flag AdminFlag, enabled bool)
func (AdminId) HasFlag(flag AdminFlag, mode AdmAccessMode) bool
func (AdminId) GetFlags(mode AdmAccessMode) int
func (AdminId) InheritGroup(gid GroupId) bool
func (AdminId) GetGroup(index int, name []char, maxlength int) GroupId
func (AdminId) SetPassword(password string)
func (AdminId) GetPassword(buffer []char, maxlength int) bool
func (AdminId) CanTarget(other AdminId) bool


func (GroupId) HasFlag(flag AdminFlag) bool
func (GroupId) SetFlag(flag AdminFlag, enabled bool)
func (GroupId) GetFlags() int
func (GroupId) GetGroupImmunity(index int) GroupId
func (GroupId) AddGroupImmunity(other GroupId)
func (GroupId) GetCommandOverride(name string, overtype OverrideType, rule *OverrideRule) bool
func (GroupId) AddCommandOverride(name string, overtype OverrideType, rule OverrideRule)

func DumpAdminCache(part AdminCachePart, rebuild bool)
func AddCommandOverride(cmd string, overtype OverrideType, flags int)
func GetCommandOverride(cmd string, overtype OverrideType, flags *int) bool
func UnsetCommandOverride(cmd string, overtype OverrideType)

func CreateAdmGroup(group_name string) GroupId
func FindAdmGroup(group_name string) GroupId
func SetAdmGroupAddFlag(id GroupId, flag AdminFlag, enabled bool)
func GetAdmGroupAddFlag(id GroupId, flag AdminFlag) bool
func GetAdmGroupAddFlags(id GroupId) int
func SetAdmGroupImmunity(id GroupId, imtype ImmunityType, enabled bool)
func GetAdmGroupImmunity(id GroupId, imtype ImmunityType) bool
func SetAdmGroupImmuneFrom(id, other_id GroupId)
func GetAdmGroupImmuneCount(id GroupId) int
func GetAdmGroupImmuneFrom(id GroupId, number int) GroupId
func AddAdmGroupCmdOverride(id GroupId, name string, ovrtype OverrideType, rule OverrideRule)
func GetAdmGroupCmdOverride(id GroupId, name string, ovrtype OverrideType, rule *OverrideRule) bool
func RegisterAuthIdentType(name string)

func CreateAdmin(name string) AdminId
func GetAdminUsername(id AdminId, name []char, maxlength int) int
func BindAdminIdentity(id AdminId, auth, ident string) bool
func SetAdminFlag(id AdminId, flag AdminFlag, enabled bool)
func GetAdminFlag(id AdminId, flag AdminFlag, mode AdmAccessMode) bool
func GetAdminFlags(id AdminId, mode AdmAccessMode) int
func AdminInheritGroup(id AdminId, gid GroupId) bool
func GetAdminGroupCount(id AdminId) int
func GetAdminGroup(id AdminId, index int, name []char, maxlength int) GroupId
func SetAdminPassword(id AdminId, password string)
func GetAdminPassword(id AdminId, buffer []char, maxlength int) bool
func FindAdminByIdentity(auth, identity string) AdminId
func RemoveAdmin(id AdminId) bool

func FlagBitsToBitArray(bits int, array []bool, maxSize int) int
func FlagBitArrayToBits(array []bool, maxSize int) int
func FlagArrayToBits(array []AdminFlag, numFlags int) int
func FlagBitsToArray(bits int, array []AdminFlag, maxSize int) int
func FindFlagByName(name string, flag *AdminFlag) bool
func FindFlagByChar(c int, flag *AdminFlag) bool
func FindFlagChar(flag AdminFlag, c *int) bool
func ReadFlagString(flags string, numchars *int) int
func CanAdminTarget(admin, target AdminId) bool
func SetAdmGroupImmunityLevel(gid GroupId, level int) int
func GetAdmGroupImmunityLevel(gid GroupId) int
func SetAdminImmunityLevel(id AdminId, level int) int
func GetAdminImmunityLevel(id AdminId) int
func FlagToBit(flag AdminFlag) int
func BitToFlag(bit int, flag *AdminFlag) bool