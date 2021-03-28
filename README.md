# SQS and Kinesis

 Parse SQS events and push it to Kinesis stream
## Solution Environment
```
* Golang (1.16)
* CMake
* docker
* docker-compose
```

docker-compose.yml
```
  ....................................
  ....................................
  ....................................
  
  preprocessor:
    container_name: "preprocessor"
    build: preprocessor
    network_mode: "host"
    environment:
      access_key_id: AKIA2RVR24VPMP788U3S
      secret_access_key: mDOWS+uN7dogVkHDaTuHaoyQ29Ju7pJmsvrrug8o
      region: eu-west-1
      end_point: http://localhost:4566
      queue_url: http://localhost:4566/000000000000/submissions
      stream_name: events
      sqs_messages_batch: 10
      sqs_poll_interval: 10
      mode: standalone

  ....................................
  ....................................
  ....................................
```
## Run the Application

Please follow below steps to run the application in your local environment.
```
* tar xvf homework.zip
* cd homework/
* make clean
* make
```

## Messages

Incoming messages from SQS (valid message as described in requirements)

```
{
    "submission_id": "8be01974-0f15-4d4f-809f-257a6de4d3f9",
    "device_id": "31980191-fa99-4768-8a7b-0e81397ca6ef",
    "time_created": "2021-03-28T15:44:18.592840",
    "events": {
        "new_process": [{
            "cmdl": "notepad.exe",
            "user": "john"
        }, {
            "cmdl": "whoami",
            "user": "admin"
        }, {
            "cmdl": "calculator.exe",
            "user": "john"
        }],
        "network_connection": [{
            "source_ip": "192.168.0.2",
            "destination_ip": "23.13.252.39",
            "destination_port": 43696
        }, {
            "source_ip": "192.168.0.1",
            "destination_ip": "142.250.74.110",
            "destination_port": 46916
        }, {
            "source_ip": "192.168.0.2",
            "destination_ip": "23.13.252.39",
            "destination_port": 58976
        }, {
            "source_ip": "192.168.0.2",
            "destination_ip": "23.13.252.39",
            "destination_port": 5817
        }, {
            "source_ip": "192.168.0.1",
            "destination_ip": "142.250.74.110",
            "destination_port": 28512
        }]
    }
}
```

Outgoing message to Kinesis

```
{
    "id": "770ca48d-8fdc-11eb-bbef-68f72870fc5b",
    "device_id": "31980191-fa99-4768-8a7b-0e81397ca6ef",
    "new_process": [{
        "cmdl": "notepad.exe",
        "user": "john"
    }, {
        "cmdl": "whoami",
        "user": "admin"
    }, {
        "cmdl": "calculator.exe",
        "user": "john"
    }],
    "network_connection": [{
        "source_ip": "192.168.0.2",
        "destination_ip": "23.13.252.39",
        "destination_port": 43696
    }, {
        "source_ip": "192.168.0.1",
        "destination_ip": "142.250.74.110",
        "destination_port": 46916
    }, {
        "source_ip": "192.168.0.2",
        "destination_ip": "23.13.252.39",
        "destination_port": 58976
    }, {
        "source_ip": "192.168.0.2",
        "destination_ip": "23.13.252.39",
        "destination_port": 5817
    }, {
        "source_ip": "192.168.0.1",
        "destination_ip": "142.250.74.110",
        "destination_port": 28512
    }],
    "created": "2021-03-28T15:44:21.782645Z"
}
```

Kinesis Output

```
{
    "SequenceNumber": "49616803008866257758247545874384882562299446046301880338",
    "ShardId": "shardId-000000000001"
}
```

There is another solution where I have tried to deploy a lambda function to localstack.
Lambda deployment and creation of event-source-mapping succeeds. But,
for some reason, I couldn't get event-source-mapping working in localstack.
```
cd homework

make clean

make stack

make serverless
```

```
deploying lambda to localstack ...
aws lambda create-function --function-name preprocessor --runtime go1.x \
--zip-file fileb://preprocessor/dist/linux/amd64/preprocessor.zip \
--handler preprocessor --endpoint-url=http://localhost:4566 \
--role arn:aws:iam::kalai:role/execution_role
{
    "FunctionName": "preprocessor",
    "FunctionArn": "arn:aws:lambda:us-east-1:000000000000:function:preprocessor",
    "Runtime": "go1.x",
    "Role": "arn:aws:iam::kalai:role/execution_role",
    "Handler": "preprocessor",
    "CodeSize": 3403253,
    "Description": "",
    "Timeout": 3,
    "LastModified": "2021-03-28T21:23:58.729+0000",
    "CodeSha256": "vyJBbvM9+bQxqzshXGp6zBqwryZDCtXKhRRQoMJ8+0M=",
    "Version": "$LATEST",
    "VpcConfig": {},
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "RevisionId": "d7d2d7d4-3b91-43a7-87a2-1ab6d2489b99",
    "State": "Active",
    "LastUpdateStatus": "Successful",
    "PackageType": "Zip"
}
enable trigger to launch lambdas when message published to SQS...
aws --endpoint-url=http://localhost:4566 lambda create-event-source-mapping \
--event-source-arn arn:aws:sqs:eu-west-1:000000000000:submissions \
--function-name preprocessor
{
    "UUID": "7809a537-bffb-4c84-b98a-4674e443fdd0",
    "StartingPosition": "LATEST",
    "BatchSize": 10,
    "EventSourceArn": "arn:aws:sqs:eu-west-1:000000000000:submissions",
    "FunctionArn": "arn:aws:lambda:us-east-1:000000000000:function:preprocessor",
    "LastModified": 1616966639.0,
    "LastProcessingResult": "OK",
    "State": "Enabled",
```

