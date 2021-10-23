# About

This repository is project skeleton for Go Web API application.

# Usage

```
git clone --depth 1 ssh://git@github.com/t-kuni/go-web-api-skeleton [ProjectName]
cd [ProjectName]
rm -rf .git 
```

# Generate 

```
go generate -x -tags wireinject ./...
```

# Boot local environment

```
dockr-compose up -d
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

```
DB_HOST=localhost DB_PORT=33060 DB_DATABASE=example_test go run commands/migrate/main.go
go test -tags feature ./...
```

# Setting remote debug on GoLand

https://gist.github.com/t-kuni/1ecec9d185aac837457ad9e583af53fb#golnad%E3%81%AE%E8%A8%AD%E5%AE%9A

# See Database

http://localhost:8080

# See SQL Log

```
docker-compose exec db tail -f /tmp/query.log
```

# Migration and Seeding

```
docker-compose exec app sh
DB_HOST=localhost DB_PORT=33060 go run commands/migrate/main.go
DB_HOST=localhost DB_PORT=33060 go run commands/seed/main.go
```

# Create Scheme

```
go run entgo.io/ent/cmd/ent init [EntityName]
```

# タスク

- [x] DIコンテナ導入（google/wire導入）
- [x] モックライブラリ導入
- [x] ent導入＆DB接続周り実装
- [x] Featureテスト作成
- [ ] リクエストをバインド
- [ ] バリデーション
- [ ] エラーハンドリング
- [x] 構造体の依存を全てポインタにする？
- [ ] マイグレーションの管理を切り出し
- [ ] 認証処理のモック化
- [ ] ロギング
- [x] featureテストでテストケース毎にデータの用意
- [x] テスト対象の処理にコミットが含まれる場合のテスト
- [ ] docker-composeをenvironmentsフォルダに移動
- [ ] レスポンスがJSONではない処理のテスト（例えばファイルのダウンロードなど）
- [ ] トランザクションが複数リクエストをまたがることはある？
- [ ] テストのカバレッジの可視化