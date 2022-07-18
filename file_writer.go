package spoor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type FileWriter struct {
	*bufio.Writer
	file          *os.File
	level         Level
	num           uint64 // The number of bytes written to this file
	maxSize       uint64
	logDir        string
	bufferSize    int
	flushInterval int //second
	mu            sync.Mutex
}

func (fw *FileWriter) Sync() error {
	return fw.file.Sync()
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	if fw.num+uint64(len(p)) >= fw.maxSize {
		if err := fw.rotateFile(time.Now()); err != nil {
			//fw.logger.exit(err)
		}
	}
	n, err = fw.Writer.Write(p)
	fw.num += uint64(n)
	if err != nil {
		//fw.logger.exit(err)
	}
	return
}

// rotateFile closes the syncBuffer's file and starts a new one.
func (fw *FileWriter) rotateFile(now time.Time) error {
	if fw.file != nil {
		fw.Flush()
		fw.file.Close()
	}
	var err error
	fw.file, _, err = createLogFile(fw.level.ToString(), fw.logDir, now)
	fw.num = 0
	if err != nil {
		return err
	}

	fw.Writer = bufio.NewWriterSize(fw.file, fw.bufferSize)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(&buf, "Running on machine: %s\n", host)
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&buf, "Log line format: [IWEF]mmdd hh:mm:ss.uuuuuu threadid file:line] msg\n")
	n, err := fw.file.Write(buf.Bytes())
	fw.num += uint64(n)
	return err
}

func createLogFile(levelName, logDir string, t time.Time) (f *os.File, filename string, err error) {
	if len(logDir) == 0 {
		return nil, "", errors.New("log: no log dirs")
	}
	name, link := getLogName(levelName, t)
	var lastErr error
	fname := filepath.Join(logDir, name)
	f, err = os.Create(fname)
	if err == nil {
		symlink := filepath.Join(logDir, link)
		os.Remove(symlink)        // ignore err
		os.Symlink(name, symlink) // ignore err
		return f, fname, nil
	}
	lastErr = err
	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}

func getLogName(levelName string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.log.%s.%04d-%02d-%02d-%02d-%02d-%02d.%4d",
		program,
		levelName,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
	return name, program + "." + levelName
}

// flushTicker periodically flushes the log file buffers.
func (fw *FileWriter) flushTicker() {
	for _ = range time.NewTicker(time.Second * time.Duration(fw.flushInterval)).C {
		fw.lockAndFlush()
	}
}

// lockAndFlush is like flush but locks l.mu first.
func (fw *FileWriter) lockAndFlush() {
	fw.mu.Lock()
	fw.flush()
	fw.mu.Unlock()
}

func (fw *FileWriter) flush() {
	file := fw.file
	if file != nil {
		fw.Flush()
		fw.Sync()
	}
}
