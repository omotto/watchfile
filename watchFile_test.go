package watchfile

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func setup(filename string, body string) (err error) {
	var f *os.File
	if f, err = os.Create(filename); err == nil {
		_, err = f.WriteString(body)
		e := f.Close()
		if err == nil {
			err = e
		}
	}
	return err
}

func modifyFile(filename string) (err error) {
	var f *os.File
	if f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend); err == nil {
		_, err = f.WriteString("\r\nmodify current test file")
		e := f.Close()
		if err == nil {
			err = e
		}
	}
	return err
}

func teardown(filename string) {
	_ = os.Remove(filename)
	return
}

func TestWatchFileError(t *testing.T) {
	var fileName = "test.txt"
	if err := setup(fileName, "test file to check watchfile package"); err != nil {
		t.Error(err)
	} else {
		if _, err := NewFileWatcher("fileName", 1, func(name string) { fmt.Println(name) }, fileName); err == nil {
			t.Error("NewFileWatcher should return Error, file not found")
		}
		if _, err := NewFileWatcher(fileName, 0, func(name string) { fmt.Println(name) }, fileName); err == nil {
			t.Error("NewFileWatcher should return Error, wrong period. It must be greater than 1s")
		}
		if _, err := NewFileWatcher(fileName, 1, func() { fmt.Println("Hello, world") }, 10); err == nil {
			t.Error("NewFileWatcher should return Error, wrong number of args")
		}
		if _, err := NewFileWatcher(fileName, 1, 0); err == nil {
			t.Error("NewFileWatcher should return Error, invalid function")
		}
		if _, err := NewFileWatcher(fileName, 1, func(s string, n int) { fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n) }, "s", "s2"); err == nil {
			t.Error("NewFileWatcher should return Error, invalid args types")
		}
	}
	teardown(fileName)
}

func TestWatchFileFunction(t *testing.T) {
	var fileName = "test.txt"
	if err := setup(fileName, "test file to check watchfile package"); err != nil {
		t.Error(err)
	} else {
		val := 123
		if wf, err := NewFileWatcher(fileName, 1, func(v *int) { *v++ }, &val); err == nil {
			wf.Start()
			for c := 0; c < 5; c++ {
				if err := modifyFile(fileName); err != nil {
					t.Error(err)
				}
				time.Sleep(time.Second * 2)
			}
			wf.Stop()
			t.Log(val)
			if val != 123+5 {
				t.Error("changes not cached")
			}
		} else {
			t.Error(err)
		}
	}
	teardown(fileName)
}
