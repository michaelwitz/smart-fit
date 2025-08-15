# Smart Fit - Functional Specifications

## Application Overview

Smart Fit is an AI-powered fitness and nutrition application designed to help users achieve their fitness goals through personalized meal planning and goal tracking. The application takes into account individual fitness goals, biological sex (for hormone-related considerations), and personal preferences to generate customized meal plans that support optimal health and fitness outcomes.

## Core Purpose

The primary intention of Smart Fit is to bridge the gap between fitness goals and nutrition by providing users with:

1. **Personalized Goal Setting**: Clear, achievable fitness objectives
2. **AI-Generated Meal Plans**: Customized nutrition strategies based on goals and biological factors
3. **Hormone-Aware Planning**: Meal recommendations that consider sex-specific hormonal needs
4. **Progress Tracking**: Comprehensive monitoring of fitness and nutrition progress
5. **Educational Content**: Understanding of how nutrition impacts fitness goals

## Target Users

- **Fitness Enthusiasts**: Individuals actively working toward specific fitness goals
- **Health-Conscious Individuals**: People seeking to optimize their nutrition for better health
- **Goal-Oriented Users**: Those who benefit from structured, measurable objectives
- **Hormone-Aware Individuals**: Users who understand the impact of biological factors on fitness

## Fitness Goals System

### Goal Categories

#### **Weight Management**

- **Lose**: Gradual, sustainable weight loss
- **Maintain**: Current weight maintenance while improving fitness
- **Gain**: Healthy weight gain through muscle building

#### **Appearance Goals**

- **Bulk**: Build muscle mass and increase overall size
- **Lean**: Build lean muscle with minimal body fat
- **Definition/Cut**: Reduce body fat to show muscle definition

#### **Strength Goals**

- **Gain**: Increase strength and power output
- **Maintain**: Maintain current strength levels

#### **Endurance Goals**

- **Gain**: Improve cardiovascular endurance and stamina
- **Maintain**: Maintain current endurance fitness levels

### Goal Assignment

- Users can select multiple goals across different categories
- Goals are directly assigned to users (immediate activation)
- No survey requirement for goal setting
- Goals can be modified or updated at any time

## Biological Sex Considerations

### Sex Enum Values

- **MALE**: Male biological sex
- **FEMALE**: Female biological sex
- **OTHER**: Non-binary, intersex, or other biological variations

### Hormone-Aware Meal Planning

#### **Male Considerations**

- Higher protein requirements for muscle building
- Testosterone optimization through nutrition
- Muscle recovery and growth support
- Energy metabolism optimization

#### **Female Considerations**

- Iron and calcium requirements
- Hormonal cycle considerations
- Bone health and osteoporosis prevention
- Energy balance during different phases

#### **Other Considerations**

- Individualized approach based on specific needs
- Flexible nutritional recommendations
- Personalized health considerations

## Database Schema Diagram

```
┌──────────────────────────┐
│          USERS           │
├──────────────────────────┤
│ id (PK)                  │
│ email                    │
│ username                 │
│ password_hash            │
│ full_name                │
│ sex                      │ (MALE, FEMALE, OTHER)
│ phone_number             │
│ address_line_1           │
│ address_line_2           │
│ city                     │
│ state_province_code      │
│ country_code             │
│ postal_code              │
│ locale                   │
│ timezone                 │
│ utc_offset               │
│ created_at               │
│ updated_at               │
└──────────────────────────┘
        │
        │ 1:many
        ├────────────────────┐
        │                    │
        ▼                    ▼
┌─────────────────┐  ┌─────────────────┐
│   USER_GOALS    │  │  FOOD_USER_LIKES│
├─────────────────┤  ├─────────────────┤
│ id (PK)         │  │ id (PK)         │
│ user_id (FK)    │  │ user_id (FK)    │
│ goal_id (FK)    │  │ food_id (FK)    │
│ created_at      │  │ created_at      │
└─────────────────┘  └─────────────────┘
        │                    │
        ▼                    ▼
┌─────────────────┐  ┌─────────────────┐
│      GOALS      │  │  FOOD_CATALOG   │
├─────────────────┤  ├─────────────────┤
│ id (PK)         │  │ id (PK)         │
│ category        │  │ food_name       │
│ name            │  │ category        │
│ description     │  │ serving_units   │
│ created_at      │  │ calories        │
│ updated_at      │  │ protein_grams   │
└─────────────────┘  │ carbs_grams     │
        │            │ fat_grams       │
        │            │ is_non_inflammatory│
        │            │ is_probiotic    │
        │            │ is_prebiotic    │
        │            │ notes           │
        │            │ created_at      │
        │            │ updated_at      │
        │            └─────────────────┘
        │                    │
        │                    │ 1:many
        │                    ▼
        │            ┌─────────────────┐
        │            │ MEAL_INGREDIENTS│
        │            ├─────────────────┤
        │            │ id (PK)         │
        │            │ meal_id (FK)    │
        │            │ food_id (FK)    │
        │            │ quantity        │
        │            │ unit            │
        │            │ notes           │
        │            │ created_at      │
        │            └─────────────────┘
        │                    │
        │                    ▼
        │            ┌─────────────────┐
        │            │      MEALS      │
        │            ├─────────────────┤
        │            │ id (PK)         │
        │            │ name            │
        │            │ description     │
        │            │ total_calories  │
        │            │ total_protein   │
        │            │ total_carbs     │
        │            │ total_fat       │
        │            │ prep_time       │
        │            │ created_at      │
        │            │ updated_at      │
        │            └─────────────────┘
        │                    │
        │                    │ 1:many
        │                    ▼
        │            ┌─────────────────┐
        │            │   USER_MEALS    │
        │            ├─────────────────┤
        │            │ id (PK)         │
        │            │ user_id (FK)    │
        │            │ meal_id (FK)    │
        │            │ date            │
        │            │ meal_number     │ (1-6 for meal ordering throughout the day)
        │            │ created_at      │
        │            └─────────────────┘
        │
        └────────────────────┘
```

