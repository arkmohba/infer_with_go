FROM nvidia/cuda:12.9.2-cudnn-devel-ubuntu24.04

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    curl \
    git \
    unzip \
    python3-pip \
    libopencv-dev \
    && rm -rf /var/lib/apt/lists/*


ENV GO_VERSION=1.26.4
RUN curl -fsSL "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -xzC /usr/local
ENV PATH=$PATH:/usr/local/go/bin

RUN mkdir -p /tmp/build
RUN git clone https://github.com/hybridgroup/gocv.git /tmp/build/gocv \
    && cd /tmp/build/gocv \
    && sed -i 's/sudo //g' Makefile \
    && make install

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED="1"
