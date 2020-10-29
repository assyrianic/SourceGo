package main


const (
	a, b = "A", MAXPLAYERS
	c = a
	d string = "D"
	e = "e1"
	f = 1.00
)

type Name [64]char
type Color = [4]int


type (
	Point struct{ x, y float }
	Points struct { p [3]Point }
	polar [3]Points
)

type PlayerInfo struct {
	K polar
}

func main() {
	
}