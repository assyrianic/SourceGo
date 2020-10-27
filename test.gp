package main


const (
	a, b = "A", MAXPLAYERS
	c = a
	d string = "D"
	e = "e1"
	f = 1.00
)

type PlayerList [MAXPLAYERS]Entity

type EStruct struct {
	Clients PlayerList
}