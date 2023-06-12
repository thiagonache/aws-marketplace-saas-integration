package terraform

data: archive_file: [string]: output_path: string
data: archive_file: [string]: source_file: string
data: archive_file: [string]: type:        string

data: aws_iam_policy_document: [string]: statement: actions: [string]
data: aws_iam_policy_document: [string]: statement: effect: string
data: aws_iam_policy_document: [string]: statement: principals?: identifiers: [string]
data: aws_iam_policy_document: [string]: statement: principals?: type: string

provider: aws: region: string

terraform: backend: s3: region: string

resource: aws_cloudwatch_log_group: [string]: name:              string
resource: aws_cloudwatch_log_group: [string]: retention_in_days: uint

resource: aws_iam_role_policy_attachment: [string]: policy_arn: string
resource: aws_iam_role_policy_attachment: [string]: role:       string

resource: aws_iam_role: [string]: assume_role_policy: string
resource: aws_iam_role: [string]: name:               string

resource: aws_lambda_function_url: [string]: authorization_type: string
resource: aws_lambda_function_url: [string]: function_name:      string

resource: aws_lambda_function: [string]: filename:         string
resource: aws_lambda_function: [string]: function_name:    string
resource: aws_lambda_function: [string]: handler:          string
resource: aws_lambda_function: [string]: role:             string
resource: aws_lambda_function: [string]: runtime:          string
resource: aws_lambda_function: [string]: source_code_hash: string
