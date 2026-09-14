FROM --platform=$BUILDPLATFORM golang:1.27.1 AS builder
ARG TARGETOS TARGETARCH

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build .

FROM scratch
COPY --from=builder /go/src/app/quick-template /bin/quick-template
ENTRYPOINT ["/bin/quick-template"]
