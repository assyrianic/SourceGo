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
	
	ClientInfo struct {
		Clients [2][MAXPLAYERS+1]Entity
	}
	
	Kektus    func(i, x Vec3, b string, blocks *Name, KC *int)   Handle
	EventFunc func(event Event, name string, dontBroadcast bool) Action
)

func (pi PlayerInfo) GetOrigin(buffer *Vec3) (float,float,float) {
	*buffer = pi.Origin
	return pi.Origin[0], pi.Origin[1], pi.Origin[2] 
}

func main() {
	//var ocpis func(client Entity) = OnClientPutInServer /// => Function ocpis = OnClientPutInServer;
	
	var cinfo ClientInfo
	for _, p1 := range cinfo.Clients {
		for _, x1 := range p1 {
		}
	}
	for _, p2 := range cinfo.Clients {
		for _, x2 := range p2 {
		}
	}
	var origin Vec3
	var p PlayerInfo
	x,y,z := p.GetOrigin(&origin)
	
	is_in_game1 := IsClientInGame(5)
	
	/*for i := 1; i<=MaxClients; i++ {
		is_in_game := IsClientInGame(i)
			//ocpis(i) /// becomes:
			/// Call_StartFunction(null, ocpis);
			/// Call_PushCell(i);
			/// Call_Finish();
	}*/
	
	switch x {
		case 1, 2:
		case 3:
		default:
	}
	
	switch {
		case x < 10, x+y < 10.0:
		case x * y <= 1024.0:
		default:
	}
}