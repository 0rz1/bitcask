package bitcask

import "testing"

func TestGetPath(t *testing.T) {
	cxt := &context{
		path: "abc",
	}
	if uGetPath(FT_Data, 1, cxt) != "abc/dat0001" {
		t.Error()
	}
	if uGetPath(FT_Data, 10001, cxt) != "abc/dat10001" {
		t.Error()
	}
	if uGetPath(FT_Location, 1, cxt) != "abc/loc0001" {
		t.Error()
	}
	if uGetPath(FT_Location, 10001, cxt) != "abc/loc10001" {
		t.Error()
	}
}
