package arguments

import "testing"

func Test_oneOf(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with one element, returns it", args{items: []string{"one"}}, "one"},
		{"with all zero-values, returns zero-values", args{items: []string{"", "", ""}}, ""},
		{"returns first", args{items: []string{"first", ""}}, "first"},
		{"returns last", args{items: []string{"", "last"}}, "last"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oneOf(tt.args.items...); got != tt.want {
				t.Errorf("oneOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
