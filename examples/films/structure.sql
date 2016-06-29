BEGIN;

CREATE TABLE festival (
  name text NOT NULL PRIMARY KEY
);

CREATE TABLE competition (
  id serial PRIMARY KEY,
  name text NOT NULL,
  festival text NOT NULL REFERENCES festival (name) ON UPDATE CASCADE ON DELETE CASCADE,
  year date NOT NULL
);

CREATE TABLE director (
  name text NOT NULL PRIMARY KEY
);

CREATE TABLE film (
  id serial PRIMARY KEY,
  title text NOT NULL,
  year date NOT NULL,
  director text REFERENCES director (name) ON UPDATE CASCADE ON DELETE CASCADE,
  rating real NOT NULL DEFAULT 0,
  language text NOT NULL
);

CREATE TABLE film_nomination (
  id serial PRIMARY KEY,
  competition integer NOT NULL REFERENCES competition (id) ON UPDATE NO ACTION ON DELETE NO ACTION,
  film integer NOT NULL REFERENCES film (id) ON UPDATE CASCADE ON DELETE CASCADE,
  won boolean NOT NULL DEFAULT true
);

COMMIT;