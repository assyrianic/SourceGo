package main


type Database any
const (
	DBCB_Connect = 0
)

var hTheDB Database

func (Database) Connect(conn int, n string)

func main() {
	Database.Connect(DBCB_Connect, "maptracker")
}
