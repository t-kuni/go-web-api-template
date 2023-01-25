# Go Web API Template

This repository is project template for Go Web API application.

# Features

* [Onion Architecture](https://jeffreypalermo.com/2008/07/the-onion-architecture-part-1/)
* [DI Container](https://github.com/samber/do)
* [Validator](https://github.com/go-playground/validator)

# Requirements

* go 1.19+
  * [How to install and switch between multiple versions of golang](https://gist.github.com/t-kuni/4e23b59f16557d704974b1ce6b49e6bb)

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

# タスク

- [ ] OpenAPIと連携
- [ ] エラーハンドリング
  - [ ] スタックトレース
- [ ] マイグレーションの管理を切り出し
- [ ] 認証処理のモック化
- [ ] ロギング
- [ ] docker composeをenvironmentsフォルダに移動
- [ ] レスポンスがJSONではない処理のテスト（例えばファイルのダウンロードなど）
- [ ] テストのカバレッジの可視化
- [ ] DBを参照してからトランザクションを開始するケースに対応できるか
- [ ] 現在日時のモック化
- [ ] DB接続のタイムゾーン
- [ ] 本番環境用コンテナ
- [ ] vscode用devcontainer定義
- [ ] テストの前処理、後処理をリファクタ（txdbの初期化処理タイミングを変更）
- [ ] レスポンスを整理（Bindエラー、Validationエラー）
- [ ] coreファイルが残る問題
