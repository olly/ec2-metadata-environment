pkg/linux_amd64/update-ec2-metadata-environment: main.go
	mkdir -p $(dir $@)
	GOOS=linux GOARCH=amd64 go build -o $(basename $@) $^

.PHONY: rpm
rpm: pkg/linux_amd64/update-ec2-metadata-environment
	vagrant up --provision
