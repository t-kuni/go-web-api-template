# About

This repository is project skeleton for Go CLI application.

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

# Setting remote debug on GoLand

https://gist.github.com/t-kuni/1ecec9d185aac837457ad9e583af53fb#golnad%E3%81%AE%E8%A8%AD%E5%AE%9A

# 検討

- DB接続に対応
- リクエストをバインド
- バリデーション
- interfaceを作るのが面倒。wire.Bindも面倒。