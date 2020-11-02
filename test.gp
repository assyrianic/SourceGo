package main

/*
const (
	a, b = "A", MAXPLAYERS
	c = a
	d string = "D"
	e = "e1"
	f = 1.00
)
*/

var (
	myself = Plugin{
		name: "SrcGo Plugin",
		author: "Nergal",
		description: "Plugin made into SP from SrcGo.",
		version: "1.0a",
		url: "https://github.com/assyrianic/Go2SourcePawn",
	}
	t = [...]string{
		"kek",
		"foo",
		"bar",
		"bazz",
	}
	
	a int
	b Handle = nil
	c = 50.0
)

/*
func NBC() Entity

func main() {
	var pf func() Entity = NBC
	pf()
}
*/
/*
type (
	Point struct{ x, y float }
	Points struct { p [3*99]Point }
	polar [3*5]Points
	QAngle Vec3
	Name [64]char
	Color = [4]int
)

type PlayerInfo struct {
	K polar
	P QAngle
	M [3][10]int
}
*/

/*
type Kektus func(i, x Vec3, b string, blocks *[64]char, KC *int) Handle
type EventFunc func(event Event, name string, dontBroadcast bool) Action

func GetNil() Handle {
	return nil
}
*/