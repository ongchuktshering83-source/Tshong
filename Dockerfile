FROM golang:1.22-alpine

WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/

WORKDIR /app/backend

RUN go mod download

COPY backend/ ./

WORKDIR /app

COPY frontend/ ./frontend/

WORKDIR /app/backend

RUN go build -o tshongmart .

EXPOSE 8080

CMD ["./tshongmart"]