# build web
FROM node:22-alpine AS web-builder
RUN corepack enable

WORKDIR /web

COPY web/package.json web/package-lock.json ./
RUN npm ci --force

COPY web ./
RUN npm run build

# build server
FROM vkhobor/bullseye-gocv-multiplatform as server-builder

ENV GOPATH /go

WORKDIR /go/src/go-opencv-extractor

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG TARGETARCH
ARG TARGETOS

COPY --from=web-builder /web/dist ./web/dist

RUN CGO_ENABLED=1 GO111MODULE=on GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o /build/extractor main.go

# run server
FROM vkhobor/bullseye-gocv-multiplatform as server-runner

RUN apt-get update -qq && apt-get install ffmpeg -y

WORKDIR /build
COPY --from=server-builder /build/extractor ./
COPY --from=server-builder /go/src/go-opencv-extractor/db/migrations ./db/migrations

EXPOSE 7001

ENTRYPOINT ["/build/extractor"]
CMD ["serve", "--port", "7001", "--db", "/DATA/db.sqlite3","--blob-storage", "/DATA/blobs", "--log-folder", "/DATA/logs"]
