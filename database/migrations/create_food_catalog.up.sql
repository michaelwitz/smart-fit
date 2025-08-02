-- 001_create_reorganized_food_catalog.up.sql
-- Single, consolidated migration for FOOD_CATALOG table

-- Create enum type for serving units (if not exists)
DO $$ BEGIN
    CREATE TYPE serving_unit_type AS ENUM (
        'grams',
        'ounces',
        'fluid_ounces',
        'cup',
        'tsp',
        'tbsp',
        'scoop'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create FOOD_CATALOG table with complete, reorganized structure
CREATE TABLE FOOD_CATALOG (
    id SERIAL PRIMARY KEY,
    food_name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'Other',
    serving_units serving_unit_type NOT NULL,
    calories DECIMAL(8,2) NOT NULL DEFAULT 0,
    protein_grams DECIMAL(8,2) NOT NULL DEFAULT 0,
    carbs_grams DECIMAL(8,2) NOT NULL DEFAULT 0,
    fat_grams DECIMAL(8,2) NOT NULL DEFAULT 0,
    is_non_inflammatory BOOLEAN NOT NULL DEFAULT false,
    is_probiotic BOOLEAN NOT NULL DEFAULT false,
    is_prebiotic BOOLEAN NOT NULL DEFAULT false,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_food_catalog_food_name ON FOOD_CATALOG(food_name);
CREATE INDEX idx_food_catalog_category ON FOOD_CATALOG(category);
CREATE INDEX idx_food_catalog_non_inflammatory ON FOOD_CATALOG(is_non_inflammatory);
CREATE INDEX idx_food_catalog_probiotic ON FOOD_CATALOG(is_probiotic);
CREATE INDEX idx_food_catalog_prebiotic ON FOOD_CATALOG(is_prebiotic);

-- Add comments for documentation
COMMENT ON TABLE FOOD_CATALOG IS 'Catalog of foods with nutritional information per unit of measure';
COMMENT ON COLUMN FOOD_CATALOG.food_name IS 'Name of the food item';
COMMENT ON COLUMN FOOD_CATALOG.category IS 'Food category (e.g., Fish, Meat, Dairy, Vegetables, Fruits, etc.)';
COMMENT ON COLUMN FOOD_CATALOG.serving_units IS 'Unit of measure for the food serving (grams, ounces, fluid_ounces, cup, tsp, tbsp, scoop)';
COMMENT ON COLUMN FOOD_CATALOG.calories IS 'Calories per one unit of serving_units';
COMMENT ON COLUMN FOOD_CATALOG.protein_grams IS 'Protein in grams per one unit of serving_units';
COMMENT ON COLUMN FOOD_CATALOG.carbs_grams IS 'Carbohydrates in grams per one unit of serving_units';
COMMENT ON COLUMN FOOD_CATALOG.fat_grams IS 'Fat in grams per one unit of serving_units';
COMMENT ON COLUMN FOOD_CATALOG.is_non_inflammatory IS 'Foods that typically do not cause inflammatory response or may reduce inflammation';
COMMENT ON COLUMN FOOD_CATALOG.is_probiotic IS 'Foods that contain beneficial live bacteria for gut health';
COMMENT ON COLUMN FOOD_CATALOG.is_prebiotic IS 'Foods that feed beneficial gut bacteria and promote their growth';
COMMENT ON COLUMN FOOD_CATALOG.notes IS 'Additional notes about the food item, preparation methods, or health benefits';
