--
-- PostgreSQL database dump
--

-- Dumped from database version 14.7 (Homebrew)
-- Dumped by pg_dump version 14.7 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: characters; Type: TABLE; Schema: public; Owner: paphos
--

CREATE TABLE public.characters (
    id uuid NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    avatar_id text,
    greeting text NOT NULL,
    persona text NOT NULL,
    example_chats text,
    visibility text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    world_scenario text,
    creator_id uuid NOT NULL,
    contentious boolean NOT NULL
);


ALTER TABLE public.characters OWNER TO paphos;

--
-- Name: chat_participants; Type: TABLE; Schema: public; Owner: paphos
--

CREATE TABLE public.chat_participants (
    id uuid NOT NULL,
    user_id uuid,
    character_id uuid,
    chat_id uuid NOT NULL
);


ALTER TABLE public.chat_participants OWNER TO paphos;

--
-- Name: chats; Type: TABLE; Schema: public; Owner: paphos
--

CREATE TABLE public.chats (
    id uuid NOT NULL,
    owner_id uuid NOT NULL,
    name text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.chats OWNER TO paphos;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: paphos
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO paphos;

--
-- Name: users; Type: TABLE; Schema: public; Owner: paphos
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email text NOT NULL,
    hashed_password text NOT NULL,
    display_name text NOT NULL,
    role text NOT NULL,
    verification_token text,
    password_reset_token text,
    last_login timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.users OWNER TO paphos;

--
-- Name: characters characters_pkey; Type: CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.characters
    ADD CONSTRAINT characters_pkey PRIMARY KEY (id);


--
-- Name: chat_participants chat_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.chat_participants
    ADD CONSTRAINT chat_participants_pkey PRIMARY KEY (id);


--
-- Name: chats chats_pkey; Type: CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: characters_contentious_idx; Type: INDEX; Schema: public; Owner: paphos
--

CREATE INDEX characters_contentious_idx ON public.characters USING btree (contentious);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: paphos
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: paphos
--

CREATE UNIQUE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: characters characters_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.characters
    ADD CONSTRAINT characters_users_id_fk FOREIGN KEY (creator_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: chat_participants chat_participants_character_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.chat_participants
    ADD CONSTRAINT chat_participants_character_id_fkey FOREIGN KEY (character_id) REFERENCES public.characters(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: chat_participants chat_participants_chat_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.chat_participants
    ADD CONSTRAINT chat_participants_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.chats(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: chat_participants chat_participants_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.chat_participants
    ADD CONSTRAINT chat_participants_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: chats chats_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: paphos
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

