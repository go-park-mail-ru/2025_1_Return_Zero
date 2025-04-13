FROM postgres:16

RUN apt-get update && apt-get install -y curl

RUN apt-get -y install postgresql-16-cron

RUN echo "shared_preload_libraries = 'pg_cron'" >> /usr/share/postgresql/postgresql.conf.sample

RUN echo "cron.database_name = 'returnzero'" >> /usr/share/postgresql/postgresql.conf.sample


