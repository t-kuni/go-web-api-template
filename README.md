# Go Web API Template

This repository is project template for Go Web API application.

# Features

* [Onion Architecture](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/)
* [DI Container](https://github.com/uber-go/fx)
* [Server generation from swagger](https://github.com/go-swagger/go-swagger)
* [Validator](https://github.com/go-playground/validator)
* [ORM](https://github.com/ent/ent)
* [Logging](https://github.com/sirupsen/logrus)
* [Error Handling (Stack trace)](https://github.com/rotisserie/eris)
* Dev Container
* Seeder
* AI Agent
  * [Cline](./how-to-use-cline.md)
  * [Claude Code](./how-to-use-claude-code.md)

# Setup

1. VSCodeまたはCursorで本リポジトリを開きます

2. `Ctrl + Shift + P` でコマンドパレットを開き `Dev Containers: Reopen in Container` を実行する

Dev Container上でエディタが開き直します  
[docker-compose.yml](./docker-compose.yml) で使用しているポートが既に利用されていると起動に失敗するので注意してください

3. Terminalを開き（`Ctrl + Shift + @`）、以下のコマンドを実行します

3-1. envファイルを生成する

```bash
cp .env.example .env
cp .env.testing.example .env.testing
```

3-2. 各種ファイルを生成する

```bash
make generate
```

3-3. DBを構築＆レコードを登録する

```bash
go run commands/migrate/main.go --reset
go run commands/seed/main.go
```

3-4. 疎通確認

```bash
curl -i "http://localhost/companies"
curl -i "http://localhost/companies/UUID-1/users"
curl -i "http://localhost/users"
curl -i "http://localhost/todos"
```

# AIにタスクを依頼する

* Claude Codeを使用する場合
  * [how-to-use-claude-code.md](./how-to-use-claude-code.md)
* Clineを使用する場合
  * [how-to-use-cline.md](./how-to-use-cline.md)

# テストを実行する

```bash
# テスト用のDBを構築する
DB_DATABASE=example_test go run commands/migrate/main.go --reset
# テスト実行
make test
```

# Setting remote debug on GoLand

https://gist.github.com/t-kuni/1ecec9d185aac837457ad9e583af53fb#golnad%E3%81%AE%E8%A8%AD%E5%AE%9A

# See Database

http://localhost:8080

# See SQL Log

```
docker compose exec db tail -f /tmp/query.log
```

# Create Scheme

```
go run entgo.io/ent/cmd/ent init [EntityName]
```

# Build Container for production

```
docker build --target prod --tag go-web-api-template .
```
