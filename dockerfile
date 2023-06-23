# Defines the base image used in image creation.
FROM golang:1.20.4-alpine AS builder

# Adds metadata labels to the image
LABEL author=Mohammad_Alwi_Irfani
LABEL email=alwi.irfani1927@gmail.com
LABEL github=https://github.com/alwi09/todolist_gin_gorm

# Executes the apk command to update the Alpine Linux packages and install Git.
RUN apk update && apk add --no-cache git

# Sets the current working directory to "/home" inside the image.
WORKDIR /home

# Copies the go.mod and go.sum files from the local directory to the current working directory inside the image.
COPY go.mod go.sum ./

# command to clean up and update the Go application dependencies based on the go.mod file.
RUN go mod tidy

# command to download all the dependencies required by the Go application.
RUN go mod download

# Copies all the files from the local directory to the current working directory inside the image.
COPY . .

# Executes the go build command to build the Go application. The -o option is used to specify the output binary file name, in this case, "todolistApp". The argument ./cmd/main.go indicates the location of the main source file of the application.
RUN go build -o todolistApp ./cmd/main.go 

# Defines the second base image for the final running container.
FROM alpine:3.15

#  Executes the apk command to install Curl inside the final container. 
RUN apk --no-cache add curl

# Sets the current working directory to "/home" inside the final container.
WORKDIR /home

# Copies the built executable binary "todolistApp" from the builder image to the current working directory inside the final container.
COPY --from=builder /home/todolistApp .

# Copies the migrations directory from the builder image to the migrations directory inside the final container.
COPY --from=builder /home/internal/database/migrations/ ./internal/database/migrations/

# Exposes port 1234 in the final container.
EXPOSE 1234

# Specifies the default command to be executed when the container starts. In this case, the command "./todolistApp" will be run inside the container.
CMD ["./todolistApp"]