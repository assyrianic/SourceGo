/**
 * sourcemod/core.go
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
	"sourcemod/types"
)


type Plugin struct {
	name, description, author, version, url string
}

type Extension struct {
	name, file string
	autoload, required int
}

const (
	SOURCEMOD_PLUGINAPI_VERSION = 5
	
	SOURCEMOD_V_TAG string = "manual"
	SOURCEMOD_V_REV = 0
	SOURCEMOD_V_CSET string = "0"
	SOURCEMOD_V_MAJOR   = 1               /**< SourceMod Major version */
	SOURCEMOD_V_MINOR   = 10              /**< SourceMod Minor version */
	SOURCEMOD_V_RELEASE = 0               /**< SourceMod Release version */
	
	SOURCEMOD_VERSION string = "1.10.0-manual" /**< SourceMod version string (major.minor.release-tag) */
)


type Action int
const (
	Plugin_Continue = Action(iota)
	Plugin_Changed
	Plugin_Handled
	Plugin_Stop
)

type Identity int
const (
	Identity_Core = Identity(iota)
	Identity_Extension
	Identity_Plugin
)


type PluginStatus int
const (
	Plugin_Running = PluginStatus(iota)       /**< Plugin is running */
	/** All states below are "temporarily" unexecutable */
	Plugin_Paused          /**< Plugin is loaded but paused */
	Plugin_Error           /**< Plugin is loaded but errored/locked */
	/** All states below do not have all natives */
	Plugin_Loaded          /**< Plugin has passed loading and can be finalized */
	Plugin_Failed          /**< Plugin has a fatal failure */
	Plugin_Created         /**< Plugin is created but not initialized */
	Plugin_Uncompiled      /**< Plugin is not yet compiled by the JIT */
	Plugin_BadLoad         /**< Plugin failed to load */
	Plugin_Evicted          /**< Plugin was unloaded due to an error */
)


type PluginInfo int
const (
	PlInfo_Name = PluginInfo(iota)            /**< Plugin name */
	PlInfo_Author          /**< Plugin author */
	PlInfo_Description     /**< Plugin description */
	PlInfo_Version         /**< Plugin version */
	PlInfo_URL 
)

const NULL_STRING string = ""
var (
	NULL_VECTOR Vec3    /// can't put NULL_VECTOR as 'const'.
	INVALID_FUNCTION __function__
)


func IsNullVector(vec Vec3) bool
func IsNullString(str string) bool
func VerifyCoreVersion() int
func MarkNativeAsOptional(name string)