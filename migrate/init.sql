/* ============================================================
   Pré‑requisitos
   ============================================================ */
-- Gera a função uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/* ============================================================
   1) Tabela: public.tb_tenant_group
   ============================================================ */
CREATE TABLE public.tb_tenant_group (
  id         uuid PRIMARY KEY             DEFAULT uuid_generate_v4(),
  name       varchar(150) NOT NULL,
  cnpj       varchar(30)  NOT NULL UNIQUE,
  is_active  boolean                     DEFAULT true,
  created_at timestamp                   DEFAULT now(),
  updated_at timestamp                   DEFAULT now()
);

/* ============================================================
   2) Tabela: public.tb_tenant
   ============================================================ */
CREATE TABLE public.tb_tenant (
  id          uuid PRIMARY KEY             DEFAULT uuid_generate_v4(),

  group_id    uuid NOT NULL,
  CONSTRAINT  fk_tenant_group
    FOREIGN KEY (group_id) REFERENCES public.tb_tenant_group(id),

  name        varchar(150) NOT NULL,
  cnpj        varchar(30)  NOT NULL UNIQUE,
  schema_name varchar       NOT NULL UNIQUE,
  is_active   boolean                      DEFAULT true,
  created_at  timestamp                    DEFAULT now(),
  updated_at  timestamp                    DEFAULT now()
);

/* ============================================================
   3) Tabela: public.tb_user
   ============================================================ */
CREATE TABLE public.tb_user (
  id               uuid PRIMARY KEY          DEFAULT uuid_generate_v4(),
  id_tanant        uuid NOT NULL,
  CONSTRAINT       fk_user_tenant
    FOREIGN KEY (id_tanant) REFERENCES public.tb_tenant(id),

  username         varchar(50)  NOT NULL UNIQUE,
  name_full        varchar(100) NOT NULL,
  hashed_password  varchar,
  email            varchar(150) NOT NULL UNIQUE,
  enabled          boolean      NOT NULL DEFAULT true,
  created_at       timestamp    NOT NULL DEFAULT now(),
  updated_at       timestamp    NOT NULL DEFAULT now(),
  role_usr         varchar      NOT NULL DEFAULT 'user'
);
