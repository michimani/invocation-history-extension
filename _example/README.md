Example of Invocation History Extension
===

# Preparing

1. Login to ECR

    ```bash
    REGION='ap-northeast-1'
    AWS_ACCOUNT_ID=$(
      aws sts get-caller-identity \
      --query 'Account' \
      --output text) \
    && aws ecr get-login-password \
      --region "${REGION}" \
      | docker login \
      --username AWS \
      --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com
    ```

2. Build and zip extension

    ```bash
    make build-ex
    ```

3. Build function image

    ```bash
    make build-func
    ```
    
# Invoke function

1. Run Lambda

    ```bash
    make run
    ```

2. Invoke function

    ```bash
    curl \
    -H 'Content-Type: application/json' \
    http://localhost:9000/2015-03-31/functions/function/invocations
    ```
    
    Each time a Lambda function is invoked, the request IDs (Invocation History) of the previously invoked functions will increase. For example, the third execution will result in the following response.
    
    ```json
    {
      "message": "Current request ID is a6fd06be-6e87-4551-b65f-ff5eda49484f",
      "invocations": [
        {
          "awsRequestId": "fea355f0-b7ec-4d5f-b6b5-b6049bb317f3",
          "invocatedAt": "2022-12-04T00:17:22.31477492Z"
        },
        {
          "awsRequestId": "39a06690-fb49-444d-b8dc-11c2f7792317",
          "invocatedAt": "2022-12-04T00:17:22.922697963Z"
        },
        {
          "awsRequestId": "a6fd06be-6e87-4551-b65f-ff5eda49484f",
          "invocatedAt": "2022-12-04T00:17:23.771295398Z"
        }
      ]
    }
    ```

# Author

[michimani210](https://twitter.com/michimani210)