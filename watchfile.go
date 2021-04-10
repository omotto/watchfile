package watchfile

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"time"
)

// WatchFile struct
type WatchFile struct {
	fileName	string				// fileName to watch
	fileSize 	int64				// File Size to check
	fileDate 	time.Time			// File modification date to check
	fileHash    []byte				// File SHA256 HASH to check
	// --
	period 		time.Duration		// Analysis watch period in seconds
	function	interface{}			// function to execute when file change happens
	fparams		[]interface{}		// function params
	// --
	stop		chan bool			// Stop channel
	running		bool				// Watcher status running (true) or stopped
}

// NewFileWatcher Creates a new FileWatcher instance
func NewFileWatcher(fileName string, period int, function interface{}, fparams ...interface{}) (wf WatchFile, err error) {
	var (
		f 	*os.File
		fi 	os.FileInfo
	)
	if f, err = os.Open(fileName); err == nil {
		h := sha256.New()
		if _, err = io.Copy(h, f); err == nil {
			if fi, err = f.Stat(); err == nil {
				if period > 0 {
					if !(function == nil || reflect.ValueOf(function).Kind() != reflect.Func) {
						if len(fparams) == reflect.TypeOf(function).NumIn() {
							for i := 0; i < reflect.TypeOf(function).NumIn(); i++ {
								functionParam := reflect.TypeOf(function).In(i)
								inputParam := reflect.TypeOf(fparams[i])
								if functionParam != inputParam {
									if functionParam.Kind() != reflect.Interface { return wf, fmt.Errorf(fmt.Sprintf("param[%d] must be be `%s` not `%s`", i, functionParam, inputParam)) }
									if !inputParam.Implements(functionParam) { return wf, fmt.Errorf(fmt.Sprintf("param[%d] of type `%s` doesn't implement interface `%s`", i, functionParam, inputParam)) }
								}
							}
							wf.fileName = fileName
							wf.fileSize = fi.Size()
							wf.fileDate = fi.ModTime()
							wf.fileHash = h.Sum(nil)
							wf.period 	= time.Duration(period * 1000000000) // Convert seconds to time.Duration struct
							wf.function = function
							wf.fparams 	= fparams
							wf.stop		= make(chan bool)
							wf.running	= false
						} else {
							err = errors.New("number of function params and number of provided params doesn't match")
						}
					} else {
						err = errors.New("invalid function parameter")
					}
				} else {
					err = errors.New("minimum period time 1s")
				}
			}
		}
		e := f.Close()
		if err == nil {
			err = e
		}
	}
	return wf, err
}

// Start starts file watcher
func (w *WatchFile) Start() {
	if w.running == false {
		w.running = true
		go func() {
			ticker := time.NewTicker(w.period)
			for {
				select {
				case <-ticker.C:
					var (
						f 	*os.File
						fi 	os.FileInfo
						err error
					)
					if f, err = os.Open(w.fileName); err == nil {
						if fi, err = f.Stat(); err == nil {
							if w.fileDate != fi.ModTime() {
								h := sha256.New()
								if _, err = io.Copy(h, f); err == nil {
									if bytes.Compare(w.fileHash, h.Sum(nil)) != 0 || w.fileSize != fi.Size() {
										w.fileDate = fi.ModTime()
										w.fileSize = fi.Size()
										w.fileHash = h.Sum(nil)
										go w.execFunction()
									}
								}
							}
						}
						_ = f.Close()
					}
					if err != nil {
						log.Println(err)
					}
				case <-w.stop:
					return
				}
			}
		}()
	}
}

//Stop stops file watcher
func (w *WatchFile) Stop() {
	if w.running == true {
		w.stop <- true
		w.running = false
	}
}

// Private Methods

func (w *WatchFile) execFunction() {
	defer func() {
		if r := recover(); r != nil { log.Println(r) }
	}()
	args := make([]reflect.Value, len(w.fparams))
	for i, param := range w.fparams {
		args[i] = reflect.ValueOf(param)
	}
	_ = reflect.ValueOf(w.function).Call(args)
}