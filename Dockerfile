FROM golang:1.24.1 AS lint

RUN apt update
RUN apt install -y \
    python3-pip \
    git

RUN pip3 install pre-commit --break-system-packages
RUN pip3 install xmlformatter --break-system-packages
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.0.2

WORKDIR /app

COPY . /app

RUN git config --global http.sslverify false

RUN pre-commit run --all-files


FROM golang:1.24.1 AS build

WORKDIR /app

COPY . /app

RUN go mod download
#CGO_ENABLED=0 gives you a static linked binary
RUN CGO_ENABLED=0 go build


FROM alpine:latest AS runtime

WORKDIR /app

COPY --from=build /app/gitea_actions_prometheus_exporter gitea_actions_prometheus_exporter

COPY entrypoint.sh /app/entrypoint.sh
RUN touch .env
#dummy copy to ensure lint phase happens
COPY --from=lint /app/README.md /README.md

ENTRYPOINT ["/app/entrypoint.sh"]
