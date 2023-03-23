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

Built with Golang 1.20+ (with modules).

Run `make` to see the development targets.

`md-slides` is built and tested by Github Actions. Releases are done manually every now and then as needed:

1. Push a new tag like `vX.Y.Z`
2. [Draft a new release](https://github.com/astromechza/md-slides/releases/new) for the tag
3. On your local machine, run `make build` and update the artifacts to the release
4. Add the git diff for good measure as release notes
5. Publish it

To update the Github pages site, do the following after a release:

1. Run `temp=$(mktemp -d); md-slides html -source SLIDES.md -target-dir ${temp}; echo ${temp}` to generate the content
2. Open a new branch/PR against the `gh-pages` branch
3. Copy the content back to the branch from `${temp}`
4. Check and merge the PR

### Who uses md-slides?

Mostly just me!

This is a personal project and tool. You are most welcome to use it too but development is sporadic and tightly tied to my own wants and needs.
