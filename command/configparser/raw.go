package configparser

import (
	"io"
	"io/ioutil"
	"strings"
)

func parseRaw(fp io.Reader, prefix []string, glue string) error {
	contents, err := ioutil.ReadAll(fp)
	if err == nil {
		data[strings.Join(prefix, "/")] = contents
	}
	return err
}
