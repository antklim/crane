docker build -t go-crane -f ./deployments/Dockerfile .

docker run -d -v ~/.aws-lambda-rie:/aws-lambda --entrypoint /aws-lambda/aws-lambda-rie -p 9000:8080 go-crane:latest /aws_lambda

curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"foo": "bar123"}'

aws ecr get-login-password --region AWS_REGION --profile PROFILE | docker login --username AWS --password-stdin AWS_ACC.dkr.ecr.AWS_REGION.amazonaws.com

docker tag go-crane:latest AWS_ACC.dkr.ecr.AWS_REGION.amazonaws.com/go-crane:latest
docker push AWS_ACC.dkr.ecr.AWS_REGION.amazonaws.com/go-crane:latest