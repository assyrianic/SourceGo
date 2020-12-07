/**
 * sourcemod/entity.go
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


type PropType int
const (
	Prop_Send = PropType(0)  /**< This property is networked. */
	Prop_Data = PropType(1)   /**< This property is for save game data fields. */
)


const (
	FL_EDICT_CHANGED =                (1<<0)  /**< Game DLL sets this when the entity state changes
                                                     Mutually exclusive with FL_EDICT_PARTIAL_CHANGE. */
	FL_EDICT_FREE =                   (1<<1)  /**< this edict if free for reuse */
	FL_EDICT_FULL =                   (1<<2)  /**< this is a full server entity */
	FL_EDICT_FULLCHECK =              (0<<0)  /**< call ShouldTransmit() each time, this is a fake flag */
	FL_EDICT_ALWAYS =                 (1<<3)  /**< always transmit this entity */
	FL_EDICT_DONTSEND =               (1<<4)  /**< don't transmit this entity */
	FL_EDICT_PVSCHECK =               (1<<5)  /**< always transmit entity, but cull against PVS */
	FL_EDICT_PENDING_DORMANT_CHECK =  (1<<6)
	FL_EDICT_DIRTY_PVS_INFORMATION =  (1<<7)
	FL_FULL_EDICT_CHANGED =           (1<<8)
	MAXENTS = 2048
)


type PropFieldType int
const (
	PropField_Unsupported = PropFieldType(0)      /**< The type is unsupported. */
	PropField_Integer          /**< Valid for SendProp and Data fields */
	PropField_Float            /**< Valid for SendProp and Data fields */
	PropField_Entity           /**< Valid for Data fields only (SendProp shows as int) */
	PropField_Vector           /**< Valid for SendProp and Data fields */
	PropField_String           /**< Valid for SendProp and Data fields */
	PropField_String_T         /**< Valid for Data fields.  Read only.
	                                 Note that the size of a string_t is dynamic, and
	                                 thus FindDataMapOffs() will return the constant size
	                                 of the string_t container (which is 32 bits right now). */
	PropField_Variant           /**< Valid for Data fields only Type is not known at the field level,
                                     (for this call), but dependent on current field value. */
)


func GetMaxEntities() int
func GetEntityCount() int
func IsValidEntity(entity Entity) bool
func IsValidEdict(edict int) bool
func IsEntNetworkable(entity Entity) bool
func CreateEdict() Entity
func RemoveEdict(edict int)
func RemoveEntity(entity Entity)
func GetEdictFlags(edict int) int
func SetEdictFlags(edict, flags int)
func GetEdictClassname(edict int, clsname []char, maxlength int) bool
func GetEntityNetClass(edict int, clsname []char, maxlength int) bool
func ChangeEdictState(edict, offset int)
func GetEntData(entity, offset, size int) int
func SetEntData(entity, offset int, value any, size int, changeState bool)
func GetEntDataFloat(entity, offset int) float
func SetEntDataFloat(entity, offset int, value float, changeState bool)
func GetEntDataEnt2(entity, offset int) Entity
func SetEntDataEnt2(entity, offset, other int, changeState bool)
func GetEntDataVector(entity, offset int, vec *Vec3)
func SetEntDataVector(entity, offset int, vec Vec3, changeState bool)
func GetEntDataString(entity, offset int, buffer []char, maxlen int) int
func SetEntDataString(entity, offset int, buffer string, maxlen int, changeState bool) int
func FindSendPropInfo(cls, prop string, fieldtype *PropFieldType, num_bits, local_offset *int) int
func FindDataMapInfo(entity Entity, prop string, fieldtype *PropFieldType, num_bits, local_offset *int) int
func GetEntSendPropOffs(entity Entity, prop string, actual bool) int
func HasEntProp(entity Entity, prop_type PropType, prop string) bool
func GetEntProp(entity Entity, prop_type PropType, prop string, size, element int) int
func SetEntProp(entity Entity, prop_type PropType, prop string, value any, size, element int)
func GetEntPropFloat(entity Entity, prop_type PropType, prop string, element int) float
func SetEntPropFloat(entity Entity, prop_type PropType, prop string, value float, element int)
func GetEntPropEnt(entity Entity, prop_type PropType, prop string, element int) Entity
func SetEntPropEnt(entity Entity, prop_type PropType, prop string, other Entity, element int)
func GetEntPropVector(entity Entity, prop_type PropType, prop string, vec *Vec3, element int)
func SetEntPropVector(entity Entity, prop_type PropType, prop string, vec Vec3, element int)
func GetEntPropString(entity Entity, prop_type PropType, prop string, buffer []char, maxlen, element int) int
func SetEntPropString(entity Entity, prop_type PropType, prop, buffer string, element int) int
func GetEntPropArraySize(entity Entity, prop_type PropType, prop string) int
func GetEntDataArray(entity, offset int, array []int, arraySize, dataSize int)
func SetEntDataArray(entity, offset int, array []int, arraySize, dataSize int, changeState bool)
func GetEntityAddress(entity Entity) Address
func GetEntityClassname(entity Entity, clsname []char, maxlength int) bool