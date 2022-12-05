Invocation History Extension
===

[![codecov](https://codecov.io/gh/michimani/invocation-history-extension/branch/main/graph/badge.svg?token=6TB4W4ZUJ0)](https://codecov.io/gh/michimani/invocation-history-extension)

This is a Extension for AWS Lambda Function that records history of invocation at the same runtime environment.

- Records in memory the AWS Request IDs of Lambda functions invoked at the same runtime execution environment along with the time of the invocation.
- Listen on localhost (default port: 1203) for an http server that returns a list of invocations executed up to that point.

# Usage

## Download or build extension's binary file

First, prepare the binary file of the extension.  
You can either download it by specifying a released version, or clone this repository and build it on your end.

### Download extension binary from released assets

You can obtain pre-built extensions in zip format at the following URL.

```bash
https://github.com/michimani/invocation-history-extension/releases/download/${EXTENSION_VERSION}/extension.zip
```

Please click [here](https://github.com/michimani/invocation-history-extension/releases) to check the released versions.

### Build extension yourself

Cloning this repository and executing the following command will build the extension and generate `extension.zip` under the `bin` directory.

```bash
make build
```

## Use extension in your Lambda Function

This extension can be used in Lambda functions by deploying the extension as a Lambda Layer or by including a binary file in the container image.

### As the Lambda Layer

To use this extension as a Lambda Layer, execute the following AWS CLI (v2) command to publish a Lambda Layer.

```bash
aws lambda publish-layer-version \
--layer-name "invocation-history-extension" \
--zip-file  "fileb://extension.zip" \
--region <your region>
```

Then, use the `update-function-configuration` command to specify the value of `LayerVersionArn` to the Lambda function that wants to use this extension.

### Include in container image

Prepare a Dockerfile like the following and place the binary files of the extension in the `/opt/extensions` directory in the image.

```dockerfile
FROM public.ecr.aws/lambda/provided:al2 as build
ENV EXTENSION_VERSION 0.1.0
RUN yum install -y golang unzip
RUN go env -w GOPROXY=direct
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o /main
RUN mkdir -p /opt
ADD ./extension.zip ./
RUN unzip extension.zip -d /opt
RUN rm extension.zip

FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
COPY entry.sh /
RUN chmod 755 /entry.sh
RUN mkdir -p /opt/extensions
WORKDIR /opt/extensions
COPY --from=build /opt/extensions .
ENTRYPOINT [ "/entry.sh" ]
CMD ["/main"]
```

# IPC (Inter-Process Communication)

This extension starts an HTTP API server on runtime, listening on `localhost:1203`. You can call `GET /invocations` to get a list of AWS Request IDs of functions invoked in the same runtime environment (including currently running ones).

The port number can be changed by setting the environment variable `INVOCATION_HISTORY_EXTENSION_HTTP_PORT` to any value. The default is `1203`.

Please see [here](https://github.com/michimani/invocation-history-extension/tree/main/docs/ipc.yaml) for API server specs.

# Example for using this extension

See [_example](https://github.com/michimani/invocation-history-extension/tree/main/_example) for using this extension at the Lambda Function (Golang) using container image.

# License

[MIT](https://github.com/michimani/aws-lambda-api-go/blob/main/LICENSE)

# Author

[michimani210](https://twitter.com/michimani210)