package main


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


func OnlyScoutsLeft(team int) bool {
	for i:=0; i<=MaxClients; i++ {
		if !IsValidClient(i) || !IsPlayerAlive(i) {
			continue;
		} else if GetClientTeam(i) == team && TF2_GetPlayerClass(i) != TFClass_Scout {
			return false
		}
	}
	return true
}


func IsClientInGame(client Entity) bool

func main() {
	for i := 1; i<=MaxClients; i++ {
		if IsClientInGame(i) {
			OnClientPutInServer(i)
		}
	}
}

func OnClientPutInServer(client Entity) {
	/// do something with client.
}