FROM golang:1.16 as builder
WORKDIR /app
ARG NAME
ARG VERSION
RUN apt-get install -y curl \
  && curl -sL https://deb.nodesource.com/setup_16.x | bash - \
  && apt-get install -y nodejs 
COPY ./ ./
RUN make build

FROM ubuntu:bionic
USER 1000
COPY --from=builder /app/bin/konfig-manager /bin/
ENTRYPOINT ["/bin/konfig-manager"]