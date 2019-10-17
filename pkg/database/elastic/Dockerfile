FROM docker.elastic.co/elasticsearch/elasticsearch:7.4.0

# Add Maintainer Info
LABEL maintainer="Luca Ruggieri <luc4.ruggieri@gmail.com>"

WORKDIR /elasticsearch/server

COPY --chown=elasticsearch:elasticsearch pkg/database/elastic/elasticsearch.yml /usr/share/elasticsearch/config/

ENV ELASTIC_HTTP_PORT=8002

EXPOSE $ELASTIC_HTTP_PORT

CMD ["bin/elasticsearch"]