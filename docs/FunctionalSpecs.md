# Functional Specifications - Smart Fit Girl

## Purpose
This document consolidates all business logic, nutritional guidelines, and behavioral specifications for the Smart Fit Girl application.

## FOOD_CATALOG Table

### Description
The `FOOD_CATALOG` table maintains nutritional details for various foods, specifying quantities per unit of measure. This data allows users to easily track nutrient intake.

### Nutritional Information
- **Units**: All nutritional values are per one unit of `serving_units` to facilitate calculation for different serving sizes.
- **Measurement Units**:
  - `grams`
  - `ounces`
  - `fluid_ounces`
  - `cup`
  - `tsp`
  - `tbsp`
  - `scoop`
- **Nutrient Breakdown** (per unit):
  - `calories`
  - `protein_grams`
  - `carbs_grams`
  - `fat_grams`

### Inflammatory Response
- **is_non_inflammatory**: Indicates foods that do not cause or may reduce inflammation.

## Conventions

### Database Naming
- **Table Names**: ALL_CAPS
- **Column Names**: snake_case
- **Boolean Columns**: Prefixed with `is_`

## Business Logic & Guidelines
- **Nutritional Tracking**: Calculate total nutrient intake by multiplying per-unit values by serving size.
- **Goal Management**: Support setting and tracking custom fitness goals, including nutrient intake, weight management, and exercise targets.
