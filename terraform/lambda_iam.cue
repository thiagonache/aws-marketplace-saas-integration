package terraform

configuration: data: {
	aws_iam_policy_document: lambda_assume_role: {
		statement: {
			actions: ["sts:AssumeRole"]
			principals: {
				identifiers: ["lambda.amazonaws.com"]
			}
		}
	}
}
