package main

import (
	"sourcemod"
)

var (
	myself = Plugin{
		name:        "SrcGo Plugin",
		author:      "Nergal",
		description: "Plugin made into SP from SrcGo.",
		version:     "1.0a",
		url:         "https://github.com/assyrianic/Go2SourcePawn",
	}
	str_array = [...]string{
		"kek",
		"foo",
		"bar",
		"bazz",
	}
)

const (
	a, b = "A", MAXPLAYERS
	c = a
	d string = "D"
	e = "e1"
	f = 1.00
)

type (
	Point struct{ x, y float }
	QAngle = Vec3
	Name   = [64]char
	Color  = [0x4]int
	
	PlayerInfo struct {
		Origin Vec3
		Angle QAngle
		Weaps [3]Entity
		PutInServer func(client Entity)
	}
	
	ClientInfo struct {
		Clients [2][MAXPLAYERS+1]Entity
	}
	
	Kektus    func(i, x  Vec3,  b    string, blocks        *Name, KC *int) Handle
	EventFunc func(event Event, name string, dontBroadcast bool)           Action
	VecFunc   func(vec Vec3) (float, float, float)
)

func (pi PlayerInfo) GetOrigin(buffer *Vec3) (float,float,float) {
	*buffer = pi.Origin
	return pi.Origin[0], pi.Origin[1], pi.Origin[2] 
}

func TestOrigin() (float, float, float) {
	var (
		pi PlayerInfo
		o Vec3
	)
	return pi.GetOrigin(&o)
}

var (
	ff1 func() int
	ff2 func() int
	ff3 func() float
)

func FF1() int
func FF2() int
func FF3() float
func FF4() (int, int, float)

func GG1() (int, int, float) {
	return FF1(), FF2(), FF3() 
}
func GG2() (int, int, float) {
	return ff1(), ff2(), ff3()
}
func GG3() (int, int, float) {
	return FF4()
}

func main() {
	var cinfo ClientInfo
	for _, p1 := range cinfo.Clients {
		for _, x1 := range p1 {
			is_in_game := IsClientInGame(x1)
		}
	}
	
	var p PlayerInfo
	var origin Vec3
	x,y,z := p.GetOrigin(&origin)
	
	//var k,l int
	//k &^= l
	
	CB := IndirectMultiRet
	p.PutInServer = func(client Entity) {
	}
	for i := 1; i<=MaxClients; i++ {
		p.PutInServer(i)
		CB()
		j,k,l := CB()
	}
	
	for f := 2.0; f < 100.0; f = Pow(f, 2.0) {
		PrintToServer("%0.2f", f)
	}
	
	my_timer := CreateTimer(0.1, func(timer Timer, data any) Action {
		data++
		return Plugin_Continue
	}, 0, 0)
	
	inlined_call_res := func(a,b int) int {
		return a + b
	}(1, 2)
	
	caller := func(a,b int) int {
		return a + b
	}
	n := caller(1, 2)
}

func IndirectMultiRet() (bool, bool, bool) {
	return MultiRetFn()
}

func MultiRetFn() (bool, bool, bool) {
	return true, false, true
}

func OnClientPutInServer(client Entity) {
}

func GetProjPosToScreen(client int, vecDelta Vec3) (xpos, ypos float) {
	var playerAngles, vecforward, right, up Vec3
	GetClientEyeAngles(client, playerAngles)
	
	up[2] = 1.0
	GetAngleVectors(playerAngles, &vecforward, &NULL_VECTOR, &NULL_VECTOR)
	vecforward[2] = 0.0
	
	NormalizeVector(vecforward, &vecforward)
	GetVectorCrossProduct(up, vecforward, &right)
	
	front, side := GetVectorDotProduct(vecDelta, vecforward), GetVectorDotProduct(vecDelta, right)
	
	xpos, ypos = 360.0 * -front, 360.0 * -side
	flRotation := (ArcTangent2(xpos, ypos) + FLOAT_PI) * (57.29577951)
	yawRadians := -flRotation * 0.017453293

	/// Rotate it around the circle
	xpos, ypos = ( 500 + (360.0 * Cosine(yawRadians)) ) / 1000.0, ( 500 - (360.0 * Sine(yawRadians)) ) / 1000.0
	return
}