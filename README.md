# ezBookkeeping
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/mayswind/ezbookkeeping)](https://goreportcard.com/report/github.com/mayswind/ezbookkeeping)
[![Latest Release](https://img.shields.io/github/release/mayswind/ezbookkeeping.svg?style=flat)](https://github.com/mayswind/ezbookkeeping/releases)
[![Latest Build](https://img.shields.io/github/actions/workflow/status/mayswind/ezbookkeeping/build-snapshot.yml?branch=main)](https://github.com/mayswind/ezbookkeeping/actions)
[![Latest Docker Image Size](https://img.shields.io/docker/image-size/mayswind/ezbookkeeping.svg?style=flat)](https://hub.docker.com/r/mayswind/ezbookkeeping)
[![Docker Pulls](https://img.shields.io/docker/pulls/mayswind/ezbookkeeping)](https://hub.docker.com/r/mayswind/ezbookkeeping)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/mayswind/ezbookkeeping)

[![Recommend By HelloGitHub](https://api.hellogithub.com/v1/widgets/recommend.svg?rid=ded5af09da574ec1811ddb154f1b2093&claim_uid=LT7EZxeBukCnh0K)](https://hellogithub.com/en/repository/mayswind/ezbookkeeping)
[![Trending](https://trendshift.io/api/badge/repositories/12917)](https://trendshift.io/repositories/12917)

## Introduction
ezBookkeeping is a lightweight, self-hosted personal finance app with a user-friendly interface and powerful bookkeeping features. It helps you record daily transactions, import data from various sources, and quickly search and filter your bills. You can analyze historical data using built-in charts or perform custom queries with your own chart dimensions to better understand spending patterns and financial trends. ezBookkeeping is easy to deploy, and you can start it with just one single Docker command. Designed to be resource-efficient, it runs smoothly on devices such as Raspberry Pi, NAS, and MicroServers.

ezBookkeeping offers tailored interfaces for both mobile and desktop devices. With support for PWA (Progressive Web Apps), you can even [add it to your mobile home screen](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/add_to_home_screen.gif) and use it like a native app.

Live Demo: [https://ezbookkeeping-demo.mayswind.net](https://ezbookkeeping-demo.mayswind.net)

## Features
- **Open Source & Self-Hosted**
    - Built for privacy and control
- **Lightweight & Fast**
    - Optimized for performance, runs smoothly even on low-resource environments
- **Easy Installation**
    - Docker-ready
    - Supports SQLite, MySQL, PostgreSQL
    - Cross-platform (Windows, macOS, Linux)
    - Works on x86, amd64, ARM architectures
- **User-Friendly Interface**
    - UI optimized for both mobile and desktop
    - PWA support for native-like mobile experience
    - Dark mode
- **AI-Powered Features**
    - Receipt image recognition
    - Supports MCP (Model Context Protocol) for AI integration
- **Powerful Bookkeeping**
    - Two-level accounts and categories
    - Attach images to transactions
    - Location tracking with maps
    - Recurring transactions
    - Advanced filtering, search, visualization, and analysis
- **Localization & Globalization**
    - Multi-language and multi-currency support
    - Automatic exchange rates
    - Multi-timezone awareness
    - Custom formats for dates, numbers, and currencies
- **Security**
    - Two-factor authentication (2FA)
    - Login rate limiting
    - Application lock (PIN code / WebAuthn)
- **Data Import/Export**
    - Supports CSV, OFX, QFX, QIF, IIF, Camt.052, Camt.053, MT940, GnuCash, Firefly III, Beancount, and more

For a full list of features, visit the [Full Feature List](https://ezbookkeeping.mayswind.net/comparison/).

## Screenshots
### Desktop Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)

### Mobile Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)

## Installation
### Run with Docker
Visit [Docker Hub](https://hub.docker.com/r/mayswind/ezbookkeeping) to see all images and tags.

**Latest Release:**

    $ docker run -p8080:8080 mayswind/ezbookkeeping

**Latest Daily Build:**

    $ docker run -p8080:8080 mayswind/ezbookkeeping:latest-snapshot

### Install from Binary
Download the latest release: [https://github.com/mayswind/ezbookkeeping/releases](https://github.com/mayswind/ezbookkeeping/releases)

**Linux / macOS**

    $ ./ezbookkeeping server run

**Windows**

    > .\ezbookkeeping.exe server run

By default, ezBookkeeping listens on port 8080. You can then visit `http://{YOUR_HOST_ADDRESS}:8080/` .

### Build from Source
Make sure you have [Golang](https://golang.org/), [GCC](https://gcc.gnu.org/), [Node.js](https://nodejs.org/) and [NPM](https://www.npmjs.com/) installed. Then download the source code, and follow these steps:

**Linux / macOS**

    $ ./build.sh package -o ezbookkeeping.tar.gz

All the files will be packaged in `ezbookkeeping.tar.gz`.

**Windows**

    > .\build.bat package -o ezbookkeeping.zip

or

    PS > .\build.ps1 package -Output ezbookkeeping.zip

All the files will be packaged in `ezbookkeeping.zip`.

### Build Windows binary via GitHub Actions

仓库内提供了 Windows 自编译工作流，可在 GitHub 上自动构建 Windows 可执行包：

- **触发方式**
  - 推送 tag（如 `v1.0.1`）时自动运行，并创建 Release、附带 Windows zip。
  - 在 Actions 页选择 **Windows Build** 工作流，点击 **Run workflow** 可手动触发。
- **产物**
  - 每次运行会生成 `ezbookkeeping-<version>-windows.zip`（或带日期的 snapshot 包），可在该次运行的 **Artifacts** 中下载。
  - 包内包含 `ezbookkeeping.exe`、前端静态资源、配置与模板，解压后在 Windows 下运行 `ezbookkeeping.exe server run` 即可。

工作流文件：`.github/workflows/windows-build.yml`。

You can also build a Docker image. Make sure you have [Docker](https://www.docker.com/) installed, then follow these steps:

**Linux**

    $ ./build.sh docker

## Contributing
We welcome contributions of all kinds.

If you find a bug, please [submit an issue](https://github.com/mayswind/ezbookkeeping/issues) on GitHub.

If you would like to contribute code, you can fork the repository and open a pull request.

Improvements to documentation, feature suggestions, and other forms of feedback are also appreciated.

You can view existing contributors on the [Contributor Graph](https://github.com/mayswind/ezbookkeeping/graphs/contributors).

## Translating
Help make ezBookkeeping accessible to users around the world. We welcome help to improve existing translations or add new ones. If you would like to contribute a translation, please refer to the [translation guide](https://ezbookkeeping.mayswind.net/translating).

Currently available translations:

| Tag | Language | Contributors |
| --- | --- | --- |
| de | Deutsch | [@chrgm](https://github.com/chrgm) |
| en | English | / |
| es | Español | [@Miguelonlonlon](https://github.com/Miguelonlonlon), [@abrugues](https://github.com/abrugues), [@AndresTeller](https://github.com/AndresTeller), [@diegofercri](https://github.com/diegofercri) |
| fr | Français | [@brieucdlf](https://github.com/brieucdlf) |
| it | Italiano | [@waron97](https://github.com/waron97) |
| ja | 日本語 | [@tkymmm](https://github.com/tkymmm) |
| kn | ಕನ್ನಡ | [@Darshanbm05](https://github.com/Darshanbm05) |
| ko | 한국어 | [@overworks](https://github.com/overworks) |
| nl | Nederlands | [@automagics](https://github.com/automagics) |
| pt-BR | Português (Brasil) | [@thecodergus](https://github.com/thecodergus) |
| ru | Русский | [@artegoser](https://github.com/artegoser), [@dshemin](https://github.com/dshemin) |
| sl | Slovenščina | [@thehijacker](https://github.com/thehijacker) |
| ta | தமிழ் | [@hhharsha36](https://github.com/hhharsha36) |
| th | ไทย | [@natthavat28](https://github.com/natthavat28) |
| tr | Türkçe | [@aydnykn](https://github.com/aydnykn) |
| uk | Українська | [@nktlitvinenko](https://github.com/nktlitvinenko) |
| vi | Tiếng Việt | [@f97](https://github.com/f97) |
| zh-Hans | 中文 (简体) | / |
| zh-Hant | 中文 (繁體) | / |

## Documentation
1. [English](https://ezbookkeeping.mayswind.net)
1. [中文 (简体)](https://ezbookkeeping.mayswind.net/zh_Hans)

## License
[MIT](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
