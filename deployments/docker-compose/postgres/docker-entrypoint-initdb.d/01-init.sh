#!/bin/bash

psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" <<-EOSQL
    CREATE USER "blueprint" WITH LOGIN PASSWORD 'insecure';
    CREATE DATABASE "blueprint" OWNER "blueprint";
EOSQL
