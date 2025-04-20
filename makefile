run:
	clean
	docker-compose up --build

cleanRun:
	build
	docker-compose up

build:
	clean
	docker-compose build
	prune

clean:
	stop
	prune

prune:
	docker image prune -f

stop:
	docker-compose down -v