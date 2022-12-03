FROM scratch
COPY service /
COPY config/application.yaml /
RUN mkdir -p /messages

ENTRYPOINT [ "/service" ]
CMD [ "serve-http", "--config", "/application.yaml" ]
