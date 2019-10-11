#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER "qilin-hasura" WITH PASSWORD 'insecure';
    CREATE DATABASE "qilin-hasura" OWNER "$POSTGRES_USER";
EOSQL
