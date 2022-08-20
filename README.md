# web-service-gin
go web service practice

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
