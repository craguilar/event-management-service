# Event Management

This sample code helps get you started with a simple Go web application deployed by AWS CloudFormation to AWS Lambda and 
Amazon API Gateway.

## What's Here?

This sample includes:

* README.md - this file
* buildspec.yml - this file is used by AWS CodeBuild to package your
  application for deployment to AWS Lambda
* cmd - This directory contains go code that glues together our service implementation in internal , the code in this directory refers to execution especific implementation details of the service.
* internal - This directory contains go code  which is NON public and contains implementation details we don't want to expose as a public API.
* scripts - Yeah , these are scripts.
* template.yml - this file contains the AWS Serverless Application Model (AWS SAM) used
  by AWS CloudFormation to deploy your application to AWS Lambda and Amazon API
  Gateway.
* template-configuration.json - this file contains the project ARN with placeholders used for tagging resources with the project ID  

## Development

### AWS Lambda execution

To work on the sample code, you'll need to clone your project's repository to your
local computer. If you haven't, do that first , then:

1. Install Go.  See https://golang.org/dl/ for details.

1. Install your dependencies:

    ```bash
    go mod init
    ```

    or if already installed

    ```bash
    go mod tidy
    ```

1. Install the SAM CLI. For details see https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html

1. Run the following command in your repository to build the main.go file.

    ```bash
    GOARCH=amd64 GOOS=linux go build -o main cmd/http/lambda/*.go
    ```

1. Start the development server:

    ```bash
    sam local start-api -p 8080
    ```

### Server mode

Run server mode using mock:

```bash
go run cmd/http/server/*.go
```

Then open http://127.0.0.1:8080/ in a web browser to view your webapp or execute

  ```bash
  scripts/test-integration.sh
  ```

### Dynamo DB

Start dynamo db locally.

1. Download Dynamo DB from  https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html

2. Create the `events` table.

```bash
aws dynamodb create-table --endpoint-url http://localhost:8000 \
	--table-name events --billing-mode PAY_PER_REQUEST \
	--attribute-definitions \
	AttributeName=id,AttributeType=S \
	AttributeName=entityType,AttributeType=S \
	--key-schema \
	AttributeName=id,KeyType=HASH \
	AttributeName=entityType,KeyType=RANGE \
	--global-secondary-indexes \
	"[{\"IndexName\": \"ownerIdx\",\"KeySchema\":[{\"AttributeName\":\"entityType\",\"KeyType\":\"HASH\"}], \
        \"Projection\":{\"ProjectionType\":\"ALL\"}}]"
```

3. Check the table exists

```bash
aws dynamodb describe-table --table-name events --endpoint-url http://localhost:800
```

And validate table estructure

```bash
aws dynamodb scan --table-name events --endpoint-url http://localhost:8000
```

Other useful commands:

Scan records

```bash
aws dynamodb scan --table-name events --endpoint-url http://localhost:8000
```

## Deployment

Validate your Cloudformation template using below command:

```bash
aws cloudformation validate-template --template-body file://./template.yml
```

## API Gateway

### Authorizer

Once deployed in AWS the Lambda is securde to ONLY be callable from IAM or trusted AWS Service , that means unless you have a trusted account in AWS you
won't be able to call the Lambda.

To enable non AWS IAM based Authentication/Authorization our Lambda integates with API Gateway where we have attached the JWT Authorizer that will require a valid JWT token (in this case provided by Amazon Cognito) and pass this token to the AWS Gateway path as a query parameter named `token` , also I added a path on OPTIONS to not require Auth, see below:

OPTIONS /{proxy+} no op see https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-routes.html?icmpid=apigateway_console_help

## Contributing

### Format

```bash
gofmt -w -s .
```

### Vulnerability checking

Requires Go version 1.18 - see https://go.dev/blog/vuln

## ToDo

1. Sending messages https://www.twilio.com/docs/whatsapp
1. Harden potential abuse of parameters to introduce max size restrictions.
1. Implement pagination , see scan Limit and ExclusiveStartKey - this looks more like a Cursor basde pagination.

## References

1. Directory structure mainly based on https://www.gobeyond.dev/packages-as-layers/ , https://www.gobeyond.dev/standard-package-layout/ and  https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091 . Other useful links:   : https://leonardqmarcq.com/posts/go-project-structure-for-api-gateway-lambda-with-aws-sam 
1. Best practices for working with AWS Lambda functions - https://docs.aws.amazon.com/lambda/latest/dg/best-practices.html
1. AWS Lamdba Golang https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
1. AWS Lambda EnvVars https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html

## Pending to document:

1. https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-the-id-token.html
1. https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/using-dynamodb-with-go-sdk.html
1. https://github.com/go-playground/validator
1. Using Viper as configuration framework https://github.com/spf13/viper/blob/master/viper.go