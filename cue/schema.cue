package terraform

data: archive_file:
	[string]: {
		output_path: string
		source_file: string
		type:        "zip"
	}
data: aws_iam_policy_document:
	[string]:
		statement: {
			effect: *"Allow" | "Deny"
			principals?: {
				type: "Service"
				identifiers: [string]
			}
			actions: [string]
		}
resource: aws_iam_role:
	[string]: {
		name:               string
		assume_role_policy: "${data.aws_iam_policy_document.assume_role.json}"
	}
resource: aws_iam_role_policy_attachment:
	[string]: {
		role:       "${aws_iam_role.lambda_role.name}"
		policy_arn: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
	}
resource: aws_lambda_function:
	[string]: {
		filename:         string
		function_name:    string
		handler:          string
		role:             "${aws_iam_role.lambda_role.arn}"
		runtime:          "go1.x"
		source_code_hash: string
	}
resource: aws_cloudwatch_log_group:
	[string]: {
		name:              string
		retention_in_days: 30 & uint
	}
resource: aws_lambda_function_url:
	[string]: {
		authorization_type: string
		function_name:      string
	}
