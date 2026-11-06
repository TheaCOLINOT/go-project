FROM golang:1.27 AS build

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM scratch
COPY --from=build /server /server
ENTRYPOINT ["/server"]
