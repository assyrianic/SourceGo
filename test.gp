package main

/** Golang struct, limitations:
 *  No pointers, empty array brackets.
 */
type EStruct struct {
	Clients [MAXPLAYERS]Entity
}

func (g *EStruct) find_client(client Entity) Entity {
	for i:=1; i<MaxClients; i++ {
		if g.Clients[i]==int(client) {
			return Entity(i)
		}
	}
	return 0
}

func h() (int, int, float)
func f() (int, int, float)

func b(x, n int, p *int) int {
	a, b, c := f()
	a, b, c = h()
	return a + b + int(c)
}


func GetEntity() (Handle, Entity)
func mul_return() (Action, Map, Handle) {
	b, c := GetEntity()
	xx := Map(1000)
	return Action(0), xx, b + Handle(c)
}