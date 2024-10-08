FROM golang:1.23 AS base-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker

FROM alpine:3.14 AS release
WORKDIR /
COPY --from=base-stage /worker /worker
EXPOSE 3004
ENTRYPOINT [ "/worker" ]