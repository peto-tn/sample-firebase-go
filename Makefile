.PHONY: build gomodgen deploy

####################################
# contant
####################################
FUNCTION    := sample-go
ENTRY_POINT := OnCall
PROJECT_ID  := sample-go
REGION      := asia-northeast1

####################################
# task
####################################
build:
	go build .

gomodgen:
	GO111MODULE=on go mod init ${FUNCTION}

deploy:
	gcloud functions deploy $(FUNCTION) --project $(PROJECT_ID) --region $(REGION) --set-env-vars PROJECT_ID=$(PROJECT_ID) --entry-point $(ENTRY_POINT) --runtime go111 --trigger-http
