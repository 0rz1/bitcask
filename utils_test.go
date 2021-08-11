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

func TextGetFTAndNo(t *testing.T) {
	cases := []struct {
		name string
		ft   FileType
		no   int
	}{
		{"qwer", FT_Invalid, 0},
		{"val", FT_Invalid, 0},
		{"dat1", FT_Invalid, 0},
		{"dat0001", FT_Data, 1},
		{"loc0111", FT_Location, 111},
		{"loc10111", FT_Location, 10111},
	}
	for _, cas := range cases {
		ft, no := uGetFTAndNo(cas.name)
		if cas.ft == FT_Invalid {
			if ft != FT_Invalid {
				t.Error()
			}
		} else {
			if cas.ft != ft || cas.no != no {
				t.Errorf("%v %d", ft, no)
			}
		}
	}
}
