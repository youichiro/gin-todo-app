seed:
	psql -h 0.0.0.0 -p 5432 -U postgres -f postgresql/seed.sql go_todo_app_development
