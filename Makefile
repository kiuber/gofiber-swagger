.PHONY: install update

install:
	go mod tidy

update:
	go get -u all && go mod tidy && gofmt -w -l .

EXAMPLES := auth-bearer basic custom-config enums image-upload manually-register-routes
$(EXAMPLES):
	go run examples/$@/main.go