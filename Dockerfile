FROM golang:1.14-alpine as builder
RUN mkdir build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
ARG VERSION=dev-docker
COPY . .
RUN CGO_ENABLED=0 go build -v -trimpath -ldflags "-X 'main.version=${VERSION}'" -o bin/tabia ./cmd/tabia

FROM alpine
LABEL maintainer="marco.franssen@philips.com"
RUN mkdir -p /app/data
WORKDIR /app
VOLUME [ "/app/data" ]
ENV TABIA_BITBUCKET_API=\
    TABIA_BITBUCKET_USER=\
    TABIA_BITBUCKET_TOKEN=\
    TABIA_GITHUB_USER=\
    TABIA_GITHUB_TOKEN=
COPY --from=builder build/bin/tabia .
ENTRYPOINT [ "/app/tabia" ]
