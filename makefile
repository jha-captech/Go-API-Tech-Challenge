# ── Database ────────────────────────────────────────────────────────────────────

.PHONY: db_up
db_up:
	docker-compose up postgres

.PHONY: db_up_d
db_up_d:
	docker-compose up postgres -d

.PHONY: db_down
db_down:
	docker-compose down postgres

# ── API ─────────────────────────────────────────────────────────────────────────

.PHONY: run_app
run_app:
	docker-compose up

.PHONY: build_image

build_image:
	docker build -t jf-techchallenge .

.PHONY: run
run:
	DATABASE_NAME=coursesDB DATABASE_USER=courses-db-user DATABASE_PASSWORD=courses-db-password DATABASE_HOST=localhost DATABASE_PORT=5432 go run cmd/api/main.go