<p align="center"><img src="../img/rubick.png" height="258" alt="Gitpeek img" /> </p>
<h1 align="center">GitPeek</h1>
<p align="center">CLI tool staring at repos so you don't have to</p>

---


## ✨ Features
- **Repository**: owner/repo
- **Name** & **Description**
- **Stars** ★, **Forks** ⑂, **Open Issues** ⚠️
- **Primary Language**
- **Repository Size**
- **Created Date**
- **License**


## 🛠️ Quick Start

Clone repo
```bash
git clone https://github.com/XRenso/goshka.git
cd goshka/GitCliTool
```

Build and run

```bash
go build
./gitpeek <repo>
```


## 💡 Usage example

```bash
# full url
./gitpeek https://github.com/torvalds/linux

# short
./gitpeek torvalds/linux

# without https://
./gitpeek github.com/torvalds/linux
```

Example output:
```bash

  torvalds/linux
  Linux kernel source tree

 ────────────────────────────────────────
  ★ 220689  ⑂ 60746   C
  issues 3     size 6091469 KB  Other
  created 2011-09-04
 ────────────────────────────────────────

```


## 🔑 Optional: GitHub Token
For higher rate limits:
```bash
export GITHUB_TOKEN=your_token_here
./gitcli torvalds/linux
```


## 👨‍💻 Author

* [Oderiy Yaroslav](https://github.com/XRenso)
