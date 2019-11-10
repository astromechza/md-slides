<h1 align="center">md-slides</h1>

<p align="center">
  <img src="https://img.shields.io/travis/astromechza/md-slides" />
  <img src="https://img.shields.io/github/downloads/astromechza/md-slides/total" />
  <img src="https://img.shields.io/badge/licence-MIT-green" />
  <a href="https://twitter.com/benmeier_">
    <img alt="Twitter: benmeier_" src="https://img.shields.io/twitter/follow/benmeier_.svg?style=social" target="_blank" />
  </a>
</p>

> `md-slides` is a single binary for rendering and presenting Markdown-based slides as html.

- Supports all traditional markdown features
- Support for various alignment, sizing, and resolution options
- Support for "paged" vs "scrolling" modes
- Decent support for rendering to PDF via Chrome/Chromium

See [astromechza.github.io/md-slides/](https://astromechza.github.io/md-slides/) for an example of the standalone output and a more in depth look at the features.

### Install

The following command will find and download the latest release from the Github release page:

```
$ curl https://raw.githubusercontent.com/astromechza/md-slides/master/install.sh | INSTALL_DIRECTORY=~/bin sh
```

Binaries are available for macOS, Linux, and Windows.

### Development

Built with Golang 1.11+ (with modules).

Run `make` to see the development targets.

### Who uses md-slides?

Mostly just me!

This is a personal project and tool. You are most welcome to use it too but development is sporadic and tightly tied to my own wants and needs.
