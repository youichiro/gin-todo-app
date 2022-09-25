# go-todo-app

## run migrate
for development db

```sh
# run up migration
docker compose run --rm migrate make up

# run down migration
docker compose run --rm migrate make down
```

for test db

```sh
# run up migration
GO_ENV=test docker compose run --rm migrate make up

# run down migration
GO_ENV=test docker compose run --rm migrate make down
```


## create migrate files
```sh
docker compose run --rm migrate make create ARG=<seq>
```

## run test
```sh
cd src
go test -v ./...

# 特定のパッケージを指定する
go test -v ./path/to/package
```

## run server
```sh
cd src/cmd
go run main.go
curl http://localhost:8080
```

```
{
  "message": "hello world!"
}
```

## generate models
https://github.com/volatiletech/sqlboiler

```sh
cd src
sqlboiler psql
```

## OpenAPI
https://stoplight.io/studio で編集する

モックサーバーを立てる
```sh
docker compose up -d prism
curl http://localhost:8081/tasks
```
