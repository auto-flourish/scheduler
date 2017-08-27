build:
	docker-compose stop
	docker-compose build web
	docker-compose up
clean:
	docker-compose down