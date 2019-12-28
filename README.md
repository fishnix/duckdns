# duck

Duck is a tiny Cli to update a domain in DuckDNS.

## usage

```bash
A DuckDNS client

Usage:
  duck [flags]
  duck [command]

Available Commands:
  help        Help about any command
  update      Update an IP at duckDNS

Flags:
      --config string        config file (default is $HOME/.duck.yaml)
  -C, --continuous           Run continuously on a timer
  -d, --domain strings       domain to update (may be repeated)
  -f, --frequency duration   time between updates as golang duration (default 5m0s)
  -h, --help                 help for duck
  -i, --ipv4 string          ipv4 ip address
  -p, --ipv6 string          ipv6 ip address
  -t, --token string         authentication token
  -v, --verbose              Be more verbose

Use "duck [command] --help" for more information about a command.
```
