# Go2SourcePawn
v1.4 beta

## Introduction

**Go2SourcePawn** is a transpiler that transforms a subset of Golang code to equivalent SourcePawn. The rationale behind Go2SourcePawn is to automate as much of the boilerplate possible when creating SourcePawn plugins.

### Purpose

To increase SourcePawn plugin development time by using Golang's streamline engineered syntax.


### Features

* Abstracted common types into their own type classes such as `float[3]` aliasing as `Vec3`, etc.
Here is the current types list and what params it abstracts to:
```
int     => int
float   => float
bool    => bool
*[]char => char[]
string  => const char[]
*string => char[]
Vec3    => const float[3]
*Vec3   => float[3]
```

* Pattern matching

Array/slice types will be automatically const unless passed by reference like: `*[]type`.
So having parameters such as `[]int` will be generated as `const int[]` while `*[]int` will be generated as `int[]`.


* relative file imports are handled by using a dot `.` as the first letter
```go
import ".file"
```

Becomes:
```c
#include "file"
```

* Multiple return values are supported by mutating them into variable references.

* Range loops for arrays:
```go
var players [MAXPLAYERS+1]Entity
for index, player := range players {
	/// code;
}
```

* Switch statements with and without an expression.
```go
switch x {
	case 1, 2:
	default:
}

switch {
	case x < 10, x+y < 10.0:
		
	default:
}
```

Expression-less switchs are useful for a more compact if-else-if series.


* Function pointer calls are broken down into manual Function API calling:
```go
func main() {
	CB := OnClientPutInServer
	for i := 1; i<=MaxClients; i++ {
		CB(i)
	}
}

func OnClientPutInServer(client Entity) {}
```
Becomes:
```c
public void OnPluginStart() {
	Function CB = OnClientPutInServer;
	for (int i = 1; i <= MaxClients; i++) {
		Call_StartFunction(null, CB);
		Call_PushCell(i);
		Call_Finish();
	}
}

public void OnClientPutInServer(int client) {}
```

* Anonymous Functions (aka Function Literals) are supported:
```go
my_timer := CreateTimer(2.0, func(timer Timer, data any) Action {
	return Plugin_Continue
}, 0, TIMER_REPEAT)
```

* Inline SourcePawn code using the builtin function `__sp__` - for those parts of SourcePawn that just can't be generated (like using new or making a methodmap from scratch).

`__sp__` only takes a single string of raw SourcePawn code. Optionally, you can also use a named string constant (it will be generated into the resulting code file, so keep that in mind.)
```go
/// using raw string quotes here so that single & double quotes don't have to be escaped.
var kv KeyValues
__sp__(`kv = new KeyValues("key_value", "key", "val");`)

...
__sp__(`delete kv;`)
```

* `make` for _making_ dynamically sized, local arrays:
```go
my_str := make([]char, size)
```
becomes:
```c
char[] my_str = new char[size];
```


### Planned Features
* Generate Natives and Forwards with an include file for them.
* Abstract, type-based syntax translation for higher data types like `StringMap` and `ArrayList`.
* Handle-based Data Structures are abstracted into supportive syntax such where it's `value = Map["key"]` instead of `map.GetValue("key", value);`

### Goal
Generate SourcePawn source code that is compileable by `spcomp` without having to modify/assist the generate source code.


## Contributing

To submit a patch, file an issue and/or hit up a pull request.

## Help

Command line options:
* `--debug`, `-dbg` - prints the file's modified AST and pretty-printed version to a file for later checking.

* `--force`, `-f` - forcefully generates a SourcePawn source code file, even if errors/issues occurred during transpilation.

* `--help`, `-h` - Prints help list.

* `--version`, `-v` - Prints the version of SourceGo.

* `--no-spcomp`, `-n` - Generates a SourcePawn source-code file without trying to invoke the SourcePawn compiler.

* `--verbose`, `-v` - prints additional warnings.

If you need help or have any question, simply file an issue with **\[HELP\]** in the title.


## Installation

### Requirements
Latest Golang version.

## Credits

* Nergal - main dev.

## License
This project is licensed under MIT License.
