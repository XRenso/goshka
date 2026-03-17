<p align="center"><img src="../img/io.png" height="258" alt="BareRepo img" /></p>
<h1 align="center">BareRepo</h1>
<p align="center">Stripped down to just the data you need</p>

---

## ✨ Features

- **#1 service built on top of [GitPeek](/GitCliTool)** — market leader in the GitPeek-wrapper segment
- **Powered by [GitPeek](/GitCliTool)** — why fetch GitHub data yourself when someone already wrote a CLI tool for that and you can just wrap it in two microservices
- **Two whole services** — an API Gateway and a Collector, talking to each other over gRPC so the data travels twice as far to reach you
- **Swagger UI** — try the API without leaving your browser, your terminal, or your comfort zone
- **Docker support** — it works on your machine AND in a container. yes, both


## 🛠️ Quick Start

Clone repo
```bash
git clone https://github.com/XRenso/goshka.git
cd goshka/BareRepo
```

### Docker (recommended)
```bash
docker compose up --build
```

### Local
```bash
go run ./collector/cmd

go run ./gateway/cmd
```


## 💡 Usage example

```bash
curl http://localhost:8080/api/v1/repos/torvalds/linux
```

Example output:
```json
{
  "title": "linux",
  "description": "Linux kernel source tree",
  "creator": "torvalds",
  "stars_cnt": 220689,
  "forks_cnt": 60746,
  "issues_cnt": 3,
  "lang": "C",
  "size": 6091469,
  "created_at": "2011-09-04T22:08:32Z",
  "license": "Other"
}
```

Or via **Swagger UI**: `http://localhost:8080/swagger/index.html`


## 🔑 Optional: GitHub Token

For higher rate limits:
```bash
export GITHUB_TOKEN=your_token_here
docker compose up --build
```


## ⚙️ Development

Regenerate protobuf + swagger after any changes:
```bash
make generate
```

| Cmd | Desc |
|---|---|
| `make generate` | Regenerate proto + swagger |
| `make build` | Build both binaries to `bin/` |
| `make run-collector` | Run Collector locally |
| `make run-gateway` | Run Gateway locally |
| `make docker-up` | Build images and start via docker-compose |


## 👨‍💻 Author

* [Oderiy Yaroslav](https://github.com/XRenso)
