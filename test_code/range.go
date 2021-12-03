package main

import (
	"sourcemod"
)


func main() {
	var clients [MAXPLAYERS]int
	for i := range clients {
		if !IsValidClient(clients[i]) {
			continue
		}
	}
}
