FROM openjdk:8-jdk-stretch

ARG SSH_KEY_EMAIL

COPY ./bin/openapi-generator-cli.jar ./
COPY ./entrypoint.sh ./
RUN chmod +x ./entrypoint.sh ./openapi-generator-cli.jar && \
    ssh-keygen -t rsa -N "" -C "${SSH_KEY_EMAIL}" -b 4096 -f ~/.ssh/id_rsa -q && \
    echo "SSH-Public-Key:\n\n$(cat ~/.ssh/id_rsa.pub)\n\n"

ENTRYPOINT [ "./entrypoint.sh" ]