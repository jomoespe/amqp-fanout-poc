# AMQP Multiple Unique Consumer

A simple example of producing a message to a RabbitMQ exchange (*poc.messages*) and a consumer will consume the message. The consumer will create and bind a queue on startup time, and the queue will be removed when connection is closed. We can create as many consumers as we want.

```text
                             +----------------+                +----------+
 +---------+                 |    RabbitMQ    |                | Consumer |-+
 | produce |--( message )--> |    exchange    |==( message )=> +----------+ |-+
 +---------+                 | (poc.messages) |                  +----------+ |
                             +----------------+                    +-----------
```

## Requisites

- **Docker** to start a RabbitMQ instance
- **Go 1.7+** to compile the examples


## Building

```bash
make [clean] [[bin/produce] [bin/produce] | all]
```

Example: clean and build all: `make clean all`

## Running

### Start RabbitMQ

```bash
docker run -detach --rm \
    --hostname poc-fanout-notification-service-rabbit \
    --name poc-fanout-notification-service-rabbit \
    --publish 15672:15672 \
    --publish 5672:5672 \
    rabbitmq:3.8-management
```

Then you can [access RabbitMQ management console](http://localhost:15672), with user=`guest`, password=`guest`, and create an exchange with properties:

| name        |  value           |
|-------------|------------------|
| Name        | **poc.exchange** |
| Type        | **fanout**       |
| Durability  | Durable          |
| Auto delete | No               |
| Internal    | No               |
| Arguments   | *none*           |

### Start consumers and producers

With RabbitMQ up & running, and exchange created:

Start a consumer: 

```bash
$ ./bin/consume
2022/03/16 09:56:06 Queue name: amq.gen-tPn-i_2VtGRkFHbBcRLosA
2022/03/16 09:56:06  [*] Waiting for messages. To exit press CTRL+C
```

> The queue name will be different in every execution

Produce a message:

```bash
./bin/produce MSG
```
