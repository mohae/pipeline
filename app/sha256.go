package app

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"sort"
	"sync"

	"github.com/mohae/contour"
)

func SHA256(path string) (string, error) {
	// Calculate the MD5 sum of all files under the specified directory.
	// If the provided path is just a file, and not a directory, then the
	// MD5 sum of that will be calculated
	m, err := SHA256All(path)
	if err != nil {
		return "", err
	}

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
	return "", nil
}

// SHA256Result holds the info for SHA256 hashes
type SHA256Result struct {
	path string
	sum  [sha256.Size]byte
	err  error
}

func SHA256Digester(done <-chan struct{}, paths <-chan string, c chan<- SHA256Result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case c <- SHA256Result{path, sha256.Sum256(data), err}:
		case <-done:
			return
		}
	}
}

// bounded concurrent implementation of SHA256All
func SHA256All(root string) (map[string][sha256.Size]byte, error) {
	// SHA256All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start the goroutines for digestion
	c := make(chan SHA256Result)
	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			SHA256Digester(done, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	// range through the results and add to the map
	sha := make(map[string][sha256.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		sha[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return sha, nil
}
