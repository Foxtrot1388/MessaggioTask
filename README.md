всё поднимается `docker-compose up -d` (при наличие make - `make run`)

сервис доступен по пути `http://localhost:8080`

swagger доступен по пути `http://localhost:8080/swagger/index.html`

`MessaggioTask.postman_collection.json` содержит коллекцию запросов для postman

Параметры для env:

	LOG_LEVEL=-4 `уровень логирования LevelDebug: -4, LevelInfo: 0, LevelWarn: 4, LevelError: 8`
	POSTGRES_USER=root `юзер pg`
	POSTGRES_PASSWORD=root `пароль pg`
	POSTGRES_DB=messaggio `база для сервиса`
	KAFKA_TOPIC=tasks `топик для кафки (включено автосоздание)`