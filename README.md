
# Consulator

Consulator lets you import and synchronize your KV data from JSON or YAML sources directly to Consul. This allows you to easily version your configuration data, storing it in git or anywhere else you see fit.

[![CircleCI](https://circleci.com/gh/lewispeckover/consulator/tree/master.svg?style=shield)](https://circleci.com/gh/lewispeckover/consulator/tree/master)

## Getting Consulator

Docker is the easiest way. You can find it on the [Docker Hub](https://hub.docker.com/r/lewispeckover/consulator/).


```
docker run -it --rm lewispeckover/consulator -help
```

Alternatively, download a binary from the [releases](https://github.com/lewispeckover/consulator/releases/latest) page, or clone the repo and build it yourself.

## Running Consulator

```
Usage: consulator [--version] [--help] <command> [<options>] [<path> ...]

Available commands are:
    dump       Dumps parsed config as JSON suitable for use with consul kv import
    import     Imports data into consul
    version    Prints the version

Options:
  -glue string
    	Glue to use for joining array values (default "\n")
  -json
    	Parse stdin as JSON
  -prefix string
    	Key prefix to use for output / Consul import destination
  -yaml
    	Parse stdin as YAML

Multiple paths (files or directories) may be provided, they are parsed in order. This allows you to specify some default values in the first path.
If no paths are provided, stdin is used. In this case, -yaml or -json must be specified.

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

Running `consulator import -glue=, myapp.yaml` will result in a Consul key `myapp/tags` with the value `production,web`

When a directory is specified as the source, it is scanned for files with extensions .json, .yaml or .yml. Subdirectories and filenames are used to build key prefixes. 

Suppose that the above file is located at `/etc/consuldata/config/myapp.yaml`. When consulator is executed as `consulator import -glue=, /etc/consuldata`, it will result in a Consul key `config/myapp/tags` with value `production,web`.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

