-- Fix tenant table group_id constraint
-- This migration ensures group_id is mandatory and properly constrained

-- Drop the existing foreign key constraint
ALTER TABLE public.tb_tenant 
DROP CONSTRAINT IF EXISTS fk_tenant_group;

-- Ensure group_id is NOT NULL
ALTER TABLE public.tb_tenant 
ALTER COLUMN group_id SET NOT NULL;

-- Re-add the foreign key constraint with proper ON DELETE behavior
ALTER TABLE public.tb_tenant 
ADD CONSTRAINT fk_tenant_group 
FOREIGN KEY (group_id) 
REFERENCES public.tb_tenant_group(id) 
ON UPDATE CASCADE 
ON DELETE RESTRICT;

-- Add an index on group_id for better performance
CREATE INDEX IF NOT EXISTS idx_tenant_group_id ON public.tb_tenant(group_id); 