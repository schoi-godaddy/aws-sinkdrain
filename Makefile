ECR_REGISTRY=${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com

VIRTUAL_ENV=.venv

.PHONY: build-image
build-image:
	@cd src/ && \
		GOOS=linux go build -o main && \
		docker build -t ${APP_NAME} . && \
		rm -rf main
	docker tag ${APP_NAME}:latest ${ECR_REGISTRY}/${APP_NAME}:latest

.PHONY: ecr-login
ecr-login:
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ECR_REGISTRY}

.PHONY: deploy-image
deploy-image: ecr-login clean-image
	docker push ${ECR_REGISTRY}/${APP_NAME}:latest

.PHONY: clean-image
clean-image: ecr-login
	aws ecr batch-delete-image \
		--repository-name ${APP_NAME} \
		--image-ids imageTag=latest \
		--region ${AWS_REGION} | jq

.PHONY: update-lambda
update-lambda: deploy-image
	aws lambda update-function-code \
		--function-name ${APP_NAME} \
		--image-uri ${ECR_REGISTRY}/${APP_NAME}:latest \
		--region ${AWS_REGION}

$(VIRTUAL_ENV): requirements.txt
	python3 -m venv .venv
	source .venv/bin/activate && \
		pip install --upgrade pip && \
		pip install -r requirements.txt

.PHONY: launch
launch: $(VIRTUAL_ENV)
	source .venv/bin/activate && sceptre launch -y app

.PHONY: destroy
destroy: $(VIRTUAL_ENV)
	source .venv/bin/activate && sceptre delete -y app
