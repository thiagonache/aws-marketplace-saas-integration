package terraform

configuration: resource: {
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
