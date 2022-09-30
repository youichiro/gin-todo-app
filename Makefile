seed:
	psql -h localhost -p 5432 -U postgres -f db/seed.sql go_todo_app_development
psql:
	psql -h localhost -p 5432 -U postgres go_todo_app_development
validation-proxy:
	prism proxy openapi/openapi.yaml http://localhost:8080 --errors
