# rocker

Use your Docker registry to store anything. THIS IS AN EXPERIMENTAL PROTOTYPE.

![experimental](https://img.shields.io/badge/stability-experimental-red.svg)
[![Gitter chat](https://badges.gitter.im/coollog/rocker.png)](https://gitter.im/coollog/rocker)

## What can I do with `rocker`?

### Kubernetes manifests

Use `rocker` to store your Kubernetes manifests in your Docker registry.

```console
$ cat hello-deployment.yaml | rocker push gcr.io/my-gcp-project/k8s/hello-deployment
$ cat hello-service.yaml | rocker push gcr.io/my-gcp-project/k8s/hello-service
$ 
$ rocker pull gcr.io/my-gcp-project/k8s/hello-deployment | kubectl apply -f -
$ rocker pull gcr.io/my-gcp-project/k8s/hello-service | kubectl apply -f -
```

### Anything else

Use `rocker` to store anything else.

```console
$ echo 'Hello World' | rocker push gcr.io/my-gcp-project/hello-world
$ 
$ rocker pull gcr.io/my-gcp-project/hello-world 2> /dev/null
Hello World
```

## Usage

### 1) Install `rocker`.

#### Linux

```bash
curl -Lo ./rocker https://storage.googleapis.com/rocker-download/rocker-linux-amd64 && \
    chmod +x ./rocker && sudo mv ./rocker /usr/local/bin
```

#### macOS

```bash
curl -Lo ./rocker https://storage.googleapis.com/rocker-download/rocker-darwin-amd64 && \
    chmod +x ./rocker && sudo mv ./rocker /usr/local/bin
```

#### Windows

Download the latest Windows build: https://storage.googleapis.com/rocker-download/rocker-windows-amd64.exe

#### Build from source

```bash
go get -u github.com/coollog/rocker
# `rocker` will be at `$GOPATH/bin/rocker`
```

### How it works

`rocker` simply pushes an image with one layer containing a single file `/data` that contains your data.

`rocker` uses the [`go-containerregistry`](https://github.com/google/go-containerregistry) library.
