FROM --platform=$BUILDPLATFORM alpine AS downloader-opencv

RUN apk add --no-cache curl unzip

ARG OPENCV_VERSION="4.10.0"
ENV OPENCV_VERSION $OPENCV_VERSION

ARG OPENCV_FILE="https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip"
ENV OPENCV_FILE $OPENCV_FILE

RUN curl -Lo opencv.zip ${OPENCV_FILE} && \
    unzip -q opencv.zip && \
    mv /opencv-${OPENCV_VERSION} /opencv

FROM --platform=$BUILDPLATFORM alpine AS downloader-contrib

RUN apk add --no-cache curl unzip

ARG OPENCV_VERSION="4.10.0"
ENV OPENCV_VERSION $OPENCV_VERSION

ARG OPENCV_CONTRIB_FILE="https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip"
ENV OPENCV_CONTRIB_FILE $OPENCV_CONTRIB_FILE

RUN curl -Lo opencv_contrib.zip ${OPENCV_CONTRIB_FILE} && \
    unzip -q opencv_contrib.zip && \
    mv /opencv_contrib-${OPENCV_VERSION} /opencv_contrib

FROM --platform=linux/arm64 golang:1.23-rc-bullseye AS opencv-base-arm64

RUN apt-get update && apt-get install -y --no-install-recommends \
    git build-essential cmake pkg-config unzip libgtk2.0-dev \
    curl ca-certificates libcurl4-openssl-dev libssl-dev \
    libavcodec-dev libavformat-dev libswscale-dev \
    libjpeg62-turbo-dev libpng-dev libtiff-dev libdc1394-22-dev && \
    apt-get autoremove -y && apt-get autoclean -y

FROM --platform=linux/arm64 opencv-base-arm64 AS opencv-build-arm64

COPY --from=downloader-opencv /opencv /opencv
COPY --from=downloader-contrib /opencv_contrib /opencv_contrib

RUN cd /opencv && \
    mkdir build && cd build && \
    cmake -D CMAKE_BUILD_TYPE=RELEASE \
    -D CMAKE_INSTALL_PREFIX=/usr/local \
    -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib/modules \
    -D BUILD_TESTS=OFF \
    -D WITH_EIGEN=OFF \
    -D WITH_VTK=OFF \
    -D WITH_QT=OFF \
    -D BUILD_JPEG=ON \
    -D OPENCV_ENABLE_NONFREE=ON \
    -D BUILD_DOCS=OFF \
    -D BUILD_EXAMPLES=OFF \
    -D BUILD_TESTS=OFF \
    -D BUILD_PERF_TESTS=ON \
    -D BUILD_opencv_java=NO \
    -D BUILD_opencv_python=NO \
    -D BUILD_opencv_python2=NO \
    -D BUILD_opencv_python3=NO \
    -D OPENCV_GENERATE_PKGCONFIG=ON .. && \
    make -j $(nproc --all) && \
    make preinstall && make install && ldconfig && \
    cd / && rm -rf opencv*

FROM --platform=linux/amd64 golang:1.23-rc-bullseye AS opencv-base-amd64

RUN apt-get update && apt-get install -y \
    git build-essential cmake pkg-config unzip libgtk2.0-dev \
    curl ca-certificates libcurl4-openssl-dev libssl-dev \
    libavcodec-dev libavformat-dev libswscale-dev libtbb2 libtbb-dev \
    libjpeg62-turbo-dev libpng-dev libtiff-dev libdc1394-22-dev nasm && \
    rm -rf /var/lib/apt/lists/*

FROM --platform=linux/amd64 opencv-base-amd64 AS opencv-build-amd64

COPY --from=downloader-opencv /opencv /opencv
COPY --from=downloader-contrib /opencv_contrib /opencv_contrib

RUN cd /opencv && \
    mkdir build && cd build && \
    cmake -D CMAKE_BUILD_TYPE=RELEASE \
    -D WITH_IPP=OFF \
    -D WITH_OPENGL=OFF \
    -D WITH_QT=OFF \
    -D CMAKE_INSTALL_PREFIX=/usr/local \
    -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib/modules \
    -D OPENCV_ENABLE_NONFREE=ON \
    -D WITH_JASPER=OFF \
    -D WITH_TBB=ON \
    -D BUILD_JPEG=ON \
    -D WITH_SIMD=ON \
    -D ENABLE_LIBJPEG_TURBO_SIMD=ON \
    -D BUILD_DOCS=OFF \
    -D BUILD_EXAMPLES=OFF \
    -D BUILD_TESTS=OFF \
    -D BUILD_PERF_TESTS=ON \
    -D BUILD_opencv_java=NO \
    -D BUILD_opencv_python=NO \
    -D BUILD_opencv_python2=NO \
    -D BUILD_opencv_python3=NO \
    -D OPENCV_GENERATE_PKGCONFIG=ON .. && \
    make -j $(nproc --all) && \
    make preinstall && make install && ldconfig && \
    cd / && rm -rf opencv*

FROM opencv-build-${TARGETARCH} as bullseye-gocv-multiplatform

CMD ["opencv_version", "-b"]
