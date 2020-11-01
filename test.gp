package main


const (
	a, b = "A", MAXPLAYERS
	c = a
	d string = "D"
	e = "e1"
	f = 1.00
)

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
}

var (
	/// must be 'char t[][] = { "kek", "foo", "bar", "baz" };'
	t = [...]string{
		"kek",
		"foo",
		"bar",
		"baz",
	}
)


type Kektus func(i, x QAngle, b string, blocks *Name, KC *int) Handle
type EventFunc func(event Event, name string, dontBroadcast bool) Action

func main() {
	//var pf func() int
}