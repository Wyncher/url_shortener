FROM golang:alpine
WORKDIR /build
COPY ./go.mod ./go.sum ./main.go ./
RUN go mod download
ADD services ./services
ADD config ./config
ADD models ./models
ADD db ./db
CMD ["go", "run", "main.go"]
EXPOSE 5000
