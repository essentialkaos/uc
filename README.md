<p align="center"><a href="#readme"><img src="https://gh.kaos.st/uc.svg"/></a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/uc"><img src="https://travis-ci.com/essentialkaos/uc.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/uc"><img src="https://goreportcard.com/badge/github.com/essentialkaos/uc"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-uc-master"><img alt="codebeat badge" src="https://codebeat.co/badges/fd8a50fa-575c-47ba-8c67-1dd2f3b437f7" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#usage">Usage</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`uc` is a simple utility for counting unique lines.

### Usage demo

[![demo](https://gh.kaos.st/uc-001.gif)](#usage-demo)

### Installation

#### From sources

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `uc` from scratch, make sure you have a working Go 1.12+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/uc
```

If you want to update `uc` to latest stable release, do:

```
go get -u github.com/essentialkaos/uc
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```bash
sudo yum install -y https://yum.kaos.st/get/$(uname -r).rpm
sudo yum install uc
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

### Usage

```
Usage: uc {options} file

Options

  --dist, -d            Show number of occurrences for every line
  --max, -m num         Max number of unique lines (default: 5000)
  --no-progress, -np    Disable progress output
  --no-progress, -np    Disable progress output
  --no-color, -nc       Disable colors in output
  --help, -h            Show this help message
  --version, -v         Show version

Examples

  uc file.txt
  Count unique lines in file.txt

  uc -d file.txt
  Show distribution for file.txt

  cat file.txt | uc -
  Count unique lines in stdin data

```

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>