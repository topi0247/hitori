.PHONY: up down dev/front dev/back dev/supabase generate

DATABASE_URL ?= postgresql://postgres:postgres@localhost:54322/postgres

generate:
	cd back && PATH=$$PATH:$(shell go env GOPATH)/bin go generate ./usecase/...

up: dev/supabase
	@trap 'kill 0' SIGINT; \
	(cd front && vp dev) & \
	(cd back && DATABASE_URL=$(DATABASE_URL) go run main.go) & \
	wait

dev/supabase:
	supabase start

dev/front:
	cd front && vp dev

dev/back:
	cd back && DATABASE_URL=$(DATABASE_URL) go run main.go

down:
	supabase stop
