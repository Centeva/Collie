FROM golang:alpine as build-env

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# These options create smaller go builds.
# CGO_ENABLED=0: builds staticly linked binaries (we are in a scratch image so no libs to dynamically link to)
# -ldflags=-w: Disables DWARF generation (breakpoints, etc)
# -ldflags=-s: Disables symbol table
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/collie

FROM scratch

COPY --from=build-env /go/bin/collie /go/bin/collie

ENTRYPOINT ["/go/bin/collie"]