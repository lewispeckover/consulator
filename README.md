
# Consulator

Consulator lets you import, diff, and synchronize your KV data from JSON or YAML sources directly to Consul. This allows you to easily version your configuration data, storing it in git or anywhere else you see fit.

[![CircleCI](https://circleci.com/gh/lewispeckover/consulator/tree/master.svg?style=shield)](https://circleci.com/gh/lewispeckover/consulator/tree/master)

## Getting Consulator

Docker is the easiest way. You can find it on the [Docker Hub](https://hub.docker.com/r/lewispeckover/consulator/).


```
docker run -it --rm lewispeckover/consulator -help
```

Alternatively, download a binary from the [releases](https://github.com/lewispeckover/consulator/releases/latest) page, or clone the repo and build it yourself.

## Running Consulator

```
Usage: ./consulator [OPTIONS] [PATH]

PATH should be the path to a file or directory that contains your data.
If no path is provided, stdin is used. In this case, -format must be specified.

Options:

  -debug
    	Show debugging information
  -dump
    	Dump loaded data as JSON, suitable for using in a 'consul kv import'
  -format string
    	Specify data format(json or yaml) when reading from stdin.
  -glue string
    	Glue to use when joining array values (default "\n")
  -prefix string
    	Specifies a Consul tree to work under.
  -quiet
    	Only show errors
  -sync
    	Sync to consul
  -trace
    	Show even more debugging information
  -version
    	Show version

The usual Consul client environment variables can be used to configure the connection:

 - CONSUL_HTTP_ADDR
 - CONSUL_HTTP_TOKEN
 - CONSUL_HTTP_SSL

Etc. See https://www.consul.io/docs/commands/ for a complete list.
```


## Source data

JSON and YAML sources are supported. Data can be loaded from files or piped from standard input. Note that Consul KV values are only allowed to be strings, so non-string values are converted where possible. Most significantly, array values are joined with a configurable glue string (default: "\n").

Given a `myapp.yaml`:

```
tags:
 - production
 - web
```
or equivalent `myapp.json`:

```
{ 
 "tags": ["production", "web"]
}
```

Running `consulator -sync -glue=, myapp.yaml` will result in a Consul key `myapp/tags` with the value `production,web`

When a directory is specified as the source, it is scanned for files with extensions .json, .yaml or .yml. Subdirectories and filenames are used to build key prefixes. 

Suppose that the above file is located at `/etc/consuldata/config/myapp.yaml`. When consulator is executed as `consulator -sync -glue=, /etc/consuldata`, it will result in a Consul key `config/myapp/tags` with value `production,web`.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

