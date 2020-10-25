package main

import "sourcemod"

/** Golang struct, limitations:
 *  No pointers, empty array brackets.
 */
type EStruct struct {
	Clients [35]int
	FPtr    func()int
}

func h() (int, int, float)
func f() (int, int, float) /// => func f(*int, *float32) int

func b(x, n int, p *int) int {
	a, b, c := f()
	a, b, c = h()
	return a + b + int(c)
}

func mul_return() (SM.Action, SM.StrMap, Handle) {
	return SM.Plugin_Continue, make(SM.StrMap), Handle(-1)
}

func kek() {
}