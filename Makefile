.PHONY: install update run all

install:
	go mod tidy

update:
	go get -u all && go mod tidy

EXAMPLES := auth-bearer basic custom-config enums image-upload manually-register-routes
$(EXAMPLES):
	go run examples/$@/main.go