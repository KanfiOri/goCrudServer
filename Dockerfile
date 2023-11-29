FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

COPY wait-for-db.sh .
RUN chmod +x wait-for-db.sh

EXPOSE 8080
CMD ["./wait-for-db.sh"]