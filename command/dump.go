package command

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/lewispeckover/consulator/command/configparser"

	"github.com/mitchellh/cli"
)

type DumpCommand struct {
	Ui          cli.Ui
	name        string
	args        string
	synopsis    string
	flags       *flag.FlagSet
	parseAsYAML *bool
	parseAsJSON *bool
	parseAsTAR  *bool
	arrayGlue   *string
	keyPrefix   *string
	initialised bool
}

func (c *DumpCommand) init() {
	if c.initialised {
		return
	}
	c.name = "consulator dump"
	c.args = "[options] [path ...]"
	c.synopsis = "Dumps parsed config as JSON suitable for use with consul kv import"
	c.flags = flag.NewFlagSet("dump", flag.ContinueOnError)
	c.parseAsYAML = c.flags.Bool("yaml", false, "Parse stdin as YAML")
	c.parseAsJSON = c.flags.Bool("json", false, "Parse stdin as JSON")
	c.parseAsTAR = c.flags.Bool("tar", false, "Parse stdin as a tarball")
	c.arrayGlue = c.flags.String("glue", "\n", "Glue to use for joining array values")
	c.keyPrefix = c.flags.String("prefix", "", "Key prefix to use for output")
	c.flags.Usage = func() { c.Ui.Output(c.Help()) }
	c.initialised = true
}

func (c *DumpCommand) Run(args []string) int {
	c.init()
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	if *c.parseAsYAML && *c.parseAsJSON {
		c.Ui.Error("Only one input format may be specified")
		return 1
	}
	// clean up the prefix
	*c.keyPrefix = strings.TrimSuffix(strings.TrimSpace(*c.keyPrefix), "/")
	if *c.keyPrefix != "" {
		*c.keyPrefix = *c.keyPrefix + "/"
	}
	data := make(map[string][]byte)
	if c.flags.NArg() == 0 {
		switch {
		case *c.parseAsYAML:
			if err := configparser.ParseAsYAML("/dev/stdin", data, *c.arrayGlue); err != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", err))
				return 1
			}
		case *c.parseAsJSON:
			if err := configparser.ParseAsJSON("/dev/stdin", data, *c.arrayGlue); err != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", err))
				return 1
			}
		case *c.parseAsTAR:
			if err := configparser.ParseAsTAR("/dev/stdin", data, *c.arrayGlue); err != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", err))
				return 1
			}
		default:
			c.Ui.Error("You must specify an input format when using stdin\n")
			c.Ui.Error(c.Help())
			return 1
		}
	} else {
		for _, p := range c.flags.Args() {
			if err := configparser.Parse(p, data, *c.arrayGlue); err != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", err))
				return 1
			}
		}
	}

	exported := make([]*dumpExportEntry, len(data))
	i := 0
	for key, val := range data {
		exported[i] = c.toExportEntry(key, val)
		i++
	}
	json, err := json.MarshalIndent(exported, "", "\t")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error exporting data: %s", err))
		return 1
	}
	c.Ui.Output(string(json))
	return 0
}

type dumpExportEntry struct {
	Key   string `json:"key"`
	Flags uint64 `json:"flags"`
	Value string `json:"value"`
}

func (c *DumpCommand) toExportEntry(key string, val []byte) *dumpExportEntry {
	return &dumpExportEntry{
		Key:   *c.keyPrefix + key,
		Flags: 0,
		Value: base64.StdEncoding.EncodeToString(val),
	}
}

func (c *DumpCommand) Synopsis() string {
	c.init()
	return c.synopsis
}

func (c *DumpCommand) Help() string {
	c.init()
	flagOut := new(bytes.Buffer)
	c.flags.SetOutput(flagOut)
	c.flags.PrintDefaults()
	c.flags.SetOutput(nil)
	return fmt.Sprintf("%s %s\n\n%s\n\nOptions:\n%s", c.name, c.args, c.synopsis, flagOut.String())
}
