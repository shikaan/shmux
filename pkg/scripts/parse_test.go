package scripts

import (
	"reflect"
	"slices"
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
		wantDeps				 []string
	}{
		{"returns true with correct script", args{"start:"}, true, "start", nil},
		{"interprets further strings as dependencies", args{"start: dunno"}, true, "start", []string{"dunno"}},
		{"returns all the dependencies", args{"start: dep1 dep2 dep3"}, true, "start", []string{"dep1", "dep2", "dep3"}},
		{"dedupes dependencies", args{"start: dep1 dep1 dep1"}, true, "start", []string{"dep1"}},
		{"handles whitespaces", args{"start:dep1     dep2    dep3       "}, true, "start", []string{"dep1", "dep2", "dep3"}},
		{"returns false with blank script", args{"    start:"}, false, "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsScriptLine, gotMatch, gotDeps := readScript(tt.args.line)
			if gotIsScriptLine != tt.wantIsScriptLine {
				t.Errorf("readScript() gotIsScriptLine = %v, want %v", gotIsScriptLine, tt.wantIsScriptLine)
			}
			if gotMatch != tt.wantMatch {
				t.Errorf("readScript() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
			if !slices.Equal(gotDeps, tt.wantDeps) {
				t.Errorf("readScript() gotDeps = %v, want %v", gotDeps, tt.wantDeps)
			}
		})
	}
}

func TestReadScript(t *testing.T) {
	type args struct {
		scriptName string
		shell      string
		file       []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Script
		wantErr bool
	}{
		{
			"one script, with non-default interpreter",
			args{"script", "/bin/interpreter", []byte("script:\n\techo $HOME")},
			&Script{Name: "script", Lines: []string{"echo $HOME"}, Interpreter: "/bin/interpreter"},
			false,
		},
		{
			"one script, one line",
			args{"script", "/bin/bash", []byte("script:\n\techo $HOME")},
			&Script{Name: "script", Lines: []string{"echo $HOME"}, Interpreter: "/bin/bash"},
			false,
		},
		{
			"one script, many lines",
			args{"script", "/bin/bash", []byte("script:\n\tif true:\n\t\tprint('test2')")},
			&Script{Name: "script", Lines: []string{"if true:", "\tprint('test2')"}, Interpreter: "/bin/bash"},
			false,
		},
		{
			"one script, many lines with spaces",
			args{"script", "/bin/bash", []byte("script:\n   if true:\n     print('test2')")},
			&Script{Name: "script", Lines: []string{"if true:", "  print('test2')"}, Interpreter: "/bin/bash"},
			false,
		},
		{
			"two scripts, retuns latter",
			args{"another", "/bin/bash", []byte("script:\n\tscript1\n\nanother:\n\tscript2")},
			&Script{Name: "another", Lines: []string{"script2"}, Interpreter: "/bin/bash"},
			false,
		},
		{
			"two scripts, finds nothing",
			args{"unknown", "/bin/bash", []byte("script:\n\tscript1\n\nanother:\n\tscript2")},
			nil,
			true,
		},
		{
			"one script, with simple shebang",
			args{"script", "/bin/bash", []byte("script:\n\t#!/bin/lol\n\tline")},
			&Script{Name: "script", Lines: []string{"line"}, Interpreter: "/bin/lol"},
			false,
		},
		{
			"one script, with shebang with options",
			args{"script", "/bin/bash", []byte("script:\n\t#!/bin/lol --lol=1 --two\n\tline")},
			&Script{Name: "script", Lines: []string{"line"}, Interpreter: "/bin/lol", Options: []string{"--lol=1", "--two"}},
			false,
		},
		{
			"with one dependency",
			args{"top_level", "/bin/bash", []byte("top_level: dep1\n\tline\ndep1:\n\tline2")},
			&Script{Name: "top_level", Lines: []string{"line"}, Interpreter: "/bin/bash", Dependencies: []Script{{Name: "dep1", Lines: []string{"line2"}, Interpreter: "/bin/bash"}}},
			false,
		},
		{
			"with more dependencies",
			args{"top_level", "/bin/bash", []byte("top_level: dep1 dep2\n\tline0\ndep1:\n\tline1\ndep2:\n\tline2")},
			&Script{Name: "top_level", Lines: []string{"line0"}, Interpreter: "/bin/bash", Dependencies: []Script{
				{Name: "dep1", Lines: []string{"line1"}, Interpreter: "/bin/bash"},
				{Name: "dep2", Lines: []string{"line2"}, Interpreter: "/bin/bash"},
			}},
			false,
		},
		{
			"with nesting",
			args{"top_level", "/bin/bash", []byte("top_level: dep1\n\tline0\ndep1: dep2\n\tline1\ndep2:\n\tline2")},
			&Script{
				Name: "top_level",
				Lines: []string{"line0"},
				Interpreter: "/bin/bash",
				Dependencies: []Script{
					{
						Name: "dep1",
						Lines: []string{"line1"},
						Interpreter: "/bin/bash",
						Dependencies: []Script{
							{
								Name: "dep2",
								Lines: []string{"line2"},
								Interpreter: "/bin/bash",
							},
						},
					},
				},
			},
			false,
		},
		{
			"with circular dependency",
			args{"top_level", "/bin/bash", []byte("top_level: top_level\n\tline")},
			nil,
			true,
		},
		{
			"with too nested dependencies",
			args{"top_level", "/bin/bash", []byte("top_level: dep1\n\tline\n\ndep1: dep2\n\tline2\n\ndep2: dep3\n\tline3\n\ndep3: top_level\n\tline4")},
			nil,
			true,
		},
	};

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadScript(tt.args.scriptName, tt.args.shell, tt.args.file, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == tt.want {
				return
			}

			sameName := got.Name == tt.want.Name
			sameLines := slices.Equal(got.Lines, tt.want.Lines)
			sameInterpreter := got.Interpreter == tt.want.Interpreter
			sameOptions := slices.Equal(got.Options, tt.want.Options)
			sameDeps := slices.EqualFunc(got.Dependencies, tt.want.Dependencies, func(a Script, b Script) bool {
				return reflect.DeepEqual(a,b)
			})
			
			if !sameName || !sameLines || !sameInterpreter || !sameOptions || !sameDeps {
				t.Errorf("ReadScript() = %v, want %v", got, tt.want)
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