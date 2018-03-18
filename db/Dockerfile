FROM postgres
ENV POSTGRES_DB auth-recipe
ENV POSTGRES_USER admin
ENV POSTGRES_PASSWORD hunter2
COPY init.sql /docker-entrypoint-initdb.d/
