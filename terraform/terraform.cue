package terraform

import "github.com/thiagonache/aws-marketplace-saas-integration/cnifm"

cnifm.#Validate

configuration: terraform: {
	backend: s3: {
		bucket: "aws-marketplace-saas-go-integration"
		key:    "dev.tfstate"
		region: "us-east-1"
	}
}
