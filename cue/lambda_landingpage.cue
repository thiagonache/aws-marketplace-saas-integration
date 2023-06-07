package terraform

data: archive_file:
	lambda_landingpage: {
		output_path: "/tmp/aws-marketplace-saas-integration/landingPage.zip"
		source_file: "/tmp/aws-marketplace-saas-integration/landingPage"
	}

data: aws_iam_policy_document:
	resolve_customer:
		statement: {
			actions: ["aws-marketplace:ResolveCustomer"]
		}

output: 
	marketplace_fulfillment_url: {
		description: "Lambda Public Endpoint to be configured on Marketplace Fulfillment URL"
		value: "${aws_lambda_function_url.redirect.function_url}"
	}

resource: aws_lambda_function:
	landingpage: {
		filename:         "${data.archive_file.lambda_landingpage.output_path}"
		function_name:    "aws-marketplace-saas-integration-landingpage"
		handler:          "landingPage"
		source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_landingpage.output_path}\")}"
		environment: {
			variables: {
				"AMSI_ENTITLEMENT_QUEUE_URL" : "test"
				"AMSI_SUBSCRIBERS_TABLE_NAME" : "test"
			}
		}
	}

resource: aws_cloudwatch_log_group:
	landingpage: {
		name: "/aws/lambda/${aws_lambda_function.landingpage.function_name}"
	}

resource: aws_lambda_function_url:
	landingpage: {
		authorization_type: "NONE"
		function_name: "${aws_lambda_function.landingpage.function_name}"
	}
