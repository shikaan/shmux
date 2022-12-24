package scripts

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_readScript(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name             string
		args             args
		wantIsScriptLine bool
		wantMatch        string
	}{
		{"returns true with correct script", args{"start:"}, true, "start"},
		{"handles and ignores foreign characters", args{"start: dunno"}, true, "start"},
		{"returns false with blank start", args{"    start:"}, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsScriptLine, gotMatch := readScript(tt.args.line)
			if gotIsScriptLine != tt.wantIsScriptLine {
				t.Errorf("readScript() gotIsScriptLine = %v, want %v", gotIsScriptLine, tt.wantIsScriptLine)
			}
			if gotMatch != tt.wantMatch {
				t.Errorf("readScript() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestReadScript(t *testing.T) {
	type args struct {
		scriptName string
		file       io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Script
		wantErr bool
	}{
		{"one script, one line", args{"script", strings.NewReader("script:\n\techo $HOME")}, []string{"echo $HOME"}, false},
		{"one script, many lines", args{"script", strings.NewReader("script:\n\tif true:\n\t\tprint('test2')")}, []string{"if true:", "\tprint('test2')"}, false},
		{"one script, many lines with spaces", args{"script", strings.NewReader("script:\n   if true:\n     print('test2')")}, []string{"if true:", "  print('test2')"}, false},
		{"two scripts, retuns latter", args{"another", strings.NewReader("script:\n\tscript1\n\nanother:\n\tscript2")}, []string{"script2"}, false},
		{"two scripts, finds nothing", args{"unknown", strings.NewReader("script:\n\tscript1\n\nanother:\n\tscript2")}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadScript(tt.args.scriptName, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadScript() = %v, want %v", got, tt.want)
			}
		})
	}
}
