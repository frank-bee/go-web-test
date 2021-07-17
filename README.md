# Just to try iris framework

## Swagger UI

Generate files via  `swag init`

http://localhost:8080/swagger/index.html

User data structure
```
{
	"firstname": "a",
	"lastname": "b",
	"age": 22,
	"email": "foo"
}
```

## Database

Run the test mongo db via `sudo docker run -d -p 27017:27017 -v ~/data:/data/db mongo`