# rocker

Use `rocker` to store your Kubernetes manifests in your Docker registry.

```console
$ go get github.com/coollog/rocker/cmd/rocker
$ 
$ cat hello-deployment.yaml | rocker push gcr.io/my-gcp-project/k8s/hello-deployment
$ cat hello-service.yaml | rocker push gcr.io/my-gcp-project/k8s/hello-service
$ 
$ rocker pull gcr.io/my-gcp-project/k8s/hello-deployment | kubectl apply -f -
$ rocker pull gcr.io/my-gcp-project/k8s/hello-service | kubectl apply -f -
```

Use `rocker` to store anything else.

```console
$ go get github.com/coollog/rocker/cmd/rocker
$ 
$ echo 'Hello World' | rocker push gcr.io/my-gcp-project/hello-world
$ 
$ rocker pull gcr.io/my-gcp-project/hello-world 2> /dev/null
Hello World
```

`rocker` uses the [`go-containerregistry`](https://github.com/google/go-containerregistry) library.
