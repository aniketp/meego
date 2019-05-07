package codegen

import (
	"bytes"
	"fmt"
)

var TMP_COUNT int

func write(b *bytes.Buffer, code string, args ...interface{}) {
	b.WriteString(fmt.Sprintf(code, args...))
}

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func freshTemp() string {
	TMP_COUNT += 1
	return fmt.Sprintf("tmp_%d", TMP_COUNT)
}
