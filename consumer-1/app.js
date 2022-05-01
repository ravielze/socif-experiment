#!/usr/bin/env node

var amqp = require('amqplib/callback_api');

const amqpURL = process.env.AMQP_SERVER_URL || 'amqp://test:test@localhost:5672/';

function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

const mainProgram = async () => {
    amqp.connect(amqpURL, async function (error, connection) {
        if (error) {
            console.log('Failed to connect RabbitMQ, waiting 5 seconds...');
            await sleep(5000);
            await mainProgram();
            return;
        }

        connection.createChannel(function (err, channel) {
            if (err) {
                console.log(err);
                return;
            }
            const QUEUE_NAME = process.env.QUEUE_NAME || 'queue-1';

            console.log(' [*] Waiting for messages in %s. To exit press CTRL+C', QUEUE_NAME);
            channel.consume(
                QUEUE_NAME,
                function (msg) {
                    console.log(' [x] Received %s', msg.content.toString());
                },
                {
                    noAck: true,
                }
            );
        });
    });
};
mainProgram();
