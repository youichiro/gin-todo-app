seed:
	psql -h localhost -p 5432 -U postgres -f postgresql/seed.sql go_todo_app_development
psql:
	psql -h localhost -p 5432 -U postgres go_todo_app_development
