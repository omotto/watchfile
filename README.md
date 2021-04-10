[![GoDoc](http://godoc.org/github.com/omotto/watchfile?status.png)](http://godoc.org/github.com/omotto/watchfile)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/omotto/watchfile)](https://pkg.go.dev/github.com/omotto/watchfile)
[![Build Status](https://travis-ci.com/omotto/watchfile.svg?branch=master)](https://travis-ci.com/omotto/watchfile)
[![Coverage Status](https://coveralls.io/repos/github/omotto/watchfile/badge.svg)](https://coveralls.io/github/omotto/watchfile)
[![Go Report Card](https://goreportcard.com/badge/github.com/omotto/watchfile)](https://goreportcard.com/report/github.com/omotto/watchfile)

# FileWatcher

Package watchfile implements a file watcher; if file is modified executes a custom function

### Installation

To download the specific tagged release, run:

```
go get github.com/omotto/watchfile
```

Import it in your program as:

```
import "github.com/omotto/watchfile"
```

### Usage

```
    num := 0
    if wf, err := NewFileWatcher(fileName, 1, func(v *int) { *v++ }, &num); err == nil {
        wf.Start()
        ....
        wf.Stop()
    }
```