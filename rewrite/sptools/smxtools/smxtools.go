package SMXTools

// saving this for the future.
// no affiliation with sourcepawn's smxtools.

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"strings"
	"bytes"
)


type (
	SmxNameTable struct {
		NameTable map[string]uint
		Names   []string
		BufSize   uint
	}
)



func CompactEncodeUint32(value uint32) []byte {
	var out []byte
	for copy := value; copy > 0; copy >>= 7 {
		b := byte(copy & 0x7f)
		if copy > 0x7f {
			b |= 0x80
		}
		out = append(out, b)
	}
	return out
}