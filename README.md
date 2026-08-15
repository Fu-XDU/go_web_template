# go_web_template

Go + Gin 后端、Vue 3 + Vite 前端的 Web 项目模版。默认不连接 MySQL。

<!-- scaffold-only:start -->
## 生成新项目

别人安装到 `GOPATH/bin` 后，在任意目录执行 `create-go-web` 即可（会询问项目名）：

```sh
go install github.com/Fu-XDU/go_web_template/cmd/create-go-web@latest
create-go-web
```

请先把本仓库推到 GitHub，并保证根目录 `go.mod` 的 module 是 `github.com/Fu-XDU/go_web_template`（与仓库地址一致），然后打 tag：

```sh
git tag v0.0.1
git push origin main --tags
```

本机开发：

```sh
make create          # 打包模版、编译并启动
make install-create  # 安装到 $(go env GOPATH)/bin
```

`GOPATH/bin` 需要在 `PATH` 里：

```sh
export PATH="$(go env GOPATH)/bin:$PATH"
```

项目名必须以字母开头，只能包含字母、数字、`.`、`_`、`-`。不合法或目录已存在会重新询问。模版已打进二进制，安装后不需要再带源码仓库。

### 修改模版后如何更新

`create-go-web` 用的是打包进二进制的 `cmd/create-go-web/template.zip`，只改源码不会自动生效。

1. 重新打包并安装到本机：

```sh
make pack-create
make install-create
```

`make create` 会先 pack 再编译，本机试生成时用这个即可。

2. 发布给别人：把更新后的 `template.zip` 一并提交，再打**新的** tag 并推送（`@latest` 指向最新 semver tag）：

```sh
git add .
git commit -m "update template snapshot"
git tag v0.0.2
git push origin main --tags
```

3. 别人更新已安装的命令：

```sh
go install github.com/Fu-XDU/go_web_template/cmd/create-go-web@latest
```
<!-- scaffold-only:end -->

## 环境要求

- Go 1.25+
- Node.js 22.18+ 或 24.12+
- 可选：Docker / Docker Compose

首次使用：

```sh
cp .env.example .env
```

## 本地开发

```sh
make dev
```

会同时启动：

| 进程 | 地址 | 用途 |
| --- | --- | --- |
| Vite | http://127.0.0.1:5173/ | 前端热更新。`/v1` 会代理到后端 |
| Go | http://127.0.0.1:1423/ | API。若已有构建产物，也会提供 `/v1/web/` 静态页 |

开发时请打开 **5173**。1423 上的 `/v1/web/` 来自上次 `web/dist` 或 `assets/web/v1`，不是热更新页面。

示例接口：`GET http://127.0.0.1:1423/v1/hello/`

## 本地生产构建并运行

```sh
make start
```

会构建前后端、把前端同步到 `assets/web/v1`，再启动二进制。浏览器打开：

http://127.0.0.1:1423/v1/web/

根路径 `/` 会跳转到 `/v1/web/`。

## Make 目标

| 目标 | 说明 |
| --- | --- |
| `make dev` | `go run` + Vite 开发服务器 |
| `make start` | 构建并运行，页面在 `/v1/web/` |
| `make build` | 构建后端 + 前端，并同步静态资源 |
| `make build-backend` | 只构建 Go 二进制到 `bin/go_web_template` |
| `make build-frontend` | 只构建 Vue 到 `web/dist` |
| `make sync-frontend` | 把 `web/dist` 复制到 `assets/web/v1` |
| `make test` | 运行 Go 测试 |
| `make clean` | 删除构建产物 |
| `make docker` | 构建镜像 `go_web_template:$(VERSION)` |

构建时可覆盖版本信息：

```sh
make build-backend VERSION=1.0.0
```

`--help` / `--version` 以及启动日志会打印 `version`、`commit`、`build time`。

## Docker

```sh
cp .env.example .env
make docker
docker compose up -d
```

`docker-compose.yaml` 依赖当前目录的 `.env`。

## 启用 MySQL

默认关闭。`database/`、`model.User` 已预留，按下面三步打开即可。

`MYSQL_PASSWD` 和 `MYSQL_DB` 在 common flags 里是必填项，启用后必须写入 `.env`，否则进程起不来。

1. 在 `.env` 中取消注释并填写：

```sh
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWD=your_password
MYSQL_DB=go_web_template
```

2. 打开 `server.go` 里三处 `MYSQL支持` 注释：

```go
// 1/3  import
"github.com/Fu-XDU/go_web_template/database"

// 2/3  flags
app.Flags = append(app.Flags, mingfuflags.MysqlFlags...)

// 3/3  启动时建连
if err = database.Setup(); err != nil {
    return
}
```

3. 按需在 `database/mysql/mysql.go` 的 `initMysql` 里打开建表，例如：

```go
err = db.AutoMigrate(&model.User{})
```

并取消该文件顶部 `github.com/Fu-XDU/go_web_template/model` 的 import 注释。
