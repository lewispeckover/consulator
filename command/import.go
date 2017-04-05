package command

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/lewispeckover/consulator/command/configparser"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type ImportCommand struct {
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
	Purge       bool
}

func (c *ImportCommand) init() {
	if c.initialised {
		return
	}
	if c.Purge {
		c.name = "consulator sync"
		c.synopsis = "Syncs data into consul (like import, but with deletes)"
	} else {
		c.name = "consulator import"
		c.synopsis = "Imports data into consul"
	}
	c.args = "[options] [path ...]"
	c.flags = flag.NewFlagSet("import", flag.ContinueOnError)
	c.parseAsYAML = c.flags.Bool("yaml", false, "Parse stdin as YAML")
	c.parseAsJSON = c.flags.Bool("json", false, "Parse stdin as JSON")
	c.parseAsTAR = c.flags.Bool("tar", false, "Parse stdin as a tarball")
	c.arrayGlue = c.flags.String("glue", "\n", "Glue to use for joining array values")
	c.keyPrefix = c.flags.String("prefix", "", "Consul tree to work under")
	c.flags.Usage = func() { c.Ui.Output(c.Help()) }
	c.initialised = true
}

func (c *ImportCommand) Run(args []string) int {
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
	if err := c.syncConsul(data); err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err))
		return 1
	}
	return 0
}

func (c *ImportCommand) syncConsul(data map[string][]byte) error {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	kv := client.KV()
	deleted := 0
	updated := 0
	if c.Purge {
		pairs, _, err := kv.List(*c.keyPrefix, &api.QueryOptions{})
		if err != nil {
			return err
		}
		for _, pair := range pairs {
			// if there was a prefix, we need to strip it
			relativeKey := strings.TrimPrefix(pair.Key, *c.keyPrefix)
			if val, ok := data[relativeKey]; ok {
				if bytes.Equal(val, pair.Value) {
					delete(data, relativeKey)
				}
			} else if c.Purge {
				_, err := kv.Delete(pair.Key, nil)
				if err != nil {
					return err
				} else {
					deleted++
				}
			}
		}
	}
	for key, val := range data {
		_, err := kv.Put(c.toKVPair(key, val), nil)
		if err != nil {
			return err
		} else {
			updated++
		}
	}
	if c.Purge {
		c.Ui.Output(fmt.Sprintf("Sync completed. %d keys deleted, %d keys updated.", deleted, updated))
	} else {
		c.Ui.Output(fmt.Sprintf("Import completed. %d keys set.", updated))
	}
	return nil
}

func (c *ImportCommand) toKVPair(key string, val []byte) *api.KVPair {
	return &api.KVPair{
		Key:   *c.keyPrefix + key,
		Flags: 0,
		Value: val,
	}
}

func (c *ImportCommand) Synopsis() string {
	c.init()
	return c.synopsis
}

func (c *ImportCommand) Help() string {
	c.init()
	flagOut := new(bytes.Buffer)
	c.flags.SetOutput(flagOut)
	c.flags.PrintDefaults()
	c.flags.SetOutput(nil)
	return fmt.Sprintf("%s %s\n\n%s\n\nOptions:\n%s", c.name, c.args, c.synopsis, flagOut.String())
}
