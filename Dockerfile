FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o wantbox .
EXPOSE 8089
CMD ["./wantbox"]
