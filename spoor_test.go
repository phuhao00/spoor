package spoor

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"
)

func TestName(t *testing.T) {
	NewSpoor(DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, WithConsoleWriter(os.Stdout))
}

func TestFileWriter(t *testing.T) {
	fileWriter := NewFileWriter("log", 0, 0, 0)
	NewSpoor(DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, WithFileWriter(fileWriter))
	select {}

}

func TestName1(t *testing.T) {
	_, file, line, ok := runtime.Caller(-1)
	fmt.Println(file, line, ok)

}
