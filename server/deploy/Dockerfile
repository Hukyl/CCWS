FROM golang:1.23.1-alpine AS build_stage

RUN apk --no-cache add ca-certificates
WORKDIR /go/src/github.com/Hukyl/CCWS
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -o /api-server cmd/main.go

FROM scratch AS api_stage
COPY --from=build_stage /api-server /api-server
COPY --from=build_stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/api-server"]
