Invocation History Extension
===

[![codecov](https://codecov.io/gh/michimani/invocation-history-extension/branch/main/graph/badge.svg?token=6TB4W4ZUJ0)](https://codecov.io/gh/michimani/invocation-history-extension)

This is a Extension for AWS Lambda Function that records history of invocation at the same runtime environment.

- Records in memory the AWS Request IDs of Lambda functions invoked at the same runtime execution environment along with the time of the invocation.
- Listen on localhost (default port: 1203) for an http server that returns a list of invocations executed up to that point.


# Example

See [_example](https://github.com/michimani/invocation-history-extension/tree/main/_example) for using this extension at the Lambda Function using container image.

# License

[MIT](https://github.com/michimani/aws-lambda-api-go/blob/main/LICENSE)

# Author

[michimani210](https://twitter.com/michimani210)