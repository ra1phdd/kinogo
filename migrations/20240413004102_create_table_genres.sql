CREATE TABLE IF NOT EXISTS public.genres
(
    id integer NOT NULL DEFAULT nextval('genres_id_seq'::regclass),
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT genres_pkey PRIMARY KEY (id)
)