package hb

import "testing"

func Test_asJStoJSON(t *testing.T) {
	type args struct {
		names      []string
		values     []string
		withQuotes bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{

			name: "happy case",
			args: args{
				names:      []string{"a", "b", "c"},
				values:     []string{"a", "b", "c"},
				withQuotes: true,
			},
			want: "JSON.stringify({ a: 'a', b: 'b', c: 'c', })",
		},
		{

			name: "happy case without quotes",
			args: args{
				names:      []string{"a", "b", "c"},
				values:     []string{"a", "b", "c"},
				withQuotes: false,
			},
			want: "JSON.stringify({ a: a, b: b, c: c, })",
		},
		{

			name: "less names",
			args: args{
				names:      []string{"a", "b"},
				values:     []string{"a", "b", "c"},
				withQuotes: true,
			},
			want: "",
		},
		{

			name: "less values",
			args: args{
				names:      []string{"a", "b", "c"},
				values:     []string{"a", "b"},
				withQuotes: true,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		if got := asJStoJSON(tt.args.names, tt.args.values, tt.args.withQuotes); got != tt.want {
			t.Errorf("%q. asJStoJSON() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
