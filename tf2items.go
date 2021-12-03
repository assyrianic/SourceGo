/**
 * tf2items.go
 * 
 * Copyright 2020 Nirari Technologies, Alliedmodders LLC.
 * All Credit to Asher 'Asherkin' Baker.
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

const (
	OVERRIDE_CLASSNAME =    (1 << 0)    // Item will override the entity's classname.
	OVERRIDE_ITEM_DEF =     (1 << 1)    // Item will override the item's definition index.
	OVERRIDE_ITEM_LEVEL =   (1 << 2)    // Item will override the entity's level.
	OVERRIDE_ITEM_QUALITY = (1 << 3)    // Item will override the entity's quality.
	OVERRIDE_ATTRIBUTES =   (1 << 4)    // Item will override the attributes for the item with the given ones.
	OVERRIDE_ALL =          (0b11111)   // Magically enables all the other flags.

	PRESERVE_ATTRIBUTES =   (1 << 5)
	FORCE_GENERATION =      (1 << 6)
)


func TF2Items_GiveNamedItem(Client Entity, hItem Handle) Entity
func TF2Items_CreateItem(iFlags int) Handle
func TF2Items_SetFlags(hItem Handle, iFlags int)
func TF2Items_SetClassname(hItemOverride Handle, strClassName string)
func TF2Items_SetItemIndex(hItem Handle, iItemDefIndex int)
func TF2Items_SetQuality(hItem Handle, iEntityQuality int)
func TF2Items_SetLevel(hItem Handle, iEntityLevel int)
func TF2Items_SetNumAttributes(hItem Handle, iNumAttributes int)
func TF2Items_SetAttribute(hItem Handle, iSlotIndex, iAttribDefIndex int, flValue float)
func TF2Items_GetFlags(hItem Handle) int
func TF2Items_GetClassname(hItem Handle, strDest []char, iDestSize int)
func TF2Items_GetItemIndex(hItem Handle) int
func TF2Items_GetQuality(hItem Handle) int
func TF2Items_GetLevel(hItem Handle) int
func TF2Items_GetNumAttributes(hItem Handle) int
func TF2Items_GetAttributeId(hItem Handle, iSlotIndex int) int
func TF2Items_GetAttributeValue(hItem Handle, iSlotIndex int) float