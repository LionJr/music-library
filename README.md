# music-library

1. git clone https://github.com/LionJr/music-library
2. docker run --name container_name \               
   -e POSTGRES_USER=username \
   -e POSTGRES_PASSWORD=password \
   -e POSTGRES_DB=db_name \
   -p port:port \
   -d postgres
3. Create .env file with:
   - HTTP_HOST
   - HTTP_PORT
   - DB_HOST
   - DB_PORT
   - DB_USER
   - DB_PASSWORD
   - DB_NAME
   - EXTERNAL_API_URL (external API to add songs)
4. go run cmd/main.go
