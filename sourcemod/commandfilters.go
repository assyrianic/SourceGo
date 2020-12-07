/**
 * sourcemod/commandfilters.go
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


const (
	MAX_TARGET_LENGTH =           64

	COMMAND_FILTER_ALIVE =        (1<<0)      /**< Only allow alive players */
	COMMAND_FILTER_DEAD =         (1<<1)      /**< Only filter dead players */
	COMMAND_FILTER_CONNECTED =    (1<<2)      /**< Allow players not fully in-game */
	COMMAND_FILTER_NO_IMMUNITY =  (1<<3)      /**< Ignore immunity rules */
	COMMAND_FILTER_NO_MULTI =     (1<<4)      /**< Do not allow multiple target patterns */
	COMMAND_FILTER_NO_BOTS =      (1<<5)      /**< Do not allow bots to be targetted */

	COMMAND_TARGET_NONE =          0          /**< No target was found */
	COMMAND_TARGET_NOT_ALIVE =    -1          /**< Single client is not alive */
	COMMAND_TARGET_NOT_DEAD =     -2          /**< Single client is not dead */
	COMMAND_TARGET_NOT_IN_GAME =  -3          /**< Single client is not in game */
	COMMAND_TARGET_IMMUNE =       -4          /**< Single client is immune */
	COMMAND_TARGET_EMPTY_FILTER = -5          /**< A multi-filter (such as @all) had no targets */
	COMMAND_TARGET_NOT_HUMAN =    -6          /**< Target was not human */
	COMMAND_TARGET_AMBIGUOUS =    -7          /**< Partial name had too many targets */
)


func ProcessTargetString(pattern string, admin int, targets []int, max_targets, filter_flags int, target_name []char, tn_maxlength int, tn_is_ml *bool) int
func ReplyToTargetError(client, reason int)


type MultiTargetFilter func(pattern string, clients ArrayList) bool

func AddMultiTargetFilter(pattern string, filter MultiTargetFilter, phrase string, phraseIsML bool)
func RemoveMultiTargetFilter(pattern string, filter MultiTargetFilter)