FROM golang:1.18-alpine as debug

# Instalar git
RUN apk update && apk upgrade && \
    apk add --no-cache git \
        dpkg \
        gcc \
        git \
        musl-dev

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/cosmtrek/air@latest
RUN echo "alias air='$(go env GOPATH)/bin/air'" >> /root/.bashrc


# Work Directory
WORKDIR /go/src/work

COPY . /go/src/work/

# compile application
RUN go build -o app

### execution debugger Delve ###
ENV DEBUG_MODE=false
COPY dlv.sh /
RUN chmod +x /dlv.sh
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh


# ENTRYPOINT ["/dlv.sh"]
ENTRYPOINT ["/wait-for-it.sh", "db_go:5432", "--", "/dlv.sh"]

