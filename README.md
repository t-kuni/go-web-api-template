# Go Web API Template

This repository is project template for Go Web API application.

# Features

* [Onion Architecture](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/)
* [DI Container](https://github.com/samber/do)
* [Server generation from swagger](https://github.com/go-swagger/go-swagger)
* [Validator](https://github.com/go-playground/validator)
* [ORM](https://github.com/ent/ent)
* [Logging](https://github.com/sirupsen/logrus)
* [Error Handling (Stack trace)](https://github.com/rotisserie/eris)
* Seeder

# Requirements

* go 1.19+
  * [How to install and switch between multiple versions of golang](https://gist.github.com/t-kuni/4e23b59f16557d704974b1ce6b49e6bb)

# Usage

```
swagger generate server -A App -f ./swagger.yml --model-package=restapi/models
go run cmd/app-server/main.go --scheme=http --port 34567 --host 0.0.0.0
curl -i "http://localhost:34567"
curl -i "http://localhost:34567/companies"
curl -i "http://localhost:34567" -d "{\"description\":\"message $RANDOM\"}" -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
curl -i localhost:34567/1 -X DELETE -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json'
```

```
cp .env.example .env
cp .env.feature.example .env.feature
go generate -x ./...
docker compose up -d
```

DB Migration and Seeding

```
docker compose exec app sh
go run commands/migrate/main.go
go run commands/seed/main.go
```

Confirm

```
curl -XGET "http://localhost"
curl -XPOST -H "Content-Type: application/json" -d '{"name":"DUMMY", "age":50}' http://localhost/users
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
DB_DATABASE=example_test go run commands/migrate/main.go
gotestsum --hide-summary=skipped -- -tags feature ./...
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
- [ ] レスポンスがJSONではない処理のテスト（例えばファイルのダウンロードなど）
- [ ] 現在日時のモック化
- [ ] DB接続のタイムゾーン
- [x] 本番環境用コンテナ
- [ ] vscode用devcontainer定義
- [ ] coreファイルが残る問題
- [x] アクセスログミドルウェア
- [x] OpenAPIと連携
  - [x] make コマンドでコード生成できるようにする（＋ファイルのクリーンアップ）
  - [x] 既存の仕組みと統合する
  - [x] airに対応
  - [x] DB接続に対応
  - [x] テストに対応
  - [-] 共通エラーハンドラ対応
  - [x] バリデーションを試す
- [x] テストをmakeコマンドから実行できるようにする
- [x] テストのカバレッジの可視化
- [x] テストの前処理、後処理をリファクタ（txdbの初期化処理タイミングを変更）
- [x] テストのレコードの初期投入を見やすくする
- [x] CI
- [ ] CD
