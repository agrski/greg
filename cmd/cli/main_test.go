package main

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_getFiletypes(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := getFiletypes(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getFiletypes() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_getLocation(t *testing.T) {
	tests := []struct {
		name string
		want location
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := getLocation(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getLocation() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_getSearchPattern(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := getSearchPattern(); got != tt.want {
					t.Errorf("getSearchPattern() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_isEmpty(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := isEmpty(tt.args.s); got != tt.want {
					t.Errorf("isEmpty() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_makeURI(t *testing.T) {
	type args struct {
		l location
	}
	tests := []struct {
		name string
		args args
		want url.URL
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := makeURI(tt.args.l); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("makeURI() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_parseLocationFromUrl(t *testing.T) {
	type args struct {
		rawUrl string
	}
	tests := []struct {
		name string
		args args
		want location
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := parseLocationFromUrl(tt.args.rawUrl); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("parseLocationFromUrl() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
