.DEFAULT_GOAL:=help
test: # Run unit tests via gotestdox
	gotestdox ./...
build: test # Generate Lambda binaries and terraform configuration files via cue
	test -d /tmp/aws-marketplace-saas-integration || mkdir /tmp/aws-marketplace-saas-integration
	env GOARCH=amd64 GOOS=linux go build -o /tmp/aws-marketplace-saas-integration/redirect cmd/redirect/main.go
	env GOARCH=amd64 GOOS=linux go build -o /tmp/aws-marketplace-saas-integration/landingPage cmd/landingpage/main.go
	test -d temp || mkdir temp
	cd temp && test -d .terraform || terraform init
	cue vet -c ./...
	cue export ./... > temp/main.tf.json
deploy: build # Run terraform apply
	cd temp && terraform apply
shutdown: build # Run terraform destroy
	cd temp && terraform destroy
clean: build # Run terraform destroy with auto-approve and delete temp directories
	cd temp && terraform destroy -auto-approve || true
	rm -r temp
help:
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done
