package main

const (
	a, b = "A", MAXPLAYERS
	c = a
	d string = "D"
	e = "e1"
	f = 1.00
)

var (
	myself = Plugin{
		name: "SrcGo Plugin",
		author: "Nergal",
		description: "Plugin made into SP from SrcGo.",
		version: "1.0a",
		url: "https://github.com/assyrianic/Go2SourcePawn",
	}
	
	str_array = [...]string{
		"kek",
		"foo",
		"bar",
		"bazz",
	}
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

func (pi PlayerInfo) GetOrigin(buffer *Vec3) {
	*buffer = pi.Origin
}

func IsClientInGame(client Entity) bool

func main() {
	var p PlayerInfo
	for i := 1; i<=MaxClients; i++ {
		if IsClientInGame(i) {
			OnClientPutInServer(i)
		}
		var origin Vec3
		p.GetOrigin(&origin)
		PrintToServer("%f, %f, %f", origin[0], origin[1], origin[2])
	}
}

func OnClientPutInServer(client Entity) {
	/// do something with client.
}