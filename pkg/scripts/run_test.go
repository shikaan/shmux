package scripts

import "testing"

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
