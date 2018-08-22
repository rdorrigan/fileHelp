package fileDateSort

import (
	"os"
	"strings"
	"time"
)

// File is a file
type File struct {
	Info os.FileInfo
}

// Fldr is a folder
type Fldr struct {
	Files []File
}

// Less Checks the i is Before j
func (f Fldr) Less(i, j int) bool {
	// x := f.Files[i].Info.ModTime
	// y := f.Files[j].Info.ModTime
	// z := first.Before(y)
	return f.Files[i].Info.ModTime().Before(f.Files[j].Info.ModTime())
}

// Swap swaps files in a Fldr slice
func (f *Fldr) Swap(i, j int) {
	f.Files[i], f.Files[j] = f.Files[j], f.Files[i]
}

// Latest returns the most recently modified file in a folder
func (f Fldr) Latest() os.FileInfo {
	var modTime time.Time
	var names []File
	for _, i := range f.Files {
		if !i.Info.IsDir() {
			if i.Info.Mode().IsRegular() {
				if !strings.HasSuffix(i.Info.Name(), ".ini") {
					if !i.Info.ModTime().Before(modTime) {
						if i.Info.ModTime().After(modTime) {
							modTime = i.Info.ModTime()
							names = names[:0]
						}
						names = append(names, i)
					}
				}
			}
		}
	}
	return names[0].Info
}
