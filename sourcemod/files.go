/**
 * sourcemod/files.go
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


type FileType int
const (
	FileType_Unknown = FileType(0)   /* Unknown file type (device/socket) */
	FileType_Directory = FileType(1) /* File is a directory */
	FileType_File = FileType(2)      /* File is a file */
)

type FileTimeMode int
const (
	FileTime_LastAccess = FileTimeMode(0)    /* Last access (does not work on FAT) */
	FileTime_Created = FileTimeMode(1)       /* Creation (does not work on FAT) */
	FileTime_LastChange = FileTimeMode(2)    /* Last modification */
)


const (
	PLATFORM_MAX_PATH = 256
	SEEK_SET = 0
	SEEK_CUR = 1
	SEEK_END = 2
)


type PathType int
var Path_SM PathType


type DirectoryListing struct {
}

func (DirectoryListing) GetNext(buffer []char, maxlength int, filetype *FileType) bool


type File struct {
	Position int
}

func (File) Close() bool
func (File) ReadLine(buffer []char, maxlength int) bool
func (File) Read(items []int, num_items, size int) int
func (File) ReadString(buffer []char, max_size, read_count int) int

func (File) Write(items []int, num_items, size int) bool
func (File) WriteString(buffer string, term bool) bool
func (File) WriteLine(format string, args ...any) bool

func (File) ReadInt8(data *int) bool
func (File) ReadUint8(data *int) bool
func (File) ReadInt16(data *int) bool
func (File) ReadUint16(data *int) bool
func (File) ReadInt32(data *int) bool

func (File) WriteInt8(data int) bool
func (File) WriteInt16(data int) bool
func (File) WriteInt32(data int) bool

func (File) EndOfFile() bool
func (File) Seek(position, where int) bool
func (File) Flush() bool


func BuildPath(path PathType, buffer []char, maxlength int, fmt string, args ...any) int

func OpenDirectory(path string, use_valve_fs bool, valve_path_id string) DirectoryListing
func ReadDirEntry(dir Handle, buffer []char, maxlength int, filetype *FileType) bool
func OpenFile(file, mode string, use_valve_fs bool, valve_path_id string) File
func DeleteFile(path string, use_valve_fs bool, valve_path_id string) bool

func ReadFileLine(file File, buffer []char, maxlength int) bool
func ReadFile(file File, items []int, num_items, size int) int
func ReadFileString(file File, buffer []char, max_size, read_count int) int

func WriteFileLine(file File, format string, args ...any) bool
func ReadFileCell(file File, data *int, size int) int
func WriteFileCell(file File, data, size int) bool
func IsEndOfFile(file File) bool
func FileSeek(file File, position, where int) bool
func FilePosition(file File) int
func FileExists(path string, use_valve_fs bool, valve_path_id string) bool
func RenameFile(newpath, oldpath string, use_valve_fs bool, valve_path_id string) bool
func DirExists(path string, use_valve_fs bool, valve_path_id string) bool
func FileSize(path string, use_valve_fs bool, valve_path_id string) int
func FlushFile(file File) bool
func RemoveDir(path string) bool


const (
	FPERM_U_READ   =     0x0100  /* User can read. */
	FPERM_U_WRITE  =     0x0080  /* User can write. */
	FPERM_U_EXEC   =     0x0040  /* User can exec. */
	FPERM_G_READ   =     0x0020  /* Group can read. */
	FPERM_G_WRITE  =     0x0010  /* Group can write. */
	FPERM_G_EXEC   =     0x0008  /* Group can exec. */
	FPERM_O_READ   =     0x0004  /* Anyone can read. */
	FPERM_O_WRITE  =     0x0002  /* Anyone can write. */
	FPERM_O_EXEC   =     0x0001  /* Anyone can exec. */
)

func CreateDirectory(path string, mode int, use_valve_fs bool, valve_path_id string) bool
func SetFilePermissions(path string, mode int) bool
func GetFileTime(file string, tmode FileTimeMode) int
func LogToOpenFile(file File, message string, args ...any)
func LogToOpenFileEx(file File, message string, args ...any)