### **Database Relationships**

- **User ←→ Goals** (many-to-many via USER_GOALS): Direct user goal assignment
- **User ←→ FoodCatalog** (many-to-many via FOOD_USER_LIKES): User food preferences
- **User ←→ Meals** (many-to-many via USER_MEALS): User meal consumption tracking
- **Meals ←→ FoodCatalog** (many-to-many via MEAL_INGREDIENTS): Meal composition with quantities

### **Association Tables**

- **USER_GOALS**: Links users to their selected fitness goals
- **FOOD_USER_LIKES**: Links users to foods they like/prefer
- **USER_MEALS**: Links users to meals with date and meal_number tracking
- **MEAL_INGREDIENTS**: Links meals to food items with quantities and units

### **Enum Values**

- **Sex**: MALE, FEMALE, OTHER
- **Goal Categories**: Weight, Appearance, Strength, Endurance
- **Food Categories**: MEAT, FISH, GRAIN, VEGETABLE, FRUIT, DAIRY, DAIRY_ALTERNATIVE, FAT, NIGHTSHADES, OIL, SPICE_HERB, SWEETENER, CONDIMENT, SNACK, BEVERAGE, LEGUMES, NUTS, SEEDS, OTHER
- **Serving Units**: GRAMS, OUNCES, TSP, TBSP, CUPS, PIECES

---

## Database Table Descriptions

### **USERS Table**

Stores user profiles with international support and authentication:

- **id**: Primary key (auto-increment)
- **email**: Login identifier (unique, required, immutable after creation)
- **password**: Bcrypt hashed password (required)
- **full_name**: User's display name (required)
- **phone_number**: Optional contact number
- **sex**: Biological sex (MALE, FEMALE, OTHER)
- **city**: User's city
- **state_province**: State or province
- **postal_code**: Postal/ZIP code for international support
- **country_code**: ISO 2-letter country code
- **locale**: Language/locale preference (e.g., "en-US", "es-ES")
- **timezone**: IANA timezone identifier (e.g., "America/New_York")
- **utc_offset**: UTC offset in hours for convenience
- **created_at**: Account creation timestamp
- **updated_at**: Last profile update timestamp

### **GOALS Table**

Stores available fitness goals organized by categories:

- **id**: Primary key (auto-increment)
- **category**: ENUM ('Weight', 'Appearance', 'Strength', 'Endurance')
- **name**: Goal name (max 50 chars for UI display)
- **description**: Goal description (max 200 chars)
- **created_at**: Goal creation timestamp
- **updated_at**: Last goal update timestamp

### **USER_GOALS Table**

Junction table linking users to their selected goals:

- **id**: Primary key (auto-increment)
- **user_id**: Foreign key to USERS table
- **goal_id**: Foreign key to GOALS table
- **created_at**: Goal assignment timestamp
- **UNIQUE constraint**: (user_id, goal_id) - prevents duplicate goal assignments

### **FOOD_CATALOG Table**

