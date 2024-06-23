create table if not exists movies
(
    id          serial
        primary key,
    title       text             not null,
    description text,
    releasedate integer          not null,
    timemovie   integer          not null,
    scorekp     double precision not null,
    scoreimdb   double precision not null,
    poster      text,
    typemovie   integer
);

alter table movies
    owner to postgres;

create index if not exists trgm_idx
    on movies using gin (title gin_trgm_ops);

create table if not exists genres
(
    id   serial
        primary key,
    name text not null
);

alter table genres
    owner to postgres;

create table if not exists "moviesGenres"
(
    idmovie integer not null,
    idgenre integer not null
);

alter table "moviesGenres"
    owner to postgres;

create table if not exists users
(
    id         integer not null
        primary key,
    username   text,
    photourl   text,
    first_name text,
    last_name  text,
    auth_date  integer
);

alter table users
    owner to postgres;

create table if not exists "moviesLikes"
(
    "userId"    integer not null
        constraint "likes_userId_fkey"
            references users,
    "movieId"   integer not null
        constraint "likes_movieId_fkey"
            references movies,
    "likeType"  varchar(10)
        constraint "likes_likeType_check"
            check (("likeType")::text = ANY
                   ((ARRAY ['like'::character varying, 'dislike'::character varying])::text[])),
    "createdAt" timestamp default CURRENT_TIMESTAMP,
    constraint likes_pkey
        primary key ("userId", "movieId")
);

alter table "moviesLikes"
    owner to postgres;

create table if not exists "moviesViews"
(
    "userId"   integer not null
        constraint "views_userId_fkey"
            references users,
    "movieId"  integer not null
        constraint "views_movieId_fkey"
            references movies,
    "viewDate" date    not null,
    constraint views_pkey
        primary key ("userId", "movieId", "viewDate")
);

alter table "moviesViews"
    owner to postgres;

create table if not exists countries
(
    id   serial
        primary key,
    name text not null
        unique
);

alter table countries
    owner to postgres;

create table if not exists "moviesCountries"
(
    idmovie   integer not null,
    idcountry integer not null
);

alter table "moviesCountries"
    owner to postgres;

create table if not exists comments
(
    id          serial
        primary key,
    "userId"    integer
        references users,
    "movieId"   integer
        references movies,
    "parentId"  integer
        references comments,
    text        text,
    "createdAt" timestamp not null,
    "updatedAt" timestamp
);

alter table comments
    owner to postgres;

create table if not exists "commentsLikes"
(
    "userId"    integer not null
        constraint "commentLikes_userId_fkey"
            references users,
    "commentId" integer not null
        constraint "commentLikes_commentId_fkey"
            references comments,
    "likeType"  varchar(10)
        constraint "commentLikes_likeType_check"
            check (("likeType")::text = ANY
                   (ARRAY [('like'::character varying)::text, ('dislike'::character varying)::text])),
    "createdAt" timestamp default CURRENT_TIMESTAMP,
    constraint "commentLikes_pkey"
        primary key ("userId", "commentId")
);

alter table "commentsLikes"
    owner to postgres;

