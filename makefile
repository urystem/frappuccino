run:
	make clean
	docker-compose up --build

cleanRun:
	make build
	docker-compose up

build:
	make clean
	docker-compose build
	make prune

clean:
	make stop
	make prune

prune:
	docker image prune -f

stop:
	docker-compose down -v