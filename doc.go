/*

Package watchfile implements a file watcher; if file is modified executes a custom function

Installation

To download the specific tagged release, run:

	go get github.com/omotto/watchfile

Import it in your program as:

	import "github.com/omotto/watchfile"

Usage

    num := 0
    if wf, err := NewFileWatcher(fileName, 1, func(v *int) { *v++ }, &num); err == nil {
        wf.Start()
        ....
        wf.Stop()
    }

*/
package watchfile
