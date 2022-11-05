# About

This repository is project template for Go Web API application.

# Usage 

```
cp .env.example .env
cp .env.feature.example .env.feature
go generate -x -tags wireinject ./...
docker compose up -d
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

# Migration and Seeding

```
docker compose exec app sh
go run commands/migrate/main.go
go run commands/seed/main.go
```

# Create Scheme

```
go run entgo.io/ent/cmd/ent init [EntityName]
```

# タスク

- [ ] goのバージョンアップ
- [x] DIコンテナ導入（google/wire導入）
- [x] モックライブラリ導入
- [x] ent導入＆DB接続周り実装
- [x] Featureテスト作成
- [ ] OpenAPIと連携
- [ ] リクエストをバインド
- [ ] バリデーション
- [ ] エラーハンドリング
  - [ ] スタックトレース
- [x] 構造体の依存を全てポインタにする？
- [ ] マイグレーションの管理を切り出し
- [ ] 認証処理のモック化
- [ ] ロギング
- [x] featureテストでテストケース毎にデータの用意
- [x] テスト対象の処理にコミットが含まれる場合のテスト
- [ ] docker composeをenvironmentsフォルダに移動
- [ ] レスポンスがJSONではない処理のテスト（例えばファイルのダウンロードなど）
- [x] トランザクションが複数リクエストをまたがることはある？
  - ある
- [ ] テストのカバレッジの可視化
- [ ] DBを参照してからトランザクションを開始するケースに対応できるか
- [ ] 現在日時のモック化
- [ ] DB接続のタイムゾーン
- [ ] 本番環境用コンテナ
