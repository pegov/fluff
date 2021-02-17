# Fluff
Fluff is an url-shortener service with one-time self-destructing links.

fluff-go and fluff-py are just different implementations.

## Dependencies
- docker-compose
## Usage
```bash
./scripts/fluff-go.sh up --build -d
# or
./scripts/fluff-py.sh up --build -d
```

## API
Visit to be redirected to url
```
GET /{key}
```
Create redirect
```
POST /api/links
{
  key: <string>, // can be null
  url: <string>,
}
```
key - if null, the key will be generated automatically for you

url - if it doesn't start with "http://" or "https://", then server will automatically add "https://" at the start of url. 

response
```
{
  key: <string>,
  url: <string>,
}
```

