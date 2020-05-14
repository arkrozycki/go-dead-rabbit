FROM rabbitmq:3.8-management

ADD rabbitmq.conf /etc/rabbitmq/
ADD definitions.json /etc/rabbitmq/
COPY health-check /usr/local/bin/

RUN chown rabbitmq:rabbitmq /etc/rabbitmq/rabbitmq.conf /etc/rabbitmq/definitions.json

HEALTHCHECK CMD ["health-check"]