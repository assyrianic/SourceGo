# Go2SourcePawn
v0.12a

## Introduction

**Go2SourcePawn** is a transpiler that transforms a subset of Golang code to equivalent SourcePawn. The rationale behind Go2SourcePawn is to automate as much of the boilerplate possible when creating SourcePawn plugins.


### Features

* Abstracted common types into their own type classes such as `float[3]` aliasing as `Vec3`, etc.
* Data structure Handle types are abstracted into supportive syntax such as `value = Map["key"]` instead of `map.GetValue("key", value);`
Here is the current types list and what it abstracts to:
```
int    => int
float  => float
bool   => bool
int8   => char
string => char[]
Vec3, QAngle => float[3]
Map    => StringMap
Array  => ArrayList
```

* `#include <sourcemod>` and semicolon & new decls pragma are automatically added to each generated sourcepawn code file.

* Sourcemod-based constants, functions, and types are automatically added as part of the universal Go scope.

* relative file imports are handled by using a dot `.` as the first letter
```go
import ".file"
```

Becomes:
```c
#include "file"
```

### Planned Features
* Generate Natives and an include file for them.
* Abstract, type-based syntax translation for higher data types like `StringMap` and `ArrayList`.
* Patterned matching where...

`string` matches to `const char[]`
and `*string` matches to `char[]`.
array/slice types will be automatically const unless given by "pointer" like: `*[]type`

so giving something like `[]int` will be `const int[]` while `*[]int` will become `int[]`.

* Multiple return values are supported by transpiling them into variable references.
* Abstract function pointers.
* Abstraction anonymous functions. (perfect for abstracting timers)

### Goal
Generate SourcePawn source code that is compileable by `spcomp` without having to modify/assist the generate source code.


## Contributing

To submit a patch, file an issue and/or hit up a pull request.

## Help

If you need help or have any question, simply file an issue with **\[HELP\]** in the title.


## Installation

### Requirements
latest Golang version.

## Credits

* Nergal - main dev.

## License
This project is licensed under MIT License.
