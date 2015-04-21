package app

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/mohae/contour"
)

func MD5(path string) (string, error) {
	// Calculate the MD5 sum of all files under the specified directory.
	// If the provided path is just a file, and not a directory, then the
	// MD5 sum of that will be calculated
	m, err := MD5All(path)
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

// serial implementation
func SerialMD5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

// bounded, concurrent implementation
// stages: walkFiles -> digester -> results

// MD5Result holds the MD5 info for a file.
type MD5Result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

const parallelism = 20

func MD5Digester(done <-chan struct{}, paths <-chan string, c chan<- MD5Result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case c <- MD5Result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

// bounded concurrent implementation of md5all
func MD5All(root string) (map[string][md5.Size]byte, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc
	parallelism := contour.GetInt("parallel")

	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start the goroutines for digestion
	c := make(chan MD5Result)
	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			MD5Digester(done, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	// range through the results and add to the map
	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}
