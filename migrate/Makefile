up:
	migrate -database ${POSTGRESQL_URL} -path /opt/migrate/db/migrations up
down:
	migrate -database ${POSTGRESQL_URL} -path /opt/migrate/db/migrations down
create:
	migrate create -ext sql -dir /opt/migrate/db/migrations -seq $(ARG)
