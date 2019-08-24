ARG GOLANG_VERSION=1.12.9

FROM golang:${GOLANG_VERSION}-stretch as base
LABEL maintainer="Nick Calibey"
EXPOSE 4000
WORKDIR /github.com/ncalibey/hackernews-go

COPY ./go.mod go.mod
COPY ./go.sum go.sum
RUN go mod download

##############################################################################
# Builder Stage ##############################################################
FROM base as builder
COPY ./cmd cmd
COPY ./internal internal
RUN go build -o /bin/graphql-server cmd/server/main.go

##############################################################################
# Release Stage ##############################################################
FROM golang:${GOLANG_VERSION}-stretch as release

COPY --from=builder /bin/graphql-server /bin/graphql-server
CMD ["/bin/graphql-server"]
