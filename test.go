package main

import (
	"sourcemod"
	"sdktools"
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
	MakeStrMap = "StringMap smap = new StringMap();"
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
	EventFunc func(event *Event, name string, dontBroadcast bool)           Action
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
	
	var k,l int
	k &^= l
	
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
		return Plugin_Continue
	}, 0, 0)
	
	inlined_call_res := func(a,b int) int {
		return a + b
	}(1, 2)
	
	inlined_call_res1,inlined_call_res2 := func(a,b int) (int,int) {
		return a + b, a*b
	}(1, 2)
	
	caller := func(a,b int) int {
		return a + b
	}
	//n := caller(1, 2)
	__sp__(`
	int n;
	Call_StartFunction(null, caller);
	Call_PushCell(1); Call_PushCell(2);
	Call_Finish(n);`)
	
	var kv KeyValues
	/// using raw string quotes so we don't have to escape double quotes.
	__sp__(`kv = new KeyValues("kek1", "kek_key", "kek_val");
	delete kv;`)
	
	AddMultiTargetFilter("@!party", func(pattern string, clients ArrayList) bool {
		non := StrContains(pattern, "!", false) != -1
		for i:=MAX_TF_PLAYERS; i > 0; i-- {
		__sp__(`if( IsClientValid(i) && clients.FindValue(i) == -1 ) {
			if( g_cvars.enabled.BoolValue && g_dnd.IsGameMaster(i) ) {
				if( !non ) {
					clients.Push(i);
				}
			} else if( non ) {
				clients.Push(i);
			}
		}`)
		}
		return true
	}, "The D&D Quest Party", false)
	
	__sp__(MakeStrMap)
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

func KeyValuesToStringMap(kv KeyValues, stringmap map[string][]char, hide_top bool, depth int, prefix *[]char) {
	type SectStr = [128]char
	for {
		var section_name SectStr
		kv.GetSectionName(section_name, len(section_name))
		if kv.GotoFirstSubKey(false) {
			var new_prefix SectStr
			switch {
				case depth==0 && hide_top:
					new_prefix = ""
				case prefix[0] == 0:
					new_prefix = section_name
				default:
					FormatEx(new_prefix, len(new_prefix), "%s.%s", prefix, section_name)
			}
			KeyValuesToStringMap(kv, stringmap, hide_top, depth+1, new_prefix)
            kv.GoBack()
		} else {
			if kv.GetDataType(NULL_STRING) != KvData_None {
				var key SectStr
				if prefix[0] == 0 {
					key = section_name
				} else {
					FormatEx(key, len(key), "%s.%s", prefix, section_name)
				}
				
				// lowercaseify the key
				keylen := strlen(key)
				for i := 0; i < keylen; i++ {
					bytes := IsCharMB(key[i])
					if bytes==0 {
						key[i] = CharToLower(key[i])
					} else {
						i += (bytes - 1)
					}
				}
				
				var value SectStr
				kv.GetString(NULL_STRING, value, len(value), NULL_STRING)
				stringmap[key] = value
			}
		}
		
		if !kv.GotoNextKey(false) {
			break
		}
	}
}