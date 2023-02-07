package main

import (
	"bytes"
	"fmt"
	"io"
	"mkuznets.com/go/ytils/y"
	"mkuznets.com/go/ytils/ycrypto"
	"mkuznets.com/go/ytils/yerr"
	"os"
	"strings"
)

func main() {
	var buf bytes.Buffer
	yerr.Must(io.Copy(&buf, os.Stdin))
	cleanedInput := strings.TrimSpace(buf.String())
	obscured := y.Must(ycrypto.Obscure(cleanedInput))

	fmt.Println(obscured)
}
