package cnifm

entities: {
	aws_iam_policy_document: #DataSource: {
		statement: {
			actions: [string]
			effect: *"Allow" | "Deny"
			principals: {
				identifiers: [string]
				type: "Service"
			}
		}
	}
	aws_iam_role: #Resource: {
		assume_role_policy: string
		name:               string
	}
	aws_iam_role_policy_attachment: #Resource: {
		policy_arn: string
		role:       string
	}
	archive_file: #DataSource: {
		output_path: string
		source_file: string
		type:        "zip"
	}
	aws_cloudwatch_log_group: #Resource: {
		name:              string
		retention_in_days: uint
		retention_in_days: 30 & >=30
	}
	aws_lambda_function: #Resource: {
		environment: {
			...
		}
		filename:         string
		function_name:    string
		handler:          string
		role:             string
		role:             "${aws_iam_role.lambda_role.arn}"
		runtime:          string
		runtime:          "go1.x"
		source_code_hash: string
	}
	aws_lambda_function_url: #Resource: {
		authorization_type: string
		function_name:      string
	}
}
