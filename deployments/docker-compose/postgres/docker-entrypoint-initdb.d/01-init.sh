#!/bin/bash

psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" <<-EOSQL
    CREATE USER "qilin" WITH LOGIN PASSWORD 'insecure';
    CREATE DATABASE "qilin" OWNER "qilin";
EOSQL