Comprehensive food database with nutritional information and health properties:

- **id**: Primary key (auto-increment)
- **food_name**: Food item name (required)
- **category**: Food category from enum (required)
- **serving_units**: Unit of measurement from enum (required)
- **calories**: Calories per unit (required)
- **protein_grams**: Protein content per unit (required)
- **carbs_grams**: Carbohydrate content per unit (required)
- **fat_grams**: Fat content per unit (required)
- **is_non_inflammatory**: Boolean flag indicating anti-inflammatory properties
- **is_probiotic**: Boolean flag indicating probiotic content
- **is_prebiotic**: Boolean flag indicating prebiotic content
- **notes**: Additional food information
- **created_at**: Food catalog entry creation timestamp
- **updated_at**: Last update timestamp

### **FOOD_USER_LIKES Table**

Junction table tracking user food preferences:

- **id**: Primary key (auto-increment)
- **user_id**: Foreign key to USERS table
- **food_id**: Foreign key to FOOD_CATALOG table
- **created_at**: Preference recording timestamp
- **UNIQUE constraint**: (user_id, food_id) - prevents duplicate preferences

### **MEALS Table**

Stores meal definitions with nutritional totals and preparation instructions:

- **id**: Primary key (auto-increment)
- **name**: Meal name (required)
- **description**: Meal description
- **total_calories**: Total calories for the meal
- **total_protein**: Total protein content
- **total_carbs**: Total carbohydrate content
- **total_fat**: Total fat content
- **prep_time**: Preparation time in minutes
- **prep_instructions**: HTML or Markdown formatted cooking instructions (nullable)
- **created_at**: Meal creation timestamp
- **updated_at**: Last update timestamp

### **MEAL_INGREDIENTS Table**

Junction table defining meal composition with quantities:

- **id**: Primary key (auto-increment)
- **meal_id**: Foreign key to MEALS table
- **food_id**: Foreign key to FOOD_CATALOG table
- **quantity**: Amount of food item
- **unit**: Unit of measurement from serving_units enum
- **notes**: Additional ingredient information
- **created_at**: Ingredient addition timestamp

### **USER_MEALS Table**

Tracks user meal consumption by date and meal number:

- **id**: Primary key (auto-increment)
- **user_id**: Foreign key to USERS table
- **meal_id**: Foreign key to MEALS table
- **date**: Date of meal consumption (required)
- **meal_number**: Meal order (1-6 for meal ordering throughout the day)
- **created_at**: Consumption recording timestamp
- **UNIQUE constraint**: (user_id, meal_id, date, meal_number) - prevents duplicate meal records

---

## Food Categories & Nutritional Framework

### **MEAT** (Protein-Rich Animal Products)

- **Purpose**: Primary protein source for muscle building and repair
- **Examples**: Chicken Breast, Ground Beef, Pork Tenderloin, Bison
- **Serving Units**: ounces
- **Key Nutrients**: High protein, moderate fat, zero carbs
- **Goal Alignment**: Essential for strength, bulk, and lean muscle goals

### **FISH** (Aquatic Protein Sources)

- **Purpose**: Lean protein with essential fatty acids
- **Examples**: Salmon, Cod, Tuna, Sardines
- **Serving Units**: ounces
- **Key Nutrients**: High protein, omega-3 fatty acids, low fat
- **Goal Alignment**: Ideal for lean muscle and endurance goals

### **GRAIN** (Complex Carbohydrates)

- **Purpose**: Primary energy source and fiber
- **Examples**: Brown Rice, Quinoa, White Rice
- **Serving Units**: cup
- **Key Nutrients**: Complex carbs, fiber, moderate protein
- **Goal Alignment**: Essential for endurance and energy support

### **VEGETABLE** (Nutrient-Dense Plant Foods)

- **Purpose**: Vitamins, minerals, and fiber
- **Examples**: Broccoli, Spinach, Kale, Asparagus, Carrots
- **Serving Units**: cup
- **Key Nutrients**: Low calories, high fiber, rich in micronutrients
- **Goal Alignment**: Support for all fitness goals, health optimization

### **FRUIT** (Natural Sugars and Antioxidants)

- **Purpose**: Quick energy and antioxidant support
- **Examples**: Apples, Bananas, Blueberries, Strawberries, Oranges
- **Serving Units**: cup or PIECES
- **Key Nutrients**: Natural sugars, fiber, vitamins, antioxidants
- **Goal Alignment**: Pre/post workout fuel, overall health

### **DAIRY** (Calcium and Protein)

