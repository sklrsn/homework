# Homework aka Data Preprocessor

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
