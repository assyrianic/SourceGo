# Go2SourcePawn
v0.3a

## Introduction

**Go2SourcePawn** is a transpiler that transforms a subset of Golang code to equivalent SourcePawn. The rationale behind Go2SourcePawn is to automate as much of the boilerplate possible when creating SourcePawn plugins.


### Features

* Abstracted common types into their own type classes such as `float[3]` aliasing as `Vec3`, etc.
* Data structure Handle types are abstracted into supportive syntax such as `value = Map["key"]` instead of `map.GetValue("key", value);`
Here is the current types list and what it abstracts to:
```
int    => int
float  => float
bool   => float
Vec3   => float[3]
Map    => StringMap
Array  => ArrayList
Obj    => Handle
```

* `#include <sourcemod>` and semicolon & new decls pragma are automatically added to each generated sourcepawn code file.

* relative file imports are handled by using a dot `.` as the first letter
```go
import ".file"
```

Becomes:
```c
#include "file"
```

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
