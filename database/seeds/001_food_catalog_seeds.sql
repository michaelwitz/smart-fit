-- 001_food_catalog_seeds.sql
-- Sample food data for FOOD_CATALOG table with reorganized columns

INSERT INTO FOOD_CATALOG (food_name, category, serving_units, calories, protein_grams, carbs_grams, fat_grams, is_non_inflammatory, is_probiotic, is_prebiotic, notes) VALUES
-- Fish (all non-inflammatory)
('Salmon - Wild Atlantic', 'Fish', 'ounces', 155, 22.0, 0, 7.0, true, false, false, 'High in omega-3 fatty acids'),
('Cod Fillet', 'Fish', 'ounces', 70, 15.0, 0, 0.5, true, false, false, 'Lean white fish'),
('Tuna - Yellowfin', 'Fish', 'ounces', 92, 20.0, 0, 1.0, true, false, false, 'Low mercury option'),
('Sardines', 'Fish', 'ounces', 125, 15.0, 0, 7.0, true, false, false, 'Small fish with bones for calcium'),

-- Meats (all inflammatory)
('Chicken Breast - Skinless', 'Meat', 'ounces', 140, 26.0, 0, 3.0, false, false, false, 'Lean protein source'),
('Ground Beef - 85% Lean', 'Meat', 'ounces', 160, 22.0, 0, 7.0, false, false, false, 'Moderate fat content'),
('Bison - Ground', 'Meat', 'ounces', 120, 20.0, 0, 4.0, false, false, false, 'Grass-fed lean meat'),
('Pork Tenderloin', 'Meat', 'ounces', 125, 22.0, 0, 3.5, false, false, false, 'Lean cut of pork'),

-- Dairy and Milk Products
('Cheddar Cheese', 'Dairy', 'ounces', 115, 7.0, 1.0, 9.0, false, false, false, 'Aged hard cheese'),
('Goat Cheese', 'Dairy', 'ounces', 75, 5.0, 0, 6.0, true, false, false, 'Easier to digest than cow cheese'),
('Greek Yogurt - Plain', 'Dairy', 'cup', 130, 23.0, 9.0, 0, true, true, false, 'Contains beneficial bacteria'),
('Coconut Milk - Canned', 'Dairy', 'cup', 445, 5.0, 6.0, 48.0, true, false, false, 'High in saturated fats'),
('Goat Milk', 'Dairy', 'cup', 168, 9.0, 11.0, 10.0, true, false, false, 'Alternative to cow milk'),
('A2 Milk', 'Dairy', 'cup', 150, 8.0, 12.0, 8.0, true, false, false, 'A2 protein variant'),
('Cow Milk - Regular', 'Dairy', 'cup', 150, 8.0, 12.0, 8.0, false, false, false, 'Standard dairy milk'),
('Eggs - Large Chicken', 'Dairy', 'item', 70, 6.0, 0.5, 5.0, true, false, false, 'Complete protein source'),
('Eggs - Duck', 'Dairy', 'item', 130, 9.0, 1.0, 10.0, true, false, false, 'Richer flavor than chicken eggs'),
('Eggs - Quail', 'Dairy', 'item', 14, 1.2, 0, 1.0, true, false, false, 'Small gourmet eggs'),

-- Vegetables (non-inflammatory)
('Broccoli', 'Vegetables', 'cup', 25, 3.0, 5.0, 0, true, false, true, 'High in fiber and vitamins'),
('Spinach - Raw', 'Vegetables', 'cup', 7, 1.0, 1.0, 0, true, false, false, 'High in iron and folate'),
('Asparagus', 'Vegetables', 'cup', 20, 2.0, 4.0, 0, true, false, true, 'Natural diuretic properties'),
('Sweet Potato', 'Vegetables', 'cup', 180, 4.0, 41.0, 0, true, false, true, 'High in beta-carotene'),
('Kale', 'Vegetables', 'cup', 33, 2.0, 7.0, 0, true, false, false, 'Superfood high in nutrients'),
('Carrots', 'Vegetables', 'cup', 50, 1.0, 12.0, 0, true, false, true, 'High in beta-carotene'),
('Sauerkraut', 'Vegetables', 'cup', 27, 1.0, 6.0, 0, true, true, false, 'Fermented cabbage with probiotics'),
('Kimchee', 'Vegetables', 'cup', 23, 2.0, 4.0, 0, true, true, false, 'Korean fermented vegetables'),

