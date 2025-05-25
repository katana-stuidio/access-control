-- Ativa extensão para geração de UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabela de Tenants
CREATE TABLE public.tb_tenant (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    cnpj VARCHAR(20) NOT NULL UNIQUE,
    schema_name VARCHAR(200) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Tabela de Usuários
CREATE TABLE public.tb_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_tanant UUID NOT NULL,
    username VARCHAR(50) NOT NULL,
    name_full VARCHAR(120) NOT NULL,
    hashed_password VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    enabled BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    role_usr VARCHAR(30) DEFAULT 'user' NOT NULL,

    -- Restrições
    CONSTRAINT tb_user_username_unique UNIQUE (username),
    CONSTRAINT tb_user_email_unique UNIQUE (email),
    CONSTRAINT fk_user_tenant FOREIGN KEY (id_tanant)
        REFERENCES public.tenant(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