- **Purpose**: Calcium, protein, and probiotic support
- **Examples**: Milk, Cheese, Yogurt, Eggs
- **Serving Units**: cup, ounces, or PIECES
- **Key Nutrients**: High protein, calcium, vitamin D
- **Goal Alignment**: Muscle building, bone health, recovery

### **DAIRY_ALTERNATIVE** (Plant-Based Dairy Substitutes)

- **Purpose**: Dairy-free alternatives for lactose-intolerant users
- **Examples**: Coconut Milk, Goat Milk
- **Serving Units**: cup
- **Key Nutrients**: Varied protein content, alternative calcium sources
- **Goal Alignment**: Inclusive nutrition for all users

### **FAT** (Essential Fatty Acids)

- **Purpose**: Hormone production and nutrient absorption
- **Examples**: Avocado, Nuts, Seeds
- **Serving Units**: ounces or cup
- **Key Nutrients**: Healthy fats, moderate protein
- **Goal Alignment**: Hormone optimization, satiety

### **NIGHTSHADES** (Inflammatory Considerations)

- **Purpose**: Awareness of potential inflammatory foods
- **Examples**: Bell Peppers, Eggplant, Potatoes, Tomatoes
- **Serving Units**: cup
- **Key Nutrients**: Varied, but may cause inflammation in sensitive users
- **Goal Alignment**: Customized based on individual tolerance

### **OIL** (Cooking and Dressing Fats)

- **Purpose**: Essential fatty acids and cooking medium
- **Examples**: Olive Oil, Coconut Oil, Avocado Oil
- **Serving Units**: tbsp
- **Key Nutrients**: Pure fat, 120 calories per tablespoon
- **Goal Alignment**: Hormone support, nutrient absorption

### **SPICE_HERB** (Flavor and Health Benefits)

- **Purpose**: Anti-inflammatory properties and flavor enhancement
- **Examples**: Garlic, Ginger, Turmeric
- **Serving Units**: tsp
- **Key Nutrients**: Low calories, medicinal properties
- **Goal Alignment**: Recovery, inflammation management

### **SWEETENER** (Natural Sweetening Options)

- **Purpose**: Alternative to refined sugars
- **Examples**: Honey, Maple Syrup, Stevia
- **Serving Units**: tsp or tbsp
- **Key Nutrients**: Varied glycemic impact
- **Goal Alignment**: Energy management, blood sugar control

### **CONDIMENT** (Flavor Enhancement)

- **Purpose**: Meal enjoyment and variety
- **Examples**: Mustard, Hot Sauce, Vinegar
- **Serving Units**: tsp or tbsp
- **Key Nutrients**: Low calories, flavor compounds
- **Goal Alignment**: Meal satisfaction, adherence

### **SNACK** (Convenient Nutrition)

- **Purpose**: Between-meal nutrition and energy
- **Examples**: Protein Bars, Nuts, Dried Fruit
- **Serving Units**: PIECES or ounces
- **Key Nutrients**: Balanced macros for sustained energy
- **Goal Alignment**: Consistent energy, goal adherence

### **BEVERAGE** (Hydration and Nutrition)

- **Purpose**: Fluid intake and supplemental nutrition
- **Examples**: Protein Shakes, Green Tea, Smoothies
- **Serving Units**: cup or PIECES
- **Key Nutrients**: Varied based on beverage type
- **Goal Alignment**: Hydration, recovery, convenience

### **LEGUMES** (Plant-Based Protein)

- **Purpose**: Alternative protein source for plant-based diets
- **Examples**: Black Beans, Lentils, Peanuts
- **Serving Units**: cup or ounces
- **Key Nutrients**: High protein, fiber, complex carbs
- **Goal Alignment**: Plant-based nutrition, fiber support

### **NUTS** (Healthy Fats and Protein)

- **Purpose**: Essential fatty acids and protein
- **Examples**: Almonds, Walnuts, Pistachios
- **Serving Units**: ounces
- **Key Nutrients**: Healthy fats, moderate protein, fiber
- **Goal Alignment**: Hormone support, satiety, recovery

### **SEEDS** (Micro-Nutrient Powerhouses)

- **Purpose**: Essential fatty acids and fiber
- **Examples**: Chia Seeds, Flaxseeds
- **Serving Units**: tbsp
- **Key Nutrients**: Omega-3 fatty acids, fiber, minerals
- **Goal Alignment**: Hormone optimization, digestive health

### **OTHER** (Miscellaneous Foods)

