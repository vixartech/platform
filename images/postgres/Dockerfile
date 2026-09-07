ARG POSTGRES_IMAGE=18.4-alpine3.24

FROM postgres:${POSTGRES_IMAGE}

RUN apk add --no-cache postgresql-pg_cron postgresql-pgvector \
    && cp /usr/lib/postgresql*/pg_cron.so "$(pg_config --pkglibdir)/" \
    && cp /usr/share/postgresql*/extension/pg_cron* "$(pg_config --sharedir)/extension/" \
    && cp /usr/lib/postgresql*/vector.so "$(pg_config --pkglibdir)/" \
    && cp /usr/share/postgresql*/extension/vector* "$(pg_config --sharedir)/extension/"
