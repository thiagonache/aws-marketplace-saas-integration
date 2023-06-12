package terraform

data: archive_file: lambda_redirect: output_path: "/tmp/aws-marketplace-saas-integration/redirect.zip"
data: archive_file: lambda_redirect: source_file: "/tmp/aws-marketplace-saas-integration/redirect"

resource: aws_cloudwatch_log_group: redirect: name: "/aws/lambda/${aws_lambda_function.redirect.function_name}"

resource: aws_lambda_function_url: redirect: authorization_type: "NONE"
resource: aws_lambda_function_url: redirect: function_name:      "${aws_lambda_function.redirect.function_name}"

resource: aws_lambda_function: redirect: environment: variables: "AMSI_REDIRECT_LOCATION": "${aws_lambda_function_url.landingpage.function_url}"
resource: aws_lambda_function: redirect: filename:         "${data.archive_file.lambda_redirect.output_path}"
resource: aws_lambda_function: redirect: function_name:    "aws-marketplace-saas-integration-redirect"
resource: aws_lambda_function: redirect: handler:          "redirect"
resource: aws_lambda_function: redirect: source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_redirect.output_path}\")}"
