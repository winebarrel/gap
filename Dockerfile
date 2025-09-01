FROM golang:1.25 AS build

WORKDIR /src
COPY go.* /src/
RUN go mod download

COPY *.go /src/
COPY cmd/gap/*.go /src/cmd/gap/
ARG gap_VERSION
RUN CGO_ENABLED=0 go build -o gap -ldflags "-X main.version=${GAP_VERSION}" ./cmd/gap

FROM gcr.io/distroless/static

COPY --from=build /src/gap /

ENTRYPOINT ["/gap"]
