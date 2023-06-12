package terraform

data: aws_iam_policy_document: assume_role: statement: principals: identifiers: ["lambda.amazonaws.com"]
data: aws_iam_policy_document: assume_role: statement: actions: ["sts:AssumeRole"]

resource: aws_iam_role: lambda_role: assume_role_policy: "${data.aws_iam_policy_document.assume_role.json}"
resource: aws_iam_role: lambda_role: name:               "lambda_role"

resource: aws_iam_role_policy_attachment: lambda_basic_role: policy_arn: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
resource: aws_iam_role_policy_attachment: lambda_basic_role: role:       "${aws_iam_role.lambda_role.name}"
