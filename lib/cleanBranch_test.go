package lib_test

import (
	"bitbucket.org/centeva/collie/lib"
	"testing"
)

func Test_cleanBranch(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "should not change", args: args{name: "test-branch"}, want: "test-branch"},
		{name: "should replace slashes", args: args{name: "test/branch"}, want: "test-branch"},
		{name: "should strip specials", args: args{name: "test@$#!?&branch"}, want: "testbranch"},
		{name: "should lowercase", args: args{name: "TEST-BRANCH"}, want: "test-branch"},
		{name: "should handle complicated", args: args{name: "test@$#\\/!?&BRANcH/123-2/"}, want: "test-branch-123-2-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.CleanBranch(tt.args.name); got != tt.want {
				t.Errorf("cleanBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}
