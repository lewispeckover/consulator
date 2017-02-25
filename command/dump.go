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
	Ui cli.Ui
}

var commandName = "consulator dump"
var commandArgs = "[options] [path ...]"
var synopsis = "Dumps parsed config as JSON suitable for use with consul kv import"
var dumpFlags = flag.NewFlagSet("dump", flag.ContinueOnError)
var parseAsYAML = dumpFlags.Bool("yaml", false, "Parse stdin as YAML")
var parseAsJSON = dumpFlags.Bool("json", false, "Parse stdin as JSON")
var arrayGlue = dumpFlags.String("glue", "\n", "Glue to use for joining array values")
var keyPrefix = dumpFlags.String("prefix", "", "Key prefix to use for output")

func (c *DumpCommand) Run(args []string) int {
	dumpFlags.Usage = func() { c.Ui.Output(c.Help()) }
	if err := dumpFlags.Parse(args); err != nil {
		return 1
	}
	if *parseAsYAML && *parseAsJSON {
		c.Ui.Error("Only one input format may be specified")
		return 1
	}
	// clean up the prefix
	*keyPrefix = strings.TrimSuffix(strings.TrimSpace(*keyPrefix), "/")
	if *keyPrefix != "" {
		*keyPrefix = *keyPrefix + "/"
	}
	data := make(map[string][]byte)
	if dumpFlags.NArg() == 0 {
		switch {
		case *parseAsYAML:
			if err := configparser.ParseAsYAML("/dev/stdin", data, *arrayGlue); err != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", err))
				return 1
			}
		case *parseAsJSON:
			if err := configparser.ParseAsJSON("/dev/stdin", data, *arrayGlue); err != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", err))
				return 1
			}
		default:
			c.Ui.Error("You must specify an input format when using stdin\n")
			c.Ui.Error(c.Help())
			return 1
		}
	} else {
		for _, p := range dumpFlags.Args() {
			if err := configparser.Parse(p, data, *arrayGlue); err != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", err))
				return 1
			}
		}
	}

	exported := make([]*dumpExportEntry, len(data))
	i := 0
	for key, val := range data {
		exported[i] = c.toExportEntry(key, val, *keyPrefix)
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

func (c *DumpCommand) toExportEntry(key string, val []byte, prefix string) *dumpExportEntry {
	return &dumpExportEntry{
		Key:   prefix + key,
		Flags: 0,
		Value: base64.StdEncoding.EncodeToString(val),
	}
}

func (c *DumpCommand) Synopsis() string {
	return synopsis
}

func (c *DumpCommand) Help() string {
	flagOut := new(bytes.Buffer)
	dumpFlags.SetOutput(flagOut)
	dumpFlags.PrintDefaults()
	dumpFlags.SetOutput(nil)
	return fmt.Sprintf("%s %s\n\n%s\n\nOptions:\n%s", commandName, commandArgs, synopsis, flagOut.String())
}
