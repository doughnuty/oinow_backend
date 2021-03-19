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
CREATE TABLE public."Users" (
  id serial NOT NULL,
  phone text,
  name text,
  CONSTRAINT "Users_pk" PRIMARY KEY (id)
);
-- ddl-end --
-- ALTER TABLE public."Users" OWNER TO postgres;
-- ddl-end --
-- object: public."Games" | type: TABLE --
-- DROP TABLE IF EXISTS public."Games" CASCADE;
CREATE TABLE public."Games" (
  id serial NOT NULL,
  name text,
  CONSTRAINT "Games_pk" PRIMARY KEY (id)
);
-- ddl-end --
-- ALTER TABLE public."Games" OWNER TO postgres;
-- ddl-end --
-- object: public."Score" | type: TABLE --
-- DROP TABLE IF EXISTS public."Score" CASCADE;
CREATE TABLE public."Score" (
  id serial NOT NULL,
  "id_User" int,
  "id_Game" int,
  score float,
  CONSTRAINT "Score_pk" PRIMARY KEY (id),
  CONSTRAINT "id_User_fk" FOREIGN KEY ("id_User") REFERENCES public."Users"(id),
  CONSTRAINT "id_Game_fk" FOREIGN KEY ("id_Game") REFERENCES public."Games"(id)
);
-- ddl-end --
-- ALTER TABLE public."Score" OWNER TO postgres;
-- ddl-end --