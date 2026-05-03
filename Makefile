.PHONY: up down dev/front dev/back dev/supabase

up: dev/supabase
	@trap 'kill 0' SIGINT; \
	(cd front && vp dev) & \
	(cd back && go run main.go) & \
	wait

dev/supabase:
	supabase start

dev/front:
	cd front && vp dev

dev/back:
	cd back && go run main.go

down:
	supabase stop
