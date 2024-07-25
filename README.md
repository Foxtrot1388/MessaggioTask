всё поднимается `docker-compose up -d` (при наличие make - `make run`)

сервис доступен по пути `http://localhost:8080`

swagger доступен по пути `http://localhost:8080/swagger/index.html`

`MessaggioTask.postman_collection.json` содержит коллекцию запросов для postman

Пример .env файла для docker-compose:

	LOG_LEVEL=-4
	POSTGRES_USER=root
	POSTGRES_PASSWORD=root
	POSTGRES_DB=messaggio
	KAFKA_TOPIC=tasks
