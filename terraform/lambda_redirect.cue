package terraform

configuration: data: {
	archive_file: lambda_redirect: {
		output_path: "/tmp/aws-marketplace-saas-integration/redirect.zip"
		source_file: "/tmp/aws-marketplace-saas-integration/redirect"
	}
}
configuration: resource: {
	aws_iam_role: lambda_redirect_role: {
		assume_role_policy: "${data.aws_iam_policy_document.lambda_assume_role.json}"
	}
	aws_iam_role_policy_attachment: lambda_redirect: {
		policy_arn: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
		role:       aws_iam_role.lambda_redirect_role.name
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
		role:             "${\(resource.aws_iam_role.lambda_redirect_role.#tfref).arn}"
		source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_redirect.output_path}\")}"
	}
	aws_lambda_function_url: redirect: {
		authorization_type: "NONE"
		function_name:      aws_lambda_function.redirect.function_name
	}
	aws_cloudwatch_log_group: lambda_redirect: {
		name: "/aws/lambda/\(aws_lambda_function.redirect.function_name)"
	}
}
