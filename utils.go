package fileDateSort

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

// File is a file
type File struct {
	Info    os.FileInfo
	File    *os.File
	Headers []string
}

// Fldr is a folder
type Fldr struct {
	Files []File
}

// WritetoCSV writes the Orders to CSV
func WritetoCSV(dst string, data interface{}) error {
	o, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE| /*os.O_WRONLY|*/ os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	w := csv.NewWriter(o)
	// var toWrite []string
	// toWrite = append(toWrite, f.Headers[0:]...)
	// w.Write(toWrite)
	// Default
	w.UseCRLF = false
	refdata := reflect.TypeOf(data)
	nf := refdata.NumField()
	var wr []string
	for i := 0; i < nf; i++ {
		wr = append(wr, strings.TrimSpace(refdata.Field(i).Name))
		if len(wr) > 0 {
			w.Write(wr)
		}
		wr = append(wr[:0], wr[:0]...)
	}
	ref := reflect.ValueOf(data)
	for x := 0; x < nf; x++ {
		wr = append(wr, strings.TrimSpace(ref.Field(x).String()))
		if len(wr) > 0 {
			w.Write(wr)
		}
		wr = append(wr[:0], wr[:0]...)
	}
	w.Flush()
	return o.Close()
}

// FileExists checks that a file usually the log exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// LogWriter writes to the given log
func LogWriter(file string, content string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	log.Println(content)
	return f.Close()
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

// Cleaning cleans paths
func Cleaning(c string) string {
	cleaned := path.Clean(c)
	return cleaned
}

// CopyFileContents copies files
func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// TSVReadWriter reads and writes tsv files
func TSVReadWriter(r io.Reader, dst string) error {
	var out *os.File
	var err error
	//Implemented for removing a file that is downloaded daily
	// and does not need to be archived.
	// Comment out next block if archiving is necessary
	if _, err = os.Stat(dst); err == nil {
		err = os.Remove(dst)
		if err != nil {
			fmt.Println(err)
		}
	}
	out, err = os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	tsv := csv.NewReader(r)
	tsv.Comma = '\t'
	// LazyQuotes is a solution to the parse error
	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	// parse error on line 1414, column 78: bare " in non-quoted-field
	tsv.LazyQuotes = true
	tsvOut := csv.NewWriter(out)
	tsvOut.Comma = '\t'
	read, err := tsv.ReadAll()
	if err != nil {
		log.Fatalf("error'd at ReadAll: %v\n", err)
	}
	// WriteAll Calls Flush()
	err = tsvOut.WriteAll(read)
	if err != nil {
		log.Fatalf("Error'd at WriteAll: %v\n", err)
	}
	// tsvOut.Flush()
	//Would have worked if not tsv
	// _, err = io.Copy(out, f)
	// if err != nil {
	// 	return err
	// }
	s := out.Sync()
	if s != nil {
		log.Fatalf("Didn't sync %v\n", s)
	}
	return s
}
