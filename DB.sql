-- Table: public.channels

-- DROP TABLE public.channels;

CREATE TABLE  IF NOT EXISTS public.channels
(
    created_at timestamp without time zone NOT NULL,
    description character varying(400) COLLATE pg_catalog."default" NOT NULL,
    id bigint NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    url character varying(255) COLLATE pg_catalog."default" NOT NULL,
    last_entry_id bigint NOT NULL,
    CONSTRAINT channels_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.channels
    OWNER to root;


-- Table: public.sensor_values

-- DROP TABLE public.sensor_values;

CREATE TABLE  IF NOT EXISTS public.sensor_values
(
    entry_id bigint NOT NULL,
    value character varying(100) COLLATE pg_catalog."default",
    created_at timestamp without time zone NOT NULL,
    sensor_id bigint NOT NULL,
    CONSTRAINT sensor_values_pkey PRIMARY KEY (entry_id, sensor_id)
)

TABLESPACE pg_default;

ALTER TABLE public.sensor_values
    OWNER to root;
-- Index: entry_id

-- DROP INDEX public.entry_id;

CREATE INDEX IF NOT EXISTS entry_id
    ON public.sensor_values USING btree
    (entry_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: sensor_id

-- DROP INDEX public.sensor_id;

CREATE INDEX IF NOT EXISTS sensor_id
    ON public.sensor_values USING btree
    (sensor_id ASC NULLS LAST)
    TABLESPACE pg_default;


-- Table: public.sensors

-- DROP TABLE public.sensors;

CREATE TABLE  IF NOT EXISTS public.sensors
(
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    channel_id bigint NOT NULL,
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    CONSTRAINT sensors_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public.sensors
    OWNER to root;
-- Index: channel_id

-- DROP INDEX public.channel_id;

CREATE INDEX IF NOT EXISTS channel_id
    ON public.sensors USING btree
    (channel_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: name

-- DROP INDEX public.name;

CREATE INDEX IF NOT EXISTS name
    ON public.sensors USING btree
    (name COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS name_channel_id
    ON public.sensors USING btree
    (name COLLATE pg_catalog."default" ASC NULLS LAST, channel_id ASC NULLS LAST)
    TABLESPACE pg_default;

-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE  IF NOT EXISTS public.users
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    username character varying(255) COLLATE pg_catalog."default" NOT NULL,
    password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    token text COLLATE pg_catalog."default",
    "tokenHash" character varying(255) COLLATE pg_catalog."default",
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT username_unique UNIQUE (username)
)

TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to root;
-- Index: tokenHash

-- DROP INDEX public."tokenHash";

CREATE INDEX IF NOT EXISTS "tokenHash"
    ON public.users USING btree
    ("tokenHash" COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: username

-- DROP INDEX public.username;

CREATE INDEX IF NOT EXISTS username
    ON public.users USING btree
    (username COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;


-- Table: public.push_settings

-- DROP TABLE public.push_settings;


CREATE TABLE  IF NOT EXISTS public.push_settings
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    sensor_id bigint NOT NULL,
    target_device_id character varying(255) COLLATE pg_catalog."default" NOT NULL,
    target_sensor_id character varying(255) COLLATE pg_catalog."default" NOT NULL,
    active boolean,
    last_pushed_entry_id bigint,
    push_interval integer NOT NULL,
    last_push_time timestamp without time zone,
    use_original_time boolean,
    pushed_count bigint NOT NULL DEFAULT 0,
    CONSTRAINT push_settings_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;


ALTER TABLE public.push_settings
    OWNER to root;
-- Index: active

-- DROP INDEX public.active;

CREATE INDEX IF NOT EXISTS active
    ON public.push_settings USING btree
    (active ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: push_interval

-- DROP INDEX public.push_interval;

CREATE INDEX IF NOT EXISTS push_interval
    ON public.push_settings USING btree
    (push_interval ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: user_sensor

-- DROP INDEX public.user_sensor;

CREATE INDEX IF NOT EXISTS user_sensor
    ON public.push_settings USING btree
    (user_id ASC NULLS LAST, sensor_id ASC NULLS LAST)
    TABLESPACE pg_default;