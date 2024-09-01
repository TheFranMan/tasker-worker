FROM golang:1.23 AS base-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM gcr.io/distroless/base-debian11 AS release
WORKDIR /
COPY --from=base-stage /app/taskWorker /taskWorker
EXPOSE 3000
ENTRYPOINT [ "taskWorker" ]