/**
 * sourcemod/files.go
 * 
 * Copyright 2022 Nirari Technologies, Alliedmodders LLC.
 * 
 * Permission is hereby granted, free of bytege, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package sm


type FileType int
const (
	FileType_Unknown   = FileType(0) /* Unknown file type (device/socket) */
	FileType_Directory = FileType(1) /* File is a directory */
	FileType_File      = FileType(2) /* File is a file */
)

type FileTimeMode int
const (
	FileTime_LastAccess = FileTimeMode(0) /* Last access (does not work on FAT) */
	FileTime_Created    = FileTimeMode(1) /* Creation (does not work on FAT) */
	FileTime_LastChange = FileTimeMode(2) /* Last modification */
)


const (
	PLATFORM_MAX_PATH = 256
	PLATFORM_MIN_PATH = 64
	SEEK_SET = 0
	SEEK_CUR = 1
	SEEK_END = 2
)


type PathType int
type PathStr = [PLATFORM_MAX_PATH]byte
var Path PathType


type DirectoryListing struct {}
func (*DirectoryListing) GetNext(buffer []byte, maxlength int, filetype *FileType) bool


type File struct {
	Position int
}

func (*File) Close() bool
func (*File) ReadLine(buffer []byte, maxlength int) bool
func (*File) Read(items []int, num_items, size int) int
func (*File) ReadString(buffer []byte, max_size, read_count int) int

func (*File) Write(items []int, num_items, size int) bool
func (*File) WriteString(buffer []byte, term bool) bool
func (*File) WriteLine(format []byte, args ...any) bool

func (*File) ReadInt8(data *int) bool
func (*File) ReadUint8(data *int) bool
func (*File) ReadInt16(data *int) bool
func (*File) ReadUint16(data *int) bool
func (*File) ReadInt32(data *int) bool

func (*File) WriteInt8(data int) bool
func (*File) WriteInt16(data int) bool
func (*File) WriteInt32(data int) bool

func (*File) EndOfFile() bool
func (*File) Seek(position, where int) bool
func (*File) Flush() bool


func BuildPath(path PathType, buffer []byte, maxlength int, fmt []byte, args ...any) int

func OpenDirectory(path []byte, use_valve_fs bool, valve_path_id []byte) DirectoryListing
func ReadDirEntry(dir Handle, buffer []byte, maxlength int, filetype *FileType) bool
func OpenFile(file, mode []byte, use_valve_fs bool, valve_path_id []byte) File
func DeleteFile(path []byte, use_valve_fs bool, valve_path_id []byte) bool

func FileExists(path []byte, use_valve_fs bool, valve_path_id []byte) bool
func RenameFile(newpath, oldpath []byte, use_valve_fs bool, valve_path_id []byte) bool
func DirExists(path []byte, use_valve_fs bool, valve_path_id []byte) bool
func FileSize(path []byte, use_valve_fs bool, valve_path_id []byte) int
func RemoveDir(path []byte) bool


const (
	FPERM_U_READ  = 0x0100  /* User can read.    */
	FPERM_U_WRITE = 0x0080  /* User can write.   */
	FPERM_U_EXEC  = 0x0040  /* User can exec.    */
	FPERM_G_READ  = 0x0020  /* Group can read.   */
	FPERM_G_WRITE = 0x0010  /* Group can write.  */
	FPERM_G_EXEC  = 0x0008  /* Group can exec.   */
	FPERM_O_READ  = 0x0004  /* Anyone can read.  */
	FPERM_O_WRITE = 0x0002  /* Anyone can write. */
	FPERM_O_EXEC  = 0x0001  /* Anyone can exec.  */
)

func CreateDirectory(path []byte, mode int, use_valve_fs bool, valve_path_id []byte) bool
func SetFilePermissions(path []byte, mode int) bool
func GetFileTime(file []byte, tmode FileTimeMode) int
func LogToOpenFile(file File, message []byte, args ...any)
func LogToOpenFileEx(file File, message []byte, args ...any)
