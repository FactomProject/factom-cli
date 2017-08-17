package main

import "testing"

func TestGetFactomdServer(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{

		{name: "test", want: "localhost:8088"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFactomdServer(); got != tt.want {
				t.Errorf("GetFactomdServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extidsASCII_Set(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		e       *extidsASCII
		args    args
		wantErr bool
	}{
		{name: "test1", e: "HI", {s: "123"}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("extidsASCII.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extidsHex_Set(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		e       *extidsHex
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("extidsHex.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_namesASCII_Set(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		n       *namesASCII
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("namesASCII.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_namesHex_Set(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		n       *namesHex
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("namesHex.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extidsASCII_String(t *testing.T) {
	tests := []struct {
		name string
		e    *extidsASCII
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("extidsASCII.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extidsHex_String(t *testing.T) {
	tests := []struct {
		name string
		e    *extidsHex
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("extidsHex.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_namesASCII_String(t *testing.T) {
	tests := []struct {
		name string
		n    *namesASCII
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("namesASCII.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_namesHex_String(t *testing.T) {
	tests := []struct {
		name string
		n    *namesHex
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("namesHex.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_factoshiToFactoid(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := factoshiToFactoid(tt.args.v); got != tt.want {
				t.Errorf("factoshiToFactoid() = %v, want %v", got, tt.want)
			}
		})
	}
}
