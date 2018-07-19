--           _
--  ___  ___| |_ _   _ _ __
-- / __|/ _ \ __| | | | '_ \
-- \__ \  __/ |_| |_| | |_) |
-- |___/\___|\__|\__,_| .__/
--                    |_|
--
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';

CREATE TYPE rank AS ENUM
    ('user',
    'admin',
    'owner');

--  _        _     _
-- | |_ __ _| |__ | | ___  ___
-- | __/ _` | '_ \| |/ _ \/ __|
-- | || (_| | |_) | |  __/\__ \
--  \__\__,_|_.__/|_|\___||___/
--
CREATE TABLE public.users (
    id bigserial NOT NULL,
    email citext NOT NULL,
    display_name text,
    password bytea DEFAULT ''::bytea NOT NULL,
    rank public.rank DEFAULT 'user'::rank NOT NULL,
    locked boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (id)
);
COMMENT ON TABLE public.users IS
    'contains users that are allowed to sign in to the application.';


CREATE TABLE public.invitations (
    token uuid DEFAULT uuid_generate_v4() NOT NULL,
    email citext NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    valid_until timestamp without time zone DEFAULT now() NOT NULL,
    PRIMARY KEY (token)
);
COMMENT ON TABLE public.invitations IS
    'invitations are tokens that are sent to an email address, which allow
    users to create a new account.';
COMMENT ON COLUMN public.invitations.valid_until IS
    'time at which the invitation will become stale, and can no longer be claimed';
--                      _             _       _
--   ___ ___  _ __  ___| |_ _ __ __ _(_)_ __ | |_ ___
--  / __/ _ \| '_ \/ __| __| '__/ _` | | '_ \| __/ __|
-- | (_| (_) | | | \__ \ |_| | | (_| | | | | | |_\__ \
--  \___\___/|_| |_|___/\__|_|  \__,_|_|_| |_|\__|___/
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE ONLY public.invitations
    ADD CONSTRAINT invitations_email_key UNIQUE (email);