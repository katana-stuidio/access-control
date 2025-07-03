-- Migration to add TenantGroup support
-- This migration adds the tenant_group table and updates the tenant table

-- Create the tenant_group table
CREATE TABLE public.tb_tenant_group (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    cnpj VARCHAR(20) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Add group_id column to the tenant table
ALTER TABLE public.tb_tenant 
ADD COLUMN group_id UUID;

-- Add foreign key constraint for group_id
ALTER TABLE public.tb_tenant 
ADD CONSTRAINT fk_tenant_group 
FOREIGN KEY (group_id) 
REFERENCES public.tb_tenant_group(id) 
ON UPDATE CASCADE 
ON DELETE SET NULL;

-- Create index for better performance on group_id lookups
CREATE INDEX idx_tenant_group_id ON public.tb_tenant(group_id);

-- Create index for better performance on tenant_group cnpj lookups
CREATE INDEX idx_tenant_group_cnpj ON public.tb_tenant_group(cnpj); 