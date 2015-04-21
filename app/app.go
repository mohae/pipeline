package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mohae/contour"
)

const Name = "pipeline"

func SetCfg() {
	contour.RegisterIntFlag("parallel", "p", 10, "10", "sets the parallelism of the file processing")
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	// For each regular file, start a goroutine that sums the file and sends
	// the result on c. Sent the result of the walk on errc.
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		// close the paths channel after walk returns
		defer close(paths)
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			// Abort the walk if done is closed.
			select {
			case paths <- path:
			case <-done:
				return fmt.Errorf("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}
