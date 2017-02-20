package main

import (
	"encoding/base64"
	//"fmt"
	"bytes"
	"github.com/hashicorp/consul/api"
	//"strings"
)

var config = api.DefaultConfig()
var client, _ = api.NewClient(config)

func syncConsul() {
	Info.Println("Walking consul")
	kv := client.KV()
	pairs, _, err := kv.List("/", &api.QueryOptions{})
	if err != nil {
		Error.Fatal(err)
	}
	Info.Printf("Syncing %d keys to consul...\n", len(data))
	d := 0
	u := 0
	for _, pair := range pairs {
		if val, ok := data[pair.Key]; ok {
			if bytes.Equal(val, pair.Value) {
				Trace.Printf("key %s can be ignored\n", pair.Key)
				delete(data, pair.Key)
			}
		} else {
			Trace.Printf("key %s should be removed\n", pair.Key)
			_, err := kv.Delete(pair.Key, nil)
			if err != nil {
				Error.Printf("Error removing %s: %s\n", pair.Key, err)
			} else {
				d++
			}
		}
	}
	for key, val := range data {
		Trace.Printf("set %s to %s", key, string(val))
		_, err := kv.Put(toKVPair(key, val), nil)
		if err != nil {
			Error.Printf("Error putting %s: %s: %s\n", key, string(val), err)
		} else {
			u++
		}
	}
	Info.Printf("Sync completed. %d keys deleted, %d keys updated.\n", d, u)
}

type kvExportEntry struct {
	Key   string `json:"key"`
	Flags uint64 `json:"flags"`
	Value string `json:"value"`
}

func toExportEntry(key string, val []byte) *kvExportEntry {
	return &kvExportEntry{
		Key:   key,
		Flags: 0,
		Value: base64.StdEncoding.EncodeToString(val),
	}
}

func toKVPair(key string, val []byte) *api.KVPair {
	return &api.KVPair{
		Key:   key,
		Flags: 0,
		Value: val,
	}
}
