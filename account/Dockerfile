FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY account ./account
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/account ./account/cmd/server

FROM scratch

COPY --from=build /out/account /usr/local/bin/account

EXPOSE 8081

ENTRYPOINT ["account"]
