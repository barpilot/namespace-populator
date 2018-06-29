FROM golang:1.10
WORKDIR /go/src/github.com/barpilot/namespace-populator/
COPY . /go/src/github.com/barpilot/namespace-populator/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/

FROM scratch
COPY --from=0 /go/src/github.com/barpilot/namespace-populator/app /bin/app
ENTRYPOINT ["/bin/app"]
