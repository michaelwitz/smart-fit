-- Add 'item' to serving_unit_type ENUM
-- This allows foods to be measured by individual items (e.g., eggs, apples, etc.)

ALTER TYPE serving_unit_type ADD VALUE 'item';
