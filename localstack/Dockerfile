FROM localstack/localstack

COPY init-resources.sh /docker-entrypoint-initaws.d/

ENV SERVICES=kinesis,lambda,sqs,cloudformation,s3

ENV KINESIS_ERROR_PROBABILITY=0.1
