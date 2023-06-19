package terraform

import "github.com/thiagonache/aws-marketplace-saas-integration/cnifm"

cnifm.#Validate

configuration: data: {
	aws_iam_policy_document: lambda_assume_role: {
		statement: {
			actions: ["sts:AssumeRole"]
			principals: {
				identifiers: ["lambda.amazonaws.com"]
			}
		}
	}
	archive_file: lambda_redirect: {
		output_path: "/tmp/aws-marketplace-saas-integration/redirect.zip"
		source_file: "/tmp/aws-marketplace-saas-integration/redirect"
	}
	archive_file: lambda_landingpage: {
		output_path: "/tmp/aws-marketplace-saas-integration/landingpage.zip"
		source_file: "/tmp/aws-marketplace-saas-integration/landingpage"
	}
}
configuration: output: {
	marketplace_fulfillment_url: {
		description: "Lambda Public Endpoint to be configured on Marketplace Fulfillment URL"
		value:       "${aws_lambda_function_url.landingpage.function_url}"
	}
}
configuration: resource: {
	aws_iam_role: lambda_role: {
		assume_role_policy: "${data.aws_iam_policy_document.lambda_assume_role.json}"
		name:               "lambda_role"
	}
	aws_iam_role_policy_attachment: lambda_basic_role: {
		policy_arn: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
		role:       "${aws_iam_role.lambda_role.name}"
	}
	aws_lambda_function: redirect: {
		environment: {
			variables: {
				"AMSI_REDIRECT_LOCATION": "${aws_lambda_function_url.landingpage.function_url}"
			}
		}
		filename:         "${data.archive_file.lambda_redirect.output_path}"
		function_name:    "aws-marketplace-saas-integration-redirect"
		handler:          "redirect"
		source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_redirect.output_path}\")}"
	}
	aws_lambda_function_url: redirect: {
		authorization_type: "NONE"
		function_name:      "${aws_lambda_function.redirect.function_name}"
	}
	aws_cloudwatch_log_group: lambda_redirect: {
		name: "/aws/lambda/${aws_lambda_function.redirect.function_name}"
	}
	aws_lambda_function: landingpage: {
		environment: {
			variables: {}
		}
		filename:         "${data.archive_file.lambda_landingpage.output_path}"
		function_name:    "aws-marketplace-saas-integration-landingpage"
		handler:          "landingpage"
		source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_landingpage.output_path}\")}"
	}
	aws_lambda_function_url: landingpage: {
		authorization_type: "NONE"
		function_name:      "${aws_lambda_function.landingpage.function_name}"
	}
	aws_cloudwatch_log_group: lambda_landingpage: {
		name: "/aws/lambda/${aws_lambda_function.landingpage.function_name}"
	}
	aws_dynamodb_table: subscribers: {
		attribute: {
			name: "customerIdentifier"
			type: "S"
		}
		billing_mode: "PAY_PER_REQUEST"
		hash_key:     "customerIdentifier"
		name:         "subscribers"
	}
}
configuration: terraform: {
	backend: s3: {
		bucket: "aws-marketplace-saas-go-integration"
		key:    "dev.tfstate"
		region: "us-east-1"
	}
}
