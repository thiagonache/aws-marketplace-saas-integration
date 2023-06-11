package terraform

provider: aws: region: "us-east-1"
terraform: backend: s3: bucket: "aws-marketplace-saas-integration-go"
terraform: backend: s3: key:    "dev.tfstate"
terraform: backend: s3: region: "us-east-1"
terraform: required_version: ">=1.0.0"
