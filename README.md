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

# Usage

```
cp .env.example .env
cp .env.feature.example .env.feature
go generate -x ./...
docker compose up -d
```

DB Migration and Seeding

```
docker compose exec app sh
make generate
go run commands/migrate/main.go --reset
go run commands/seed/main.go
```

Generate Server

```
swagger generate server -A App -f ./swagger.yml --model-package=restapi/models
```

Confirm

```
curl -i "http://localhost"
curl -i "http://localhost/companies"
curl -i "http://localhost/companies/UUID-1/users"
curl -i "http://localhost/users"
curl -i "http://localhost/todos"
curl -i "http://localhost" -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
curl -i localhost/1 -X DELETE -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```

# Tests

Unit test

```
go test ./...
```

Feature test  
https://localhost:8080 に接続し、`example_test`データベースを作成してから以下のコマンドを実行する

```
docker compose exec app sh
DB_DATABASE=example_test go run commands/migrate/main.go --reset
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

# タスク

- [ ] polluterからtestfixturesに載せ替え
- [ ] マイグレーションの管理を切り出し
- [ ] 認証処理のモック化
- [ ] 現在日時のモック化
- [ ] DB接続のタイムゾーン
- [ ] coreファイルが残る問題
- [ ] CD
