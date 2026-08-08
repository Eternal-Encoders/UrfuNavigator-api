# UrfuNavigator-api

List of requiered env's:
* MODE - Dev или Prod
* PORT - Порт для запуска
* CORS - Исключения из CORS'ов
* DEFAULT_PATH - Стандартный путь API
* DATABASE_URI - URI для подключения к MongoDB
* DATABASE_COLLECTION - Название коллекции
* BUCKET_ENDPOINT - Эндпоинт с API Minio
* BUCKET_ACCESS_KEY - Id ключа Minio
* BUCKET_SECRET_KEY - Secret ключа Minio
* BUCKET_NAME - Название бакета Minio

# Run

Запуск Dev:

```bash
air -c .air.windows.conf
```