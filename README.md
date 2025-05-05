# `lambda-protoc`

A wrapper for executing `protoc` generators as part of an AWS Lambda function.

The lambda expects to receive as input a JSON string representing the base64
encoded `google.protobuf.compiler.CodeGeneratorRequest` message and will return
a JSON string representing the base64 encoded
`google.protobuf.compiler.CodeGeneratorResponse`.

## Example

```dockerfile
FROM golang:1.24.1-bookworm AS build

ENV CGO_ENABLED=0 GOBIN=/go/bin

RUN go install -ldflags "-s -w" google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6
RUN go install -ldflags "-s -w" github.com/CGA1123/lambda-protoc@HEAD

FROM scratch
WORKDIR /
COPY --from=build --link /etc/passwd /etc/passwd
COPY --from=build --link --chown=root:root /go/bin/protoc-gen-go .
COPY --from=build --link --chown=root:root /go/bin/lambda-protoc .
USER nobody
ENTRYPOINT [ "/lambda-protoc", "/protoc-gen-go" ]
```
