build:
	docker-compose up --build
compose:
	docker-compose up
run:
	go run cmd/http/main.go