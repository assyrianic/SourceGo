package main

var (
	myself = Plugin{
		Name: "SrcGo Plugin",
		Author: "Nergal",
		Description: "Plugin made into SP from SrcGo.",
		Version: "1.0a",
		Url: "https://github.com/assyrianic/Go2SourcePawn",
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
	}
	
	Kektus    func(i, x Vec3, b string, blocks *Name, KC *int)   Handle
	EventFunc func(event Event, name string, dontBroadcast bool) Action
)

func (pi PlayerInfo) GetOrigin(buffer *Vec3) (float,float,float) {
	*buffer = pi.Origin
	return pi.Origin[0], pi.Origin[1], pi.Origin[2] 
}

func IsClientInGame(client Entity) bool

func main() {
	//var ocpis func(client Entity) = OnClientPutInServer /// => Function ocpis = OnClientPutInServer;
	var clients [2][MAXPLAYERS+1]Entity
	for _, p1 := range clients {
		for _, x1 := range p1 {
		}
	}
	for _, p2 := range clients {
		for _, x2 := range p2 {
		}
	}
	var origin Vec3
	var p PlayerInfo
	x,y,z := p.GetOrigin(&origin)
	/*
	for i := 1; i<=MaxClients; i++ {
		if IsClientInGame(i) {
			
			//ocpis(i) /// becomes:
			/// Call_StartFunction(null, ocpis);
			/// Call_PushCell(i);
			/// Call_Finish();
		}
	}
	*/
}