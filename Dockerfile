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

WORKDIR /go/src/go-opencv-extractor

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG TARGETARCH TARGETOS

COPY --from=web-builder /web/dist ./web/dist

RUN CGO_ENABLED=1 GO111MODULE=on GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o /build/extractor main.go

# run server
FROM ghcr.io/hybridgroup/opencv:4.9.0 as server-runner

RUN apt-get update -qq && apt-get install ffmpeg -y

WORKDIR /build
COPY --from=server-builder /build/extractor ./
COPY --from=server-builder /go/src/go-opencv-extractor/db/migrations ./db/migrations

EXPOSE 7000

ENTRYPOINT ["/build/extractor"]
CMD ["serve", "--port", "7000", "--db", "/DATA/db.sqlite3","--blob-storage", "/DATA/blobs"]"--db", "/DATA/db.sqlite3","--blob-storage", "/DATA/blobs"
