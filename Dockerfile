FROM golang:1.16 as builder
WORKDIR /app
ARG NAME
ARG VERSION
COPY ./ ./
RUN make build

FROM ubuntu:bionic
USER 1000
COPY --from=builder /app/bin/konfig-manager /bin/
ENTRYPOINT ["/bin/konfig-manager"]