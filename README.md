# Fluff
Сокращатель ссылок

## API
- `GET /api/links/{short}` - Получить информацию о ссылке
- `POST /api/links` - Создать новую ссылку
	- `curl -X POST -H "Content-Type: application/json" -d '{"long": "https://vk.com"}' http://127.0.0.1:8080/api/links`
- `DELETE /api/links/{short}` - Удалить ссылку из базы

## Redirect
- `GET /{short}` - Редирект на сокращённую ссылку

## Deploy
```
docker build -t fluff-go .
docker run -it -p 8080:8080 fluff-go
```