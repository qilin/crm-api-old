#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER "qilin" WITH LOGIN PASSWORD 'insecure';
    CREATE DATABASE "qilin" OWNER "qilin";
    CREATE USER "natss" WITH LOGIN PASSWORD 'insecure';
    CREATE DATABASE "natss" OWNER "natss";
EOSQL
