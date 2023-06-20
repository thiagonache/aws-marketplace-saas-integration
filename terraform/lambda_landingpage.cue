package terraform

configuration: data: {
	// TODO: replace data.aws_iam_policy_document by JSON directly in the aws_iam_policy
	aws_iam_policy_document: lambda_landingpage_resolvecustomer: {
		statement: {
			actions: ["aws-marketplace:ResolveCustomer"]
			resources: ["*"]
		}
	}
	aws_iam_policy_document: lambda_landingpage_sendmessage: {
		statement: {
			actions: ["sqs:SendMessage"]
			resources: ["${aws_sqs_queue.marketplace_entitlement.arn}"]
		}
	}
	aws_iam_policy_document: lambda_landingpage_putitem: {
		statement: {
			actions: ["dynamodb:PutItem"]
			resources: ["${aws_dynamodb_table.subscribers.arn}"]
		}
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
	aws_iam_role: lambda_landingpage_role: {
		assume_role_policy: "${data.aws_iam_policy_document.lambda_assume_role.json}"
	}
	aws_iam_role_policy_attachment: lambda_landingpage: {
		policy_arn: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
		role:       resource.aws_iam_role.lambda_landingpage_role.name
	}
	aws_iam_policy: marketplace_metering_resolve_customer: {
		path:        "/"
		description: "Policy to allow lambda to call resolve customer marketplace metering API"
		policy:      "${data.aws_iam_policy_document.lambda_landingpage_resolvecustomer.json}"
	}
	aws_iam_role_policy_attachment: lambda_landingpage_resolvecustomer: {
		policy_arn: "${\(resource.aws_iam_policy.marketplace_metering_resolve_customer.#tfref).arn}"
		role:       resource.aws_iam_role.lambda_landingpage_role.name
	}
	aws_iam_policy: sqs_send_message: {
		path:        "/"
		description: "Policy to allow lambda to call send message SQS API"
		policy:      "${data.aws_iam_policy_document.lambda_landingpage_sendmessage.json}"
	}
	aws_iam_role_policy_attachment: lambda_landingpage_sendmessage: {
		policy_arn: "${\(resource.aws_iam_policy.sqs_send_message.#tfref).arn}"
		role:       resource.aws_iam_role.lambda_landingpage_role.name
	}
	aws_iam_policy: dynamodb_put_item: {
		path:        "/"
		description: "Policy to allow lambda to call put item Dynamodb API"
		policy:      "${data.aws_iam_policy_document.lambda_landingpage_sendmessage.json}"
	}
	aws_iam_role_policy_attachment: lambda_landingpage_putitem: {
		policy_arn: "${\(resource.aws_iam_policy.dynamodb_put_item.#tfref).arn}"
		role:       resource.aws_iam_role.lambda_landingpage_role.name
	}
	aws_lambda_function: landingpage: {
		environment: {
			variables: {
				"AMSI_ENTITLEMENT_QUEUE_URL":  "${resource.aws_sqs_queue.marketplace_entitlement.url}"
				"AMSI_SUBSCRIBERS_TABLE_NAME": resource.aws_dynamodb_table.subscribers.name
			}
		}
		filename:         "${data.archive_file.lambda_landingpage.output_path}"
		function_name:    "aws-marketplace-saas-integration-landingpage"
		handler:          "landingpage"
		role:             "${\(resource.aws_iam_role.lambda_landingpage_role.#tfref).arn}"
		source_code_hash: "${filebase64sha256(\"${data.archive_file.lambda_landingpage.output_path}\")}"
	}
	aws_lambda_function_url: landingpage: {
		authorization_type: "NONE"
		function_name:      aws_lambda_function.landingpage.function_name
	}
	aws_cloudwatch_log_group: lambda_landingpage: {
		name: "/aws/lambda/\(resource.aws_lambda_function.landingpage.function_name)"
	}
}
