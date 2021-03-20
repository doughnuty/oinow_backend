-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.2
-- PostgreSQL version: 12.0
-- Project Site: pgmodeler.io
-- Model Author: ---
-- Database creation must be done outside a multicommand file.
-- These commands were put in this file only as a convenience.
-- -- object: new_database | type: DATABASE --
-- -- DROP DATABASE IF EXISTS new_database;
-- CREATE DATABASE new_database;
-- -- ddl-end --
-- 
-- object: public."Users" | type: TABLE --
-- DROP TABLE IF EXISTS public."Users" CASCADE;
CREATE TABLE IF NOT EXISTS public."users" (
  id serial NOT NULL,
  aituID text,
  name text,
  phone text UNIQUE,
  CONSTRAINT "users_pk" PRIMARY KEY (id)
);
-- ddl-end --
-- ALTER TABLE public."Users" OWNER TO postgres;
-- ddl-end --
-- object: public."Games" | type: TABLE --
-- DROP TABLE IF EXISTS public."Games" CASCADE;
CREATE TABLE IF NOT EXISTS public."games" (
  id serial NOT NULL,
  name text UNIQUE,
  CONSTRAINT "games_pk" PRIMARY KEY (id)
);
-- ddl-end --
-- ALTER TABLE public."Games" OWNER TO postgres;
-- ddl-end --
-- object: public."Score" | type: TABLE --
-- DROP TABLE IF EXISTS public."Score" CASCADE;
CREATE TABLE IF NOT EXISTS "game_scores" (
  id serial NOT NULL,
  userID int,
  gameID int,
  score float,
  CONSTRAINT "game_scores_pk" PRIMARY KEY (id),
  CONSTRAINT "userID_fk" FOREIGN KEY (userID) REFERENCES public."users"(id),
  CONSTRAINT "gameID_fk" FOREIGN KEY (gameID) REFERENCES public."games"(id)
);
-- ddl-end --
-- ALTER TABLE public."Score" OWNER TO postgres;
-- ddl-end --
INSERT INTO games (name) VALUES('shop');