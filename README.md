# go-wg-wrapper

[![Go Report Card](https://goreportcard.com/badge/github.com/aschmidt75/go-wg-wrapper)](https://goreportcard.com/report/github.com/aschmidt75/go-wg-wrapper)
![Go](https://github.com/aschmidt75/go-wg-wrapper/workflows/Go/badge.svg)

This is a convenience wrapper around

* `/sbin/ip`
* the go [wgctrl](https://github.com/WireGuard/wgctrl-go) lib

to easily 

* create and delete wireguard interfaces,
* configure them with keys, endpoints etc.,
* manage peers

It does not have a CLI but is intended as a library only.
See [examples/main.go](examples/main.go) for details.

# Build 

This builds on Linux only because it is intended primarily for linux only.

# Test

e.g. use multipass to launch an ubuntu lts named wgtest:

```bash
$ multipass launch -c 1 -m 512M -n wgtest --cloud-init scripts/multipass-cloudinit.yaml lts
$ multipass mount $(pwd) wgtest:/mnt
$ multipass shell wgtest
ubuntu@wgtest:~$ sudo -i
root@wgtest:~# cd /mnt
```

Cloud-init script will install wireguard tools and go. Please run go / the binary as root since
it calls `/sbin/ip` and accesses wireguard via netlink.

```bash
$ go test -cover ./pkg/...
ok  	github.com/aschmidt75/go-wg-wrapper/pkg/wgwrapper	(cached)	coverage: 72.0% of statements
```

# Contribute

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

# License

Copyright 2020,2021 @aschmidt75 Licensed under the Apache License, Version 2.0
