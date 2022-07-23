package spoor

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"
)

func TestName(t *testing.T) {
	l := NewSpoor(DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, WithConsoleWriter(os.Stdout))
	l.DebugF("hhhh")
	l.InfoF("jjjjj")
}

func TestFileWriter(t *testing.T) {
	fileWriter := NewFileWriter("log", 0, 0, 0)
	l := NewSpoor(DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, WithFileWriter(fileWriter))
	l.DebugF("hhhh")
	select {}

}

func TestName1(t *testing.T) {
	_, file, line, ok := runtime.Caller(-1)
	fmt.Println(file, line, ok)

}
