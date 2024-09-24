FROM cgr.dev/chainguard/go:latest AS gobuilder

# use static link build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/local/src/repo
COPY . /usr/local/src/repo
RUN go build ./cmd/aictl

FROM cgr.dev/chainguard/wolfi-base

RUN apk add --no-cache shadow
# requires uid 1001 for writing to /github/* directories in actions
RUN useradd -u 1001 -m ghactions

COPY --from=gobuilder /usr/local/src/repo/aictl /bin/aictl
ENTRYPOINT ["/bin/aictl"]
