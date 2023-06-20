package terraform

configuration: data: {
	archive_file: lambda_entitlement: {
		output_path: "/tmp/aws-marketplace-saas-integration/entitlement.zip"
		source_file: "/tmp/aws-marketplace-saas-integration/entitlement"
	}
	aws_iam_policy_document: lambda_entitlement_getentitlements: {
		statement: {
			actions: ["aws-marketplace:GetEntitlements"]
			resources: ["*"]
		}
	}
}
configuration: resource: {
	aws_iam_role: lambda_entitlement_role: {
		assume_role_policy: "${data.aws_iam_policy_document.lambda_assume_role.json}"
	}
	aws_iam_role_policy_attachment: lambda_entitlement: {
		policy_arn: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
		role:       aws_iam_role.lambda_entitlement_role.name
	}
	aws_iam_policy: marketplace_entitlement_get_entitlements: {
		path:        "/"
		description: "Policy to allow lambda to call get entitlements marketplace entitlement API"
		policy:      "${data.aws_iam_policy_document.lambda_entitlement_getentitlements.json}"
	}
	aws_iam_role_policy_attachment: lambda_entitlement_getentitlements: {
		policy_arn: "${\(resource.aws_iam_policy.marketplace_entitlement_get_entitlements.#tfref).arn}"
		role:       resource.aws_iam_role.lambda_entitlement_role.name
	}
	aws_lambda_function: entitlement: {
		environment: {
			variables: {
				"AMSI_SUBSCRIBERS_TABLE_NAME": resource.aws_dynamodb_table.subscribers.name
			}
		}
		filename:         "${data.archive_file.lambda_entitlement.output_path}"
		function_name:    "aws-marketplace-saas-integration-entitlement"
		handler:          "entitlement"
		role:             "${\(resource.aws_iam_role.lambda_entitlement_role.#tfref).arn}"
		source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_entitlement.output_path}\")}"
	}
	aws_lambda_function_url: entitlement: {
		authorization_type: "NONE"
		function_name:      aws_lambda_function.entitlement.function_name
	}
	aws_cloudwatch_log_group: lambda_entitlement: {
		name: "/aws/lambda/\(aws_lambda_function.entitlement.function_name)"
	}
}