-- Nightshades (all inflammatory - Solanaceae family)
('Bell Pepper - Red', 'Nightshades', 'cup', 30, 1.0, 7.0, 0, false, false, false, 'Nightshade family - may cause inflammation'),
('Tomatoes', 'Nightshades', 'cup', 32, 2.0, 7.0, 0, false, false, false, 'Nightshade family - acidic'),
('Eggplant', 'Nightshades', 'cup', 20, 1.0, 5.0, 0, false, false, false, 'Nightshade family - contains solanine'),
('Potato - Russet', 'Nightshades', 'cup', 160, 4.0, 37.0, 0, false, false, false, 'Nightshade family - high glycemic'),

-- Grains and Starches
('Brown Rice - Cooked', 'Grains', 'cup', 220, 5.0, 45.0, 2.0, true, false, false, 'Whole grain with fiber'),
('White Rice - Cooked', 'Grains', 'cup', 205, 4.0, 45.0, 0, false, false, false, 'Refined grain - higher glycemic'),
('Quinoa - Cooked', 'Grains', 'cup', 220, 8.0, 39.0, 4.0, true, false, false, 'Complete protein grain'),

-- Nuts and Seeds (mixed)
('Almonds - Whole', 'Nuts', 'ounces', 160, 6.0, 6.0, 14.0, false, false, false, 'Skin may cause inflammation for some'),
('Almonds - Blanched', 'Nuts', 'ounces', 160, 6.0, 6.0, 14.0, true, false, false, 'Skin removed - less inflammatory'),
('Walnuts', 'Nuts', 'ounces', 185, 4.0, 4.0, 18.0, true, false, false, 'Rich in omega-3 fatty acids'),
('Pistachios', 'Nuts', 'ounces', 160, 6.0, 8.0, 13.0, true, false, false, 'High in antioxidants'),
('Chia Seeds', 'Seeds', 'tbsp', 60, 3.0, 5.0, 4.0, true, false, true, 'High in omega-3 and fiber'),
('Flaxseeds - Ground', 'Seeds', 'tbsp', 37, 1.3, 2.0, 3.0, true, false, true, 'Must be ground for absorption'),

-- Legumes (inflammatory for some)
('Peanuts', 'Legumes', 'ounces', 160, 7.0, 5.0, 14.0, false, false, false, 'Legume not a nut - may cause inflammation'),
('Black Beans - Cooked', 'Legumes', 'cup', 230, 15.0, 41.0, 1.0, false, false, false, 'May cause digestive issues'),
('Lentils - Cooked', 'Legumes', 'cup', 230, 18.0, 40.0, 1.0, false, false, false, 'High protein legume'),

-- Oils and Fats
('Olive Oil - Extra Virgin', 'Oils', 'tbsp', 120, 0, 0, 14.0, true, false, false, 'Anti-inflammatory properties'),
('Coconut Oil', 'Oils', 'tbsp', 120, 0, 0, 14.0, true, false, false, 'Medium chain triglycerides'),
('Avocado Oil', 'Oils', 'tbsp', 120, 0, 0, 14.0, true, false, false, 'High smoke point for cooking'),

-- Fruits (mostly non-inflammatory)
('Blueberries', 'Fruits', 'cup', 80, 1.0, 21.0, 0, true, false, false, 'High in antioxidants'),
('Strawberries', 'Fruits', 'cup', 50, 1.0, 12.0, 0, true, false, false, 'High in vitamin C'),
('Avocado', 'Fruits', 'cup', 240, 3.0, 12.0, 22.0, true, false, false, 'Healthy monounsaturated fats'),
('Apple - Medium', 'Fruits', 'cup', 95, 0, 25.0, 0, true, false, true, 'High in pectin fiber'),
('Banana - Medium', 'Fruits', 'cup', 105, 1.0, 27.0, 0, true, false, true, 'Good source of potassium'),
('Grapes', 'Fruits', 'cup', 62, 0, 16.0, 0, true, false, false, 'Natural sugars'),
('Orange', 'Fruits', 'cup', 62, 1.0, 15.0, 0, true, false, false, 'High in vitamin C'),

-- Herbs and Spices (anti-inflammatory)
('Turmeric - Ground', 'Herbs', 'tsp', 8, 0, 1.0, 0, true, false, false, 'Powerful anti-inflammatory compound'),
('Ginger - Fresh', 'Herbs', 'tsp', 1, 0, 0, 0, true, false, false, 'Digestive aid and anti-inflammatory'),
('Garlic - Fresh', 'Herbs', 'tsp', 4, 0, 1.0, 0, true, false, true, 'Immune system support');
