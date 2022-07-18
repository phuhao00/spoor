package spoor

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	pid      = 0
	program  = ""
	host     = ""
	userName = ""
	pidStr   = ""
)

func init() {
	pid = os.Getpid()
	program = filepath.Base(os.Args[0])
	pidStr = fmt.Sprintf("pid:%05d", pid)
}
