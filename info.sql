-- Database: golang-user-api

-- DROP DATABASE golang-user-api;

CREATE DATABASE user_golang_api
  WITH OWNER = root
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       CONNECTION LIMIT = -1;

-- Table: books

-- DROP TABLE books;

CREATE TABLE books
(
 id serial NOT NULL,
 name character varying NOT NULL,
 author character varying,
 pages integer,
 publication_date date,
 CONSTRAINT pk_books PRIMARY KEY (id )
)
WITH (
 OIDS=FALSE
);
ALTER TABLE books
 OWNER TO root;