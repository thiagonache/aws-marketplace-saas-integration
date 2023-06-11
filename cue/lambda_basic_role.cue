package terraform

data: aws_iam_policy_document: assume_role: statement: principals: identifiers: ["lambda.amazonaws.com"]
data: aws_iam_policy_document: assume_role: statement: actions: ["sts:AssumeRole"]
resource: aws_iam_role: lambda_role: name: "lambda_role"
resource: aws_iam_role_policy_attachment: lambda_basic_role: {}
