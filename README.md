<h1 align="center">
    <a name="top" title="dev-gadgets">git ⚙️</a><br/>Dev-Gadgets<br/> <sup><sub>Bootstrapping git project configurations</sub></sup>
</h1>
<div align=center>

[![CodeFactor][badge-codefactor]][link-codefactor]
[![License][badge-license]][link-license]
[![Repo analytics][badge-analytics]][link-analytics]
[![Tweet][badge-twitter]][link-twitter]

</div>
<div align="center">
    <p><strong>Be sure to <a href="#" title="star">⭐️</a> or <a href="#" title="fork">🔱</a> this repo if you find it useful! 😃</strong></p>
</div>

---

## Go quickstart (new CLI)

- Build: run `make build-go` (outputs `dist/dev-gadgets`).
- Run help: `go run ./cmd/dev-gadgets --help`.
- Dev container: open in VS Code with Dev Containers; recommended extensions auto-install.
- Catalog: edit `config/catalog.yaml` to add items and strategies.

---

## Install (one-liner)

You can install the latest release with:

```sh
curl -fsSL https://raw.githubusercontent.com/pirpedro/dev-gadgets/main/install.sh | bash
```

Options:

- Install to a custom dir:

```sh
PREFIX="$HOME/.local/bin" curl -fsSL https://raw.githubusercontent.com/pirpedro/dev-gadgets/main/install.sh | bash
```

- Pin a specific version:

```sh
VERSION="v0.1.0" curl -fsSL https://raw.githubusercontent.com/pirpedro/dev-gadgets/main/install.sh | bash
```

Releases publish ZIP archives per OS/arch using GoReleaser; the installer fetches the right one and places `dev-gadgets` in `~/.local/bin` by default.

<p align="center"><strong>Don't forget to <a href="#" title="star">⭐️</a> or <a href="#" title="fork">🔱</a> this repo! 😃<br/><sub>Assembled with <b title="love">❤️</b> in Rio de Janeiro.</sub></strong></p>

[badge-analytics]: https://img.shields.io/badge/repo%20analytics-public-informational?logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48IURPQ1RZUEUgc3ZnIFBVQkxJQyAiLS8vVzNDLy9EVEQgU1ZHIDEuMS8vRU4iICJodHRwOi8vd3d3LnczLm9yZy9HcmFwaGljcy9TVkcvMS4xL0RURC9zdmcxMS5kdGQiPjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCI+PHBhdGggZD0iTTIxIDhDMTkuNSA4IDE4LjcgOS40IDE5LjEgMTAuNUwxNS41IDE0LjFDMTUuMiAxNCAxNC44IDE0IDE0LjUgMTQuMUwxMS45IDExLjVDMTIuMyAxMC40IDExLjUgOSAxMCA5QzguNiA5IDcuNyAxMC40IDguMSAxMS41TDMuNSAxNkMyLjQgMTUuNyAxIDE2LjUgMSAxOEMxIDE5LjEgMS45IDIwIDMgMjBDNC40IDIwIDUuMyAxOC42IDQuOSAxNy41TDkuNCAxMi45QzkuNyAxMyAxMC4xIDEzIDEwLjQgMTIuOUwxMyAxNS41QzEyLjcgMTYuNSAxMy41IDE4IDE1IDE4QzE2LjUgMTggMTcuMyAxNi42IDE2LjkgMTUuNUwyMC41IDExLjlDMjEuNiAxMi4yIDIzIDExLjQgMjMgMTBDMjMgOC45IDIyLjEgOCAyMSA4TTE1IDlMMTUuOSA2LjlMMTggNkwxNS45IDUuMUwxNSAzTDE0LjEgNS4xTDEyIDZMMTQuMSA2LjlMMTUgOU0zLjUgMTFMNCA5TDYgOC41TDQgOEwzLjUgNkwzIDhMMSA4LjVMMyA5TDMuNSAxMVoiIGZpbGw9IiNmZmZmZmYiIC8+PC9zdmc+&maxAge=86400
[badge-codefactor]: https://img.shields.io/codefactor/grade/github/pirpedro/dotfiles?logo=codefactor&logoColor=white&cacheSeconds=300
[badge-license]: https://img.shields.io/github/license/renemarc/dotfiles.svg?logo=data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48IURPQ1RZUEUgc3ZnIFBVQkxJQyAiLS8vVzNDLy9EVEQgU1ZHIDEuMS8vRU4iICJodHRwOi8vd3d3LnczLm9yZy9HcmFwaGljcy9TVkcvMS4xL0RURC9zdmcxMS5kdGQiPjxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB4bWxuczp4bGluaz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94bGluayIgdmVyc2lvbj0iMS4xIiB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCI+PHBhdGggZD0iTTE3LjgsMjBDMTcuNCwyMS4yIDE2LjMsMjIgMTUsMjJINUMzLjMsMjIgMiwyMC43IDIsMTlWMThINUwxNC4yLDE4QzE0LjYsMTkuMiAxNS43LDIwIDE3LDIwSDE3LjhNMTksMkMyMC43LDIgMjIsMy4zIDIyLDVWNkgyMFY1QzIwLDQuNCAxOS42LDQgMTksNEMxOC40LDQgMTgsNC40IDE4LDVWMThIMTdDMTYuNCwxOCAxNiwxNy42IDE2LDE3VjE2SDVWNUM1LDMuMyA2LjMsMiA4LDJIMTlNOCw2VjhIMTVWNkg4TTgsMTBWMTJIMTRWMTBIOFoiIGZpbGw9IiNmZmZmZmYiIC8+PC9zdmc+Cg==&maxAge=86400
[badge-twitter]: https://img.shields.io/twitter/url/http/shields.io.svg?style=social&maxAge=86400
[link-analytics]: https://repo-analytics.github.io/pirpedro/dev-gadgets
[link-codefactor]: https://www.codefactor.io/repository/github/pirpedro/dev-gadgets/overview/main
[link-license]: LICENSE.txt
[link-twitter]: https://twitter.com/intent/tweet?text=Cross-platform%2C%20cross-shell%20%23dotfiles%20using%20%23chezmoi&url=https://github.com/pirpedro/dotfiles&via=pirpedro&hashtags=dev-gadgets
