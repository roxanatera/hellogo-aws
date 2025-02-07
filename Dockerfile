FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./
COPY .env .

COPY main.go .

#Para instalar dependencias.
RUN go mod download
RUN go get github.com/joho/godotenv
RUN go get -u github.com/gin-gonic/gin 
RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/mysql

RUN go build -o bin .

ENTRYPOINT ["/app/bin"]