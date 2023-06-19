.DEFAULT_GOAL:=help
TEMP_DIR=temp
CUE_DIR=terraform

test: # Run unit tests via gotestdox
	gotestdox ./...
build: test # Build Lambda binaries 
	test -d /tmp/aws-marketplace-saas-integration || mkdir /tmp/aws-marketplace-saas-integration
	env GOARCH=amd64 GOOS=linux go build -o /tmp/aws-marketplace-saas-integration/redirect cmd/redirect/main.go
	env GOARCH=amd64 GOOS=linux go build -o /tmp/aws-marketplace-saas-integration/landingpage cmd/landingpage/main.go
generate: build # Generate terraform configuration files via CUE
	test -d $(TEMP_DIR) || mkdir $(TEMP_DIR)
	cue export ./$(CUE_DIR) -e cueniform -f -o $(TEMP_DIR)/config.tf.json
deploy: generate # Run terraform apply
	cd $(TEMP_DIR) && test -d .terraform || terraform init && terraform apply
shutdown: build # Run terraform destroy
	cd $(TEMP_DIR) && terraform destroy
clean: # Delete temp directories
	rm -r $(TEMP_DIR)
help:
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done
