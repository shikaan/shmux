package scripts

import (
	"io"
	"strings"
	"testing"
)

func TestMakeHelp(t *testing.T) {
	type args struct {
		file io.Reader
	}
	tests := []struct {
		name     string
		args     args
		includes string
	}{
		{"returns the correct text", args{strings.NewReader("script:\n\tscript1\n\nanother:\n\tscript2")}, "script1, script2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeHelp(tt.args.file); strings.Contains(got, tt.includes) {
				t.Errorf("MakeHelp() = %v, does not include %v", got, tt.includes)
			}
		})
	}
}