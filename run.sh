docker-compose down
docker-compose up -d --wait
go generate ./...
go build -o a.out .
./a.out
