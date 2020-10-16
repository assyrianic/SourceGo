# Go2SourcePawn

## Introduction

**Go2SourcePawn** is a transpiler that transforms a subset of Golang-like code to equivalent SourcePawn. The rationale behind Go2SourcePawn is to automate as much of the boilerplate possible when creating SourcePawn plugins.


### Features

* Abstracted common types into their own type classes such as `float[3]` aliasing as `vec3`, `Entity` aliasing as `int`, etc.
* Data structure Handle types are abstracted into supportive syntax such as `map["example"]` instead of `map.GetValue("example", variable);`
Here is the current types list and what it abstracts to:
```
int    => int
Entity => int
float  => float
bool   => float
vec3   => float[3]
map    => StringMap
array  => ArrayList
obj    => Handle
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
