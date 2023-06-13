package terraform

#ArchiveType:                    "zip"
#AWSRegion:                      *"us-east-1" | "af-south-1" | "ap-east-1" | "ap-northeast-1" | "ap-northeast-2" | "ap-northeast-3" | "ap-south-1" | "ap-south-2" | "ap-southeast-1" | "ap-southeast-2" | "ap-southeast-3" | "ap-southeast-4" | "ca-central-1" | "eu-central-1" | "eu-central-2" | "eu-north-1" | "eu-south-1" | "eu-south-2" | "eu-west-1" | "eu-west-2" | "eu-west-3" | "me-central-1" | "me-south-1" | "sa-east-1" | "us-east-2" | "us-west-1" | "us-west-2"
#CloudwatchLogGroupMinRetention: 30
#StatementEffect:                *"Allow" | "Deny"

data: archive_file: [string]: output_path: string
data: archive_file: [string]: source_file: string
data: archive_file: [string]: type:        string
data: archive_file: [string]: type:        #ArchiveType

data: aws_iam_policy_document: [string]: statement: actions: [string]
data: aws_iam_policy_document: [string]: statement: effect: string
data: aws_iam_policy_document: [string]: statement: effect: #StatementEffect
data: aws_iam_policy_document: [string]: statement: principals?: identifiers: [string]
data: aws_iam_policy_document: [string]: statement: principals?: type: string
data: aws_iam_policy_document: [string]: statement: principals?: type: "Service"

provider: aws: region: string
provider: aws: region: #AWSRegion

terraform: backend: s3: region: string
terraform: backend: s3: region: #AWSRegion

resource: aws_cloudwatch_log_group: [string]: name:              string
resource: aws_cloudwatch_log_group: [string]: retention_in_days: uint
resource: aws_cloudwatch_log_group: [string]: retention_in_days: #CloudwatchLogGroupMinRetention & >=#CloudwatchLogGroupMinRetention

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
resource: aws_lambda_function: [string]: role:             "${aws_iam_role.lambda_role.arn}"
resource: aws_lambda_function: [string]: runtime:          string
resource: aws_lambda_function: [string]: runtime:          "go1.x"
resource: aws_lambda_function: [string]: source_code_hash: string
