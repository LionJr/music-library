# music-library

1. git clone https://github.com/LionJr/music-library
2. docker run --name container_name \               
   -e POSTGRES_USER=username \
   -e POSTGRES_PASSWORD=password \
   -e POSTGRES_DB=db_name \
   -p port:5432 \
   -d postgres
3. Add Postgres configs to .env file (ALSO ADD EXTERNAL_API_URL)
4. go run cmd/main.go
