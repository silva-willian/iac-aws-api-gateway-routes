# iac-aws-api-gateway-routes

Project created to create the routes of an API in the AWS API Gateway through Swagger JSON

## Building and running

### Local Build

    export AWS_ACCESS_KEY_ID=your_access_key
    export AWS_SECRET_ACCESS_KEY=your_scret_access_key
    export AWS_REGION=us-east-1
    export API_HOST=http://your_host.us-east-1.elb.amazonaws.com
    export API_SWAGGER_JSON_PATH=swagger/v1/swagger.json
    export GATEWAY_NAME=your_gateway_name
    export GATEWAY_STAGE=v1
    export GATEWAY_BASE_PATH=your_base_path
    export GATEWAY_AUTHORIZER_ENABLE=true
    export GATEWAY_AUTHORIZER_NAME=lambda-authorizer

    go run main.go

Keep an eye on your console output

### Building with Docker

    docker build --tag iac-aws-gateway-routes -f devops/application/build/Dockerfile .

    docker run \
        -e AWS_ACCESS_KEY_ID="your_access_key" \
        -e AWS_SECRET_ACCESS_KEY="your_scret_access_key/yeiobKNDp2" \
        -e AWS_REGION="us-east-1" \
        -e API_HOST="http://your_host.us-east-1.elb.amazonaws.com" \
        -e API_SWAGGER_JSON_PATH="swagger/v1/swagger.json" \
        -e GATEWAY_NAME="your_gateway_name" \
        -e GATEWAY_STAGE="v1" \
        -e GATEWAY_BASE_PATH="your_base_path" \
        -e GATEWAY_AUTHORIZER_ENABLE="true" \
        -e GATEWAY_AUTHORIZER_NAME="lambda-authorizer" \
        iac-aws-gateway-routes

## Releases

### 1.0.0 (01.09.2019)

* Initial release