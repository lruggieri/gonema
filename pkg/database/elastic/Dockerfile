FROM docker.elastic.co/elasticsearch/elasticsearch:7.4.0

# Add Maintainer Info
LABEL maintainer="Luca Ruggieri <luc4.ruggieri@gmail.com>"

WORKDIR /elasticsearch/server

# required to avoid error "memory areas vm.max_map_count [65530] is too low, increase to at least [262144]"
RUN echo 262144 > /proc/sys/vm/max_map_count

COPY --chown=elasticsearch:elasticsearch pkg/database/elastic/elasticsearch.yml /usr/share/elasticsearch/config/

# Port 8080 required by Google Cloud Run
ENV ELASTIC_HTTP_PORT 8080

EXPOSE $ELASTIC_HTTP_PORT

CMD ["bin/elasticsearch"]