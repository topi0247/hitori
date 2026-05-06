.PHONY: up down dev/front dev/back dev/supabase generate generate/sqlc

DATABASE_URL ?= postgresql://postgres:postgres@localhost:54322/postgres
ALLOWED_ORIGIN ?= http://localhost:5173

generate:
	cd back && PATH=$$PATH:$(shell go env GOPATH)/bin go generate ./usecase/...

generate/sqlc:
	cd back && go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

up: dev/supabase
	@trap 'kill 0' SIGINT; \
	(cd front && vp dev) & \
	(cd back && env $$(grep -v '^\#' .env | xargs) go run main.go) & \
	wait

dev/supabase:
	supabase start

dev/front:
	cd front && vp dev

dev/back:
	cd back && env $$(grep -v '^\#' .env | xargs) go run main.go

down:
	supabase stop
