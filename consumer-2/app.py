import pika
from pika.exceptions import AMQPConnectionError
import os
import sys
import time
import logging

logging.basicConfig(level=logging.DEBUG)


def main():
    amqpURL = os.getenv('AMQP_SERVER_URL')
    if not(amqpURL):
        amqpURL = 'amqp://test:test@localhost:5672/'
    queueName = os.getenv('QUEUE_NAME')
    if not(queueName):
        queueName = 'queue-2'

    isConnected = False
    while not(isConnected):
        try:
            connection = pika.BlockingConnection(pika.URLParameters(amqpURL))
            isConnected = True
        except AMQPConnectionError:
            isConnected = False
            logging.info("Failed to connect RabbitMQ, waiting 5 seconds...")
            time.sleep(5)

    channel = connection.channel()
    channel.queue_declare(queue=queueName, durable=True)

    def callback(ch, method, properties, body):
        logging.info(" [x] Received %s" % body.decode('utf-8'))

    channel.basic_consume(queue=queueName, auto_ack=True,
                          on_message_callback=callback)

    logging.info(' [*] Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        logging.info('Interrupted')
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)
