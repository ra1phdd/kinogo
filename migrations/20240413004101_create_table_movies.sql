CREATE TABLE IF NOT EXISTS public.movies
(
    id integer NOT NULL DEFAULT nextval('movies_id_seq'::regclass),
    title text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    country text COLLATE pg_catalog."default" NOT NULL,
    releasedate integer NOT NULL,
    timemovie integer NOT NULL,
    scorekp double precision NOT NULL,
    scoreimdb double precision NOT NULL,
    poster text COLLATE pg_catalog."default",
    typemovie text COLLATE pg_catalog."default",
    views integer,
    likes integer,
    dislikes integer,
    CONSTRAINT movies_pkey PRIMARY KEY (id)
)