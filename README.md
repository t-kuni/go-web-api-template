# About

This repository is project skeleton for Go CLI application.

# Usage

```
git clone --depth 1 ssh://git@github.com/t-kuni/go-cli-app-skeleton [ProjectName]
cd [ProjectName]
rm -rf .git 
```

# Generate 

```
wire gen ./wire
```

# Generate Mock

```
mockgen -source=./domain/infrastructure/api/binanceApi.go -package=api -destination=./domain/infrastructure/api/binanceApi_mock.go
```

# Build

```
docker build --tag example-container .
```

# Run

```
docker run example-container example
```

# Setting remote debug on GoLand

https://gist.github.com/t-kuni/1ecec9d185aac837457ad9e583af53fb#golnad%E3%81%AE%E8%A8%AD%E5%AE%9A

# 検討

- wireやmockgenのコマンドを簡略化
- DB接続に対応
- バインド
- バリデーション
- 