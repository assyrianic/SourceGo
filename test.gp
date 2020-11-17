package main

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
	QAngle Vec3
	Name [64]char
	Color = [4]int
	
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

func GetFuncByName(name string) func(client Entity)

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
	
	CB := MultiRetFn
	p.PutInServer = OnClientPutInServer
	for i := 1; i<=MaxClients; i++ {
		p.PutInServer(i)
		CB()
		j,k,l := CB()
	}
	
	for f := 2.0; f < 100.0; f = Pow(f, 2.0) {
		PrintToServer("%0.2f", f)
	}
}

func MultiRetFn() (bool, bool, bool) {
	return true,false,true
}

func OnClientPutInServer(client Entity) {
}