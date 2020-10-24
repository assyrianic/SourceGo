package main

/*
import (
	//"sourcemod" /// will be automatically inserted.
	"tf2_stocks"
	"sdkhooks"
	"arraylist"
	".Impacts_suggestion"
)
*/

func f() (int, int, int)

func Func1(x, n int, p *int8) int {
	/*
	switch x {
		case 1, 2, 3, 4:
			n &= 1
			n &^ x+1 << 5
		case test:
			n += test
		default:
			foo()
	}
	*/
	a, b, c := f()   /// transform into `var a,b,c type; a = f(&b, &c)`
	return a + b + c
}

//func (n int) Receiver(c float) Action