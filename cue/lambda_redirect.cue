package terraform

data: archive_file:
	lambda_redirect: {
		output_path: "/tmp/aws-marketplace-saas-integration/redirect.zip"
		source_file: "/tmp/aws-marketplace-saas-integration/redirect"
	}
data: aws_iam_policy_document: 
	assume_role:
		statement: {
			principals: {
				identifiers: ["lambda.amazonaws.com"]
			}
			actions: ["sts:AssumeRole"]
		}
resource: aws_iam_role:
	lambda_role: {
		name: "lambda_role"
	}
resource: aws_iam_role_policy_attachment:
	lambda_redirect_basic_role: {}
resource: aws_lambda_function:
	redirect: {
		filename:         "${data.archive_file.lambda_redirect.output_path}"
		function_name:    "aws-marketplace-saas-integration-redirect"
		handler:          "redirect"
		source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_redirect.output_path}\")}"
		environment: {
			variables: {
				"AMSI_REDIRECT_LOCATION": "https://my-app-landing-page.mydomain.com"
			}
		}
	}
resource: aws_cloudwatch_log_group:
	redirect: {
		name: "/aws/lambda/${aws_lambda_function.redirect.function_name}"
	}
resource: aws_lambda_function_url:
	redirect: {
		function_name: "${aws_lambda_function.redirect.function_name}"
	}
