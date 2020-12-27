/**
 * sourcemod.go
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
	"sourcemod/floats"
	"sourcemod/vector"
	"sourcemod/strings"
	"sourcemod/handles"
	"sourcemod/functions"
	"sourcemod/files"
	"sourcemod/logging"
	"sourcemod/timers"
	"sourcemod/admin"
	"sourcemod/keyvalues"
	"sourcemod/dbi"
	"sourcemod/lang"
	"sourcemod/sorting"
	"sourcemod/textparse"
	"sourcemod/clients"
	"sourcemod/console"
	"sourcemod/convars"
	"sourcemod/events"
	"sourcemod/bitbuffer"
	"sourcemod/protobuf"
	"sourcemod/usermessages"
	"sourcemod/menus"
	"sourcemod/halflife"
	"sourcemod/adt_array"
	"sourcemod/adt_trie"
	"sourcemod/adt_stack"
	"sourcemod/banning"
	"sourcemod/commandfilters"
	"sourcemod/nextmap"
	"sourcemod/commandline"
	"sourcemod/helpers"
	"sourcemod/entity"
	"sourcemod/entity_prop_stocks"
)


type APLRes int
const (
	APLRes_Success = APLRes(0)     /**< Plugin should load */
	APLRes_Failure         /**< Plugin shouldn't load and should display an error */
	APLRes_SilentFailure    /**< Plugin shouldn't load but do so silently */
)


type GameData Handle

func LoadGameConfigFile(file string) GameData
func (GameData) GetOffset(key string) int
func (GameData) GetKeyValue(key string, buffer []char, maxlen int) bool
func (GameData) GetAddress(name string) Address


func GetMyHandle() Handle
func GetPluginIterator() Handle
func MorePlugins(iter Handle) bool
func ReadPlugin(iter Handle) Handle
func GetPluginStatus(plugin Handle) PluginStatus
func GetPluginFilename(plugin Handle, buffer []char, maxlength int)
func IsPluginDebugging(plugin Handle) bool
func GetPluginInfo(plugin Handle, info PluginInfo, buffer []char, maxlength int) bool
func FindPluginByNumber(order_num int) Handle
func SetFailState(state string, args ...any)
func ThrowError(fmt string, args ...any)
func LogStackTrace(fmt string, args ...any)
func GetTime(bigStamp [2]int) int
func FormatTime(buffer []char, maxlength int, format string, stamp int)
func GetSysTickCount() int
func AutoExecConfig(autoCreate bool, name, folder string)
func RegPluginLibrary(name string)
func LibraryExists(name string) bool
func GetExtensionFileStatus(name string, err []char, maxlength int) int


const (
	MAPLIST_FLAG_MAPSFOLDER =    (1<<0)    /**< On failure, use all maps in the maps folder. */
	MAPLIST_FLAG_CLEARARRAY =    (1<<1)    /**< If an input array is specified, clear it before adding. */
	MAPLIST_FLAG_NO_DEFAULT =    (1<<2)    /**< Do not read "default" or "mapcyclefile" on failure. */
)

func ReadMapList(array ArrayList, serial *int, str string, flags int) ArrayList
func SetMapListCompatBind(name, file string)


type FeatureType int
const (
	/**
	 * A native function call.
	 */
	FeatureType_Native = FeatureType(0)

	/**
	 * A named capability. This is distinctly different from checking for a
	 * native, because the underlying functionality could be enabled on-demand
	 * to improve loading time. Thus a native may appear to exist, but it might
	 * be part of a set of features that are not compatible with the current game
	 * or version of SourceMod.
	 */
	FeatureType_Capability
)


type FeatureStatus int
const (
	/**
	 * Feature is available for use.
	 */
	FeatureStatus_Available = FeatureStatus(0)

	/**
	 * Feature is not available.
	 */
	FeatureStatus_Unavailable

	/**
	 * Feature is not known at all.
	 */
	FeatureStatus_Unknown
)

func CanTestFeatures() bool
func GetFeatureStatus(ftr_type FeatureType, name string) FeatureStatus
func RequireFeature(ftr_type FeatureType, name, fmt string, args ...any)


type NumberType int
const (
    NumberType_Int8 = NumberType(0)
    NumberType_Int16
    NumberType_Int32
)


type Address int
const Address_Null = Address(0)    /// a typical invalid result when an address lookup fails

func LoadFromAddress(addr Address, size NumberType) int
func StoreToAddress(addr Address, data int, size NumberType)


type FrameIterator struct {
	LineNumber int
}

func (FrameIterator) Next() bool
func (FrameIterator) Reset()
func (FrameIterator) GetFunctionName(buffer []char, maxlen int)
func (FrameIterator) GetFilePath(buffer []char, maxlen int)