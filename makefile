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

.PHONY: build-image
build-image:
	podman build -t jf-techchallenge .
