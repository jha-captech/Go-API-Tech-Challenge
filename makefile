# ── Database ────────────────────────────────────────────────────────────────────

.PHONY: db_up
db_up:
	podman compose up postgres

.PHONY: db_up_d
db_up_d:
	podman compose up postgres -d

.PHONY: db_down
db_down:
	podman compose down postgres

# ── API ─────────────────────────────────────────────────────────────────────────

.PHONY: run_app
run_app:
	podman compose up

.PHONY: build_image
build_image:
	podman build -t jf-techchallenge .
