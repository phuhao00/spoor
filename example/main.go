package main

import (
	"github.com/phuhao00/spoor"
	"log"
)

func main() {
	fileWriter := spoor.NewFileWriter("log", 0, 0, 0)
	l := spoor.NewSpoor(spoor.DEBUG, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile, spoor.WithFileWriter(fileWriter))
	l.DebugF("hhhh")
	select {}
}
