package lib_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"bitbucket.org/centeva/collie/lib"
)

type ICommandParser interface {
	getBranch() *string
	parseFlags()
}

type testArgs struct {
	CleanBranch string
}

type mockCommandParser struct {
	args testArgs
}

func (c mockCommandParser) GetBranch() *string {
	val := c.args.CleanBranch
	return &val
}

func (c mockCommandParser) ParseFlags() {}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		args testArgs
		want string
	}{
		{"Should cleanBranch", testArgs{CleanBranch: "test/branch1!"}, "test-branch1\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &mockCommandParser{
				args: tt.args,
			}

			var buf bytes.Buffer

			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			lib.Entry(cmd)

			output := buf.String()

			if output != tt.want {
				t.Errorf("Entry() = %v, want %v", output, tt.want)
			}
		})
	}
}
