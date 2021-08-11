package bitcask

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

func uGetFTAndNo(name string) (ft FileType, no int) {
	if name[:3] == "dat" {
		i, err := strconv.ParseInt(name[3:], 10, 32)
		if err != nil {
			return
		} else if fmt.Sprintf("dat%04d", i) != name {
			return
		}
		ft = FT_Data
		no = int(i)
	} else if name[:3] == "loc" {
		i, err := strconv.ParseInt(name[3:], 10, 32)
		if err != nil {
			return
		} else if fmt.Sprintf("loc%04d", i) != name {
			return
		}
		ft = FT_Location
		no = int(i)
	}
	return
}

func uGetPath(ft FileType, no int, c *context) string {
	var name string
	if ft == FT_Data {
		name = fmt.Sprintf("dat%04d", no)
	} else {
		name = fmt.Sprintf("loc%04d", no)
	}
	return path.Join(c.path, name)
}

func uOpen(ft FileType, no int, c *context) (*os.File, error) {
	return os.Open(uGetPath(ft, no, c))
}

func uOpenAppend(ft FileType, no int, c *context) (*os.File, error) {
	path := uGetPath(ft, no, c)
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
