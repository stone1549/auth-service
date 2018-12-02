FROM golang:1.8

ENV AUTH_SERVICE_ENVIRONMENT=DEV
ENV AUTH_SERVICE_REPO_TYPE=IN_MEMORY
ENV AUTH_SERVICE_TIMEOUT=60
ENV AUTH_SERVICE_PORT=3333
ENV AUTH_SERVICE_PG_URL=NA
ENV AUTH_SERVICE_INIT_DATASET=/go/src/github.com/stone1549/auth-service/data/small_set.json
ENV AUTH_SERVICE_TOKEN_SECRET=SECRET!

WORKDIR /go/src/github.com/stone1549/auth-service/
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["auth-service"]

