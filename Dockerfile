# build web
FROM node:20.10.0-alpine3.19 AS web-builder
RUN corepack enable

WORKDIR /web

COPY web/package.json web/package-lock.json ./
RUN npm ci 

COPY web ./
RUN npm run build

# build server
FROM ghcr.io/hybridgroup/opencv:4.9.0 as server-builder

ENV GOPATH /go

WORKDIR /app

COPY . /go/src/go-opencv-extractor

WORKDIR /go/src/go-opencv-extractor

ARG TARGETARCH TARGETOS

RUN go mod download

COPY --from=web-builder /web/dist /go/src/go-opencv-extractor/web/dist

RUN CGO_ENABLED=1 GO111MODULE=on GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o /build/extractor main.go

# run server
FROM ghcr.io/hybridgroup/opencv:4.9.0 as server-runner

RUN apt-get update -qq && apt-get install ffmpeg -y

COPY --from=server-builder /build/extractor /build/extractor
COPY --from=server-builder /go/src/go-opencv-extractor/db_sql/migrations /build/db_sql/migrations

WORKDIR /build

EXPOSE 8080

ENTRYPOINT ["/build/extractor"]
CMD ["serve", "--port", "8080"]
