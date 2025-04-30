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
* Cline

# Setup

1. VSCodeで本リポジトリを開くと画面右下に以下のポップアップが表示されるので `Reopen in Container` をクリックする  
   （Clineも同時にインストールされます）
   
![image](https://github.com/user-attachments/assets/fc32e2ec-ffbc-403b-a14f-5abd88e26d87)

2. Terminalを開き（`Ctrl + Shift + @`）、以下のコマンドを実行します

2-1. envファイルを生成する

```bash
cp .env.example .env
cp .env.feature.example .env.feature
```

2-2. 各種ファイルを生成する

```bash
make generate
```

2-3. DBを構築＆レコードを登録する

```bash
go run commands/migrate/main.go --reset
go run commands/seed/main.go
```

2-4. 疎通確認

```bash
curl -i "http://localhost/companies"
curl -i "http://localhost/companies/UUID-1/users"
curl -i "http://localhost/users"
curl -i "http://localhost/todos"
```

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
