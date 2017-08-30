package main

import (
	"reflect"
	"testing"
)

func Test_newHelper(t *testing.T) {
	tests := []struct {
		name string
		want *helper
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHelper(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHelper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_helper_Add(t *testing.T) {
	type fields struct {
		topics map[string]*fctCmd
	}
	type args struct {
		s string
		c *fctCmd
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &helper{
				topics: tt.fields.topics,
			}
			h.Add(tt.args.s, tt.args.c)
		})
	}
}