- **Purpose**: Catch-all for unique or specialized foods
- **Examples**: Supplements, specialty products
- **Serving Units**: Varied
- **Key Nutrients**: Depends on specific food
- **Goal Alignment**: Individualized nutrition needs

## Serving Units Standardization

### **GRAMS**: Precise weight measurements

- **Use Case**: Dry ingredients, supplements
- **Precision**: Highest accuracy for nutritional calculations

### **OUNCES**: Standard weight measurements

- **Use Case**: Meat, fish, nuts, cheese
- **Precision**: Standard for protein-rich foods

### **TSP**: Small volume measurements

- **Use Case**: Spices, herbs, small amounts
- **Precision**: 1 tsp = 4.93 ml

### **TBSP**: Medium volume measurements

- **Use Case**: Oils, seeds, condiments
- **Precision**: 1 tbsp = 14.79 ml

### **CUPS**: Large volume measurements

- **Use Case**: Vegetables, fruits, grains, liquids
- **Precision**: 1 cup = 236.59 ml

### **PIECES**: Individual item counts

- **Use Case**: Eggs, whole fruits, individual items
- **Precision**: Exact count for consistent portions

## AI Meal Planning Algorithm

### **Input Factors**

1. **User Goals**: Primary and secondary fitness objectives
2. **Biological Sex**: Hormone-related nutritional considerations
3. **Activity Level**: Current exercise frequency and intensity
4. **Food Preferences**: Liked/disliked foods from user history
5. **Health Considerations**: Non-inflammatory requirements, allergies
6. **Meal Timing**: Preferred eating schedule and workout timing

### **Output Generation**

1. **Daily Calorie Targets**: Based on goals and activity level
2. **Macro Distribution**: Protein, carbs, and fat ratios
3. **Meal Structure**: Breakfast, lunch, dinner, and snacks
4. **Food Selection**: Specific foods from appropriate categories
5. **Portion Sizing**: Exact quantities based on serving units
6. **Preparation Instructions**: Cooking methods and timing

### **Optimization Criteria**

- **Goal Alignment**: Meals support stated fitness objectives
- **Hormone Balance**: Nutrition supports sex-specific hormonal needs
- **Inflammation Management**: Minimizes inflammatory foods for sensitive users
- **Variety**: Diverse food selection for nutritional completeness
- **Practicality**: Realistic preparation time and availability
- **Sustainability**: Long-term adherence and enjoyment

## User Experience Flow

### **1. Goal Setting**

- User selects primary and secondary fitness goals
- System explains goal implications and timeline
- Goals are immediately activated

### **2. Profile Completion**

- Biological sex selection (MALE, FEMALE, OTHER)
- Activity level assessment
- Food preference indication
- Health consideration identification

### **3. AI Meal Generation**

- Algorithm processes user inputs
- Generates personalized meal plan
- Provides nutritional breakdown
- Suggests preparation methods

### **4. Plan Customization**

- User can modify suggested meals
- Food substitutions based on preferences
- Portion adjustments for individual needs
- Schedule modifications

### **5. Progress Tracking**

- Meal adherence monitoring
- Goal progress assessment
- Nutritional intake analysis
- Fitness outcome correlation

### **6. Continuous Optimization**

- AI learns from user feedback
- Meal plan adjustments based on results
- Seasonal food availability updates
- New goal integration

## Success Metrics

### **User Engagement**

- Daily meal plan usage
- Food preference updates
- Goal modification frequency
- Long-term retention rates

### **Goal Achievement**

- Weight change tracking
- Strength improvement measurements
- Endurance enhancement
- Body composition changes

### **Nutritional Adherence**

- Calorie target compliance
- Macro ratio achievement
- Food category variety
- Inflammation reduction

### **Health Outcomes**

- Energy level improvements
- Sleep quality enhancement
- Recovery time reduction
- Overall well-being increase

## Technical Implementation

### **Data Requirements**

- Comprehensive food database with accurate nutritional information
- User preference and goal tracking systems
- AI algorithm for meal plan generation
- Progress monitoring and analytics

### **Integration Points**

- Fitness tracking devices and apps
- Food delivery and grocery services
- Social sharing and community features
- Healthcare provider communication

### **Scalability Considerations**

- Multiple user support
- Diverse goal and preference combinations
- Seasonal and regional food variations
- Continuous learning and improvement

This functional specification provides the foundation for developing Smart Fit as a comprehensive, AI-powered fitness and nutrition platform that truly personalizes the user experience based on individual goals, biological factors, and preferences.
