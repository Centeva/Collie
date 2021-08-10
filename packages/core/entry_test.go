package core_test

import (
	"testing"

	"bitbucket.org/centeva/collie/packages/core"
)

type mockCommandParser struct {
	called map[string]int
}

func (c mockCommandParser) ParseCommands() (err error) {
	c.called["ParseCommands"]++
	return
}

func Test_main(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		want     int
	}{
		{"Should set logger", "ParseCommands", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &mockCommandParser{
				called: make(map[string]int),
			}

			// var buf bytes.Buffer
			// log.SetOutput(&buf)
			// defer func() {
			// 	log.SetOutput(os.Stderr)
			// }()
			// output := buf.String()

			core.Entry(cmd)
			if cmd.called[tt.funcName] != tt.want {
				t.Errorf("Entry(): commandParser.%s() Should have been called (%d) times got: %v", tt.funcName, tt.want, cmd.called)
			}
		})
	}
}
