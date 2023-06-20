package terraform

configuration: {
	resource: {
		aws_iam_role: [Name=string]: {
			name: Name
		}
		aws_sqs_queue: [Name=string]: {
			name: Name
		}
		aws_dynamodb_table: [Name=string]: {
			name: Name
		}
		aws_iam_policy: [Name=string]: {
			name: Name
		}
	}
}
