FROM golang:latest

WORKDIR /app
ARG CACHEBUST=1
COPY . .

RUN go mod tidy
RUN go build -o backend

EXPOSE 8080
CMD ["./backend"]

