CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criando a tabela tb_user
CREATE TABLE public.tb_user (
	id uuid DEFAULT uuid_generate_v4() NOT NULL,
	username varchar(50) NOT NULL,
	name_full varchar(120) NOT NULL,
	hashed_password varchar(255) NULL,
	email varchar(255) NULL,
	enabled bool DEFAULT true NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	role_usr varchar(30) DEFAULT 'driver'::character varying NOT NULL,
	CONSTRAINT tb_user_pkey PRIMARY KEY (id),
	CONSTRAINT tb_user_username_unique UNIQUE (username)
);
