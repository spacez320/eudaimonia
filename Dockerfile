FROM golang:1.22
WORKDIR /app
COPY . .
RUN go install
CMD ["go", "run", "main.go"]
