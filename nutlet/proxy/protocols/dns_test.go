package protocols

import "testing"

func Test_match(t *testing.T) {
	type args struct {
		req string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"hit", args{"user.local.nutshell"}, "127.0.0.1", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := match(tt.args.req)
			if got != tt.want {
				t.Errorf("match() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
