FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY catalog ./catalog
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/catalog ./catalog/cmd/server

FROM scratch

COPY --from=build /out/catalog /usr/local/bin/catalog

EXPOSE 8082

ENTRYPOINT ["catalog"]
