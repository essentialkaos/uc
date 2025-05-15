<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/uc"><img src="https://kaos.sh/r/uc.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/w/uc/ci"><img src="https://kaos.sh/w/uc/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/uc/codeql"><img src="https://kaos.sh/w/uc/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#usage">Usage</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`uc` is a simple utility for counting unique lines.

### Usage demo

[![demo](https://github.com/user-attachments/assets/24bd2d94-cdcd-41b7-b3d4-f9ac7b0b239e)](#usage-demo)

### Benchmarks

```
$ wc -l data.txt
18408096 data.txt

$ hyperfine 'sort -u data.txt | wc -l' 'uc -np data.txt'
Benchmark 1: sort -u data.txt | wc -l
  Time (mean ± σ):     16.030 s ±  0.181 s    [User: 86.713 s, System: 1.165 s]
  Range (min … max):   15.699 s … 16.324 s    10 runs

Benchmark 2: uc -np data.txt
  Time (mean ± σ):      2.889 s ±  0.101 s    [User: 2.435 s, System: 0.454 s]
  Range (min … max):    2.721 s …  3.065 s    10 runs

Summary
  uc -np data.txt ran
    5.55 ± 0.20 times faster than sort -u data.txt | wc -l
```

### Installation

#### From sources

To build the `uc` from scratch, make sure you have a working Go [1.23+](https://github.com/essentialkaos/.github/blob/master/GO-VERSION-SUPPORT.md) workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/uc@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://pkgs.kaos.st)

```bash
sudo dnf install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo dnf install uc
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/uc/latest).

To install the latest prebuilt version, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) uc
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo uc --completion=bash 1> /etc/bash_completion.d/uc
```


ZSH:
```bash
sudo uc --completion=zsh 1> /usr/share/zsh/site-functions/uc
```


Fish:
```bash
sudo uc --completion=fish 1> /usr/share/fish/vendor_completions.d/uc.fish
```

### Man documentation

You can generate man page for `uc` using next command:

```bash
uc --generate-man | sudo gzip > /usr/share/man/man1/uc.1.gz
```

### Usage

<img src=".github/images/usage.svg" />

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/uc/ci.svg?branch=master)](https://kaos.sh/w/uc/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/uc/ci.svg?branch=master)](https://kaos.sh/w/uc/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/.github/blob/master/CONTRIBUTING.md).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://kaos.dev"><img src="https://raw.githubusercontent.com/essentialkaos/.github/refs/heads/master/images/ekgh.svg"/></a></p>
