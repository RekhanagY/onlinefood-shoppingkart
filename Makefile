ECR_REPO = 231839105289.dkr.ecr.ap-southeast-2.amazonaws.com/oolio-tech-challenge
REGION = ap-southeast-2
VERSION = 1.0.0

update_ecr_image:
	aws ecr get-login-password --region $(REGION) | docker login --username AWS --password-stdin $(ECR_REPO)
	docker build -t $(ECR_REPO):$(VERSION) .
	docker push $(ECR_REPO):$(VERSION)

.PHONY: update_ecr_image