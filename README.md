# AMQP Multiple Unique Consumer

A simple example of producing a message to a RabbitMQ exchange (*poc.messages*) and a consumer will consume the message. The consumer will create and bind a queue on startup time, and the queue will be removed when connection is closed. We can create as many consumers as we want.

Next diagram shows architecture of this PoC.

```text
                             +----------------+                +-----------+              +------------+
+-------------+              |     fanout     |            +-->|   queue   |--(message)-->|  Consumer  |
|  Publisher  |--(message)-->|    exchange    |--[ bind ]--|   +-----------+-+            +------------+-+
+-------------+              | (poc.messages) |            +-->  |   queue   |--(message)-->|  Consumer  |
                             +----------------+                  +-----------+               +------------+
                                     |                                                            ^
                            (unconsumed messages)                                                 |
                                     |                                               (read messages on connect)
                                     V                                                            |
                            +------------------- +
                            |  alternate fanout  |             +-----------+                      |
                            |      exchange      |--[ bind ]-->|   queue   |- - - - - - - - - - - +
                            | (poc.messages.alt) |             +-----------+
                            +--------------------+
```

To avoid message lost, we use an **alternate exchange** to collect unroutable messages, mainly because no queues consuming from main exchange. When no consumer

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
2022/03/17 09:56:05  [*] Waiting for alternate messages.
2022/03/16 09:56:06 Queue name: c4e7026f-a0f7-4de2-9bdc-8be0f73367e4@poc.messages bound to poc.messages
2022/03/16 09:56:06  [*] Waiting for messages. To exit press CTRL+C
```

> The queue name will be different in every execution

Produce a message:

```bash
./bin/produce MSG
```

## Related

- [RabbitMQ Exchanges, routing keys and bindings](https://www.cloudamqp.com/blog/part4-rabbitmq-for-beginners-exchanges-routing-keys-bindings.html) for a description of RabbitMQ exchange types.
- [Collecting Unroutable Messages in a RabbitMQ Alternate Exchange](https://www.cloudamqp.com/blog/collecting-unroutable-messages-in-a-rabbitmq-alternate-exchange.html)
- [RabbitMQ Documentation - Alternate Exchanges](https://www.rabbitmq.com/ae.html)
- [RabbitMQ Go Tutorial](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
