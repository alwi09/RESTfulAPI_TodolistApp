# to build image todolist (dockerfile).
sudo docker build -t alwi09/todolist-image:v1 .

## how to install mockery
go install github.com/vektra/mockery/v2@latest

## how to build binary file restfulAPI todo
go build -o myApp ./cmd/main.go
## cara run
./myApp