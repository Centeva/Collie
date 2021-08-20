package external

import "io/ioutil"

type IFileReader interface {
	ReadFile(filename string) ([]byte, error)
}

type FileReader struct{}

func (f *FileReader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
