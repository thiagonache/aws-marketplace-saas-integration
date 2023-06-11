package terraform

data: archive_file: lambda_landingpage: output_path: "/tmp/aws-marketplace-saas-integration/landingPage.zip"
data: archive_file: lambda_landingpage: source_file: "/tmp/aws-marketplace-saas-integration/landingPage"

data: aws_iam_policy_document: resolve_customer: statement: actions: ["aws-marketplace:ResolveCustomer"]

output: marketplace_fulfillment_url: description: "Lambda Public Endpoint to be configured on Marketplace Fulfillment URL"
output: marketplace_fulfillment_url: value:       "${aws_lambda_function_url.redirect.function_url}"

resource: aws_lambda_function: landingpage: filename:         "${data.archive_file.lambda_landingpage.output_path}"
resource: aws_lambda_function: landingpage: function_name:    "aws-marketplace-saas-integration-landingpage"
resource: aws_lambda_function: landingpage: handler:          "landingPage"
resource: aws_lambda_function: landingpage: source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_landingpage.output_path}\")}"
resource: aws_lambda_function: landingpage: environment: variables: "AMSI_ENTITLEMENT_QUEUE_URL":  "test"
resource: aws_lambda_function: landingpage: environment: variables: "AMSI_SUBSCRIBERS_TABLE_NAME": "test"

resource: aws_cloudwatch_log_group: landingpage: name: "/aws/lambda/${aws_lambda_function.landingpage.function_name}"

resource: aws_lambda_function_url: landingpage: authorization_type: "NONE"
resource: aws_lambda_function_url: landingpage: function_name:      "${aws_lambda_function.landingpage.function_name}"
