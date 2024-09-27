FROM cgr.dev/chainguard/go:latest AS gobuilder
# use static link build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/local/src/repo
COPY . /usr/local/src/repo
RUN go build ./cmd/aictl

FROM cgr.dev/chainguard/wolfi-base
COPY --from=gobuilder /usr/local/src/repo/aictl /bin/aictl

# TIPS: make sure to run this image with root user for github actions
# github actions requires root or uid 1001 to have an access to /github/* dirs
ENTRYPOINT ["/bin/aictl"]
