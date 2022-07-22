package spoor

import (
	"log"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	l := NewSpoor(DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, WithNormalWriter(os.Stdout))
	l.DebugF("hhhh")
}

func TestFileWriter(t *testing.T) {
	fileWriter := NewFileWriter(".", 0, 0, 0)
	l := NewSpoor(DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, WithFileWriter(fileWriter))
	l.DebugF("hhhh")
	select {}

}
