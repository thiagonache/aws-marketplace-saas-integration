package cnifm

entities: {
	aws_iam_policy: #Resource: {
		name:        string
		path:        string
		description: string
		policy:      string
	}
	aws_iam_policy_document: #DataSource: {
		statement: {
			actions: [string]
			effect: *"Allow" | "Deny"
			principals?: {
				identifiers: [string]
				type: string
				type: "Service"
			}
			resources?: [string]
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
		type:        string
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
		runtime:          string
		runtime:          "go1.x"
		source_code_hash: string
	}
	aws_lambda_function_url: #Resource: {
		authorization_type: string
		function_name:      string
	}
	aws_dynamodb_table: #Resource: {
		name!:        string // I'm not sure if this should be required
		billing_mode: "PROVISIONED" | "PAY_PER_REQUEST"
		hash_key:     string
		attribute: {
			name: string
			type: "S" | "N" | "B"
		}
	}
	aws_sqs_queue: #Resource: {

	}

}
