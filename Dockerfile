FROM golang:1.11-stretch AS build

# Копируем исходный код в Docker-контейнер
RUN ls
RUN pwd
RUN echo $GOPATH
# ADD . /opt/build/golang/
ADD . /go/src/DBCourse
#ADD common/ /opt/build/common/

# WORKDIR /opt/build/golang
WORKDIR /go/src/DBCourse

# Собираем и устанавливаем пакет
RUN go get "github.com/lib/pq" && go get "github.com/gin-gonic/gin"
RUN go build

FROM ubuntu:18.04 AS release

# Установка postgresql
ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Собранный ранее сервер
# COPY --from=build /opt/build/golang/* /home/
COPY --from=build /go/src/DBCourse/* /home/

# Back to the root user
USER root

# Объявлем порт сервера
EXPOSE 5000

# Запускаем PostgreSQL и сервер
CMD service postgresql start && export PGPASSWORD='docker' &&  psql docker < home/dump.sql -h localhost -U docker &&\
./home/DBCourse --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://docker:docker@localhost/docker