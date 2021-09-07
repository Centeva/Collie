package testutils

type MockFileReader struct {
	Called      map[string]int
	CalledWith  map[string][]interface{}
	ReadFileRes []byte
}

func NewMockFileReader(ReadFileRes string) *MockFileReader {
	return &MockFileReader{
		Called:      make(map[string]int),
		CalledWith:  make(map[string][]interface{}),
		ReadFileRes: []byte(ReadFileRes),
	}
}

type FRReadFileArgs struct {
	Filename string
}

func (m *MockFileReader) ReadFile(filename string) ([]byte, error) {
	m.Called["readfile"]++
	m.CalledWith["readfile"] = append(m.CalledWith["readfile"], &FRReadFileArgs{filename})

	return m.ReadFileRes, nil
}
