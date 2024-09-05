# Scopious

Manage scope with ease

## Install

Scopious can be installed a couple different ways.

### Download from prebuilt binaries from GitHub releases.

Current build can be found at https://github.com/analog-substance/scopious/releases

### Golang Install

```bash
go install github.com/analog-substance/scopious@latest
```

## Usage

Scopious has a few ways to manage scope. It stores scoping information in text files within the `data/` folder. Scopious accepts scope as arguments and as data streams.

### Add

You can add IP addresses and domains. Scopious will figure out where to put them.

```bash
scopious add test.dev
scopious add 127.0.0.1/28
```

You can view your domains and IPs by running:
```bash
scopious domains
scopious ips -x
```

![Scopious add](docs/images/scopius-add.gif)

### Exclude

```bash
scopious add 127.0.0.1/28
scopious ips -x
scopious exclude 127.0.0.1/30
scopious ips -x
```

![Scopious exclude](docs/images/scopius-exclude.gif)

### Expand

Sometimes you don't want to add CIDRs to scope, but you need to expand them.

```bash
echo 127.0.0.1/28 | scopious expand
```
![Scopious expand](docs/images/scopius-expand.gif)

## About

Scope is stored in text files withing the `data/` dir by default. However, this behavior can be changed with the `--scope-dir` option or within the config file.

```txt
data/
    default/
        domains.txt
        ipv4.txt
        ipv6.txt
    internal/
        domains.txt
        ipv4.txt
        ipv6.txt
    internal-aws/
        domains.txt
        ipv4.txt
        ipv6.txt
```