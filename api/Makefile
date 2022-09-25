sqlboiler:
	sqlboiler psql
test:
	go test ./... --shuffle=on -p=2 -cover -count=1
coverage:
	go test ./... --shuffle=on -p=2 -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html
