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

func Test_replaceArguments(t *testing.T) {
	type args struct {
		content    string
		arguments  []string
		scriptName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"replace one occurrence", args{"this $1 replacement, $1 it not?", []string{"is"}, "script"}, "this is replacement, is it not?"},
		{"replace script name", args{"this $1 $@", []string{"is"}, "script"}, "this is script"},
		{"replace many occurrences", args{"this $2 $1 $@", []string{"not", "is"}, "script"}, "this is not script"},
		{"replace only 9 occurrences", args{"$9$8$7$6$5$4$3$2$1", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, "script"}, "987654321"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceArguments(tt.args.content, tt.args.arguments, tt.args.scriptName); got != tt.want {
				t.Errorf("replaceArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cap(t *testing.T) {
	type args struct {
		slice []string
		n     int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"shorter", args{[]string{"1"}, 2}, []string{"1"}},
		{"equal", args{[]string{"1", "2"}, 2}, []string{"1", "2"}},
		{"longer", args{[]string{"1", "2"}, 1}, []string{"1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cap(tt.args.slice, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cap() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
