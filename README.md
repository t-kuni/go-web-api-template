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
```

# Tests

Unit test

```
go test ./...
```

Feature test

```
go test -tags feature ./...
```

# Setting remote debug on GoLand

https://gist.github.com/t-kuni/1ecec9d185aac837457ad9e583af53fb#golnad%E3%81%AE%E8%A8%AD%E5%AE%9A

# See Database

http://localhost:8080

# Seeding data

```
docker-compose exec app sh
go run commands/seed/main.go
```


# Create Scheme

```
go run entgo.io/ent/cmd/ent init YourScheme
```

# 検討

- リクエストをバインド
- バリデーション
- 構造体の依存を全てポインタにする？
- マイグレーションの管理を切り出し
- featureテストでテストケース毎にデータの用意
- テスト対象の処理にコミットが含まれる場合の対処