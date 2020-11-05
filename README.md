# Go2SourcePawn
v0.22a

## Introduction

**Go2SourcePawn** is a transpiler that transforms a subset of Golang code to equivalent SourcePawn. The rationale behind Go2SourcePawn is to automate as much of the boilerplate possible when creating SourcePawn plugins.


### Features

* Abstracted common types into their own type classes such as `float[3]` aliasing as `Vec3`, etc.
Here is the current types list and what it abstracts to:
```
int    => int
float  => float
bool   => bool
int8   => char
string => char[]
Vec3   => float[3]
Map    => StringMap
Array  => ArrayList
```

* `#include <sourcemod>` and semicolon & new decls pragma are automatically added to each generated sourcepawn code file.

* Sourcemod-based constants, functions, and types are automatically added as part of the universal Go scope.

* Patterned matching where...

`string` matches to `const char[]`
and `*string` matches to `char[]`.
array/slice types will be automatically const unless passed by reference like: `*[]type`

so giving something like `[]int` will be `const int[]` while `*[]int` will become `int[]`.


* relative file imports are handled by using a dot `.` as the first letter
```go
import ".file"
```

Becomes:
```c
#include "file"
```

* Multiple return values are supported by transpiling them into variable references.
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

### Planned Features
* Generate Natives and an include file for them.
* Abstract, type-based syntax translation for higher data types like `StringMap` and `ArrayList`.
* Abstract function pointers from manual Function API calling.
* Abstract anonymous functions into name-generated functions. (perfect for abstracting timers)
* Handle-based Data Structures are abstracted into supportive syntax such where it's `value = Map["key"]` instead of `map.GetValue("key", value);`
* Func Methods for Entities and Vectors.

### Goal
Generate SourcePawn source code that is compileable by `spcomp` without having to modify/assist the generate source code.


## Contributing

To submit a patch, file an issue and/or hit up a pull request.

## Help

Commandline options:
* `--debug`, `-dbg` - prints the file's modified AST and pretty-printed version to a file for later checking.

* `--force`, `-f` - forcefully generates a SourcePawn source code file, even if errors/issues occurred during transpilation.

* `--help`, `-h` - Prints help list.

* `--version`, `-v` - Prints the version of SourceGo.

If you need help or have any question, simply file an issue with **\[HELP\]** in the title.


## Installation

### Requirements
Latest Golang version.

## Credits

* Nergal - main dev.

## License
This project is licensed under MIT License.