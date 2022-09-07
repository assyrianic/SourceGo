package main

import (
	"./sourcemod"
)

var (
	myself = sm.Plugin{
		Name:     []byte("SrcGo Plugin"),
		Author:   []byte("Nergal"),
		Descript: []byte("Plugin made into SP from SrcGo."),
		Version:  []byte("1.0a"),
		Url:      []byte("https://github.com/assyrianic/SourceGo"),
	}
	str_array = [...]string{
		"kek",
		"foo",
		"bar",
		"bazz",
	}
)

const MAXPLAYERS = 34

const (
	a, b = "A", MAXPLAYERS
	c = a
	d = "D"
	e = "e1"
	f = 1.00
	MakeStrMap = "StringMap smap = new StringMap();"
)

type (
	Point struct{ x, y float32 }
	QAngle = [3]float32
	Name   = [64]byte
	Color  = [0x4]int
	
	PlayerInfo struct {
		Origin sm.Vec3
		Angle sm.QAngle
		Weaps [3]sm.Entity
		PutInServer func(client sm.Entity)
	}
	
	ClientInfo struct {
		Clients [2][MAXPLAYERS+1]sm.Entity
	}
	
	Kektus    func(i, x sm.Vec3, b []byte, blocks *Name, KC *int) sm.Handle
	EventFunc func(event *sm.Event, name []byte, dontBroadcast bool) sm.Action
	VecFunc   func(vec sm.Vec3) (float32, float32, float32)
)

func V() (float32, float32, float32) {
	return 1.0, 2.0, 3.0
}

func V3() sm.Vec3 {
	return sm.Vec3{1.0, 2.0, 3.0}
}

func main() {
	t, x := func(i int) int {
		return 10 + i
	}, sm.VerifyCoreVersion()
	if x > t(5) {
		func() int {
			return 5
		}()
		return
	}
}