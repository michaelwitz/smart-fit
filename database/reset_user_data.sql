-- Reset User-Related Development Data
-- This script truncates user-related tables and rebuilds from seed data
-- WARNING: This will delete ALL user, survey, and goal data!
-- FOOD_CATALOG is preserved and not affected

-- Disable foreign key checks temporarily
SET session_replication_role = replica;

-- Truncate user-related tables in dependency order
TRUNCATE TABLE USER_SURVEY_GOALS RESTART IDENTITY CASCADE;
TRUNCATE TABLE SURVEYS RESTART IDENTITY CASCADE;
TRUNCATE TABLE USERS RESTART IDENTITY CASCADE;
TRUNCATE TABLE GOALS RESTART IDENTITY CASCADE;

-- Re-enable foreign key checks
SET session_replication_role = DEFAULT;

-- Note: Run seed files separately using:
-- docker exec -i smart-fit-postgres psql -U smartfit -d smartfitgirl < database/seeds/001_initial_goals.sql
-- docker exec -i smart-fit-postgres psql -U smartfit -d smartfitgirl < database/seeds/002_test_users.sql

-- Show final counts
\echo 'User data reset complete. Final counts:'
SELECT 'Goals' as table_name, COUNT(*) as count FROM GOALS
UNION ALL
SELECT 'Users' as table_name, COUNT(*) as count FROM USERS
UNION ALL
SELECT 'Surveys' as table_name, COUNT(*) as count FROM SURVEYS
UNION ALL
SELECT 'Survey Goals' as table_name, COUNT(*) as count FROM USER_SURVEY_GOALS
UNION ALL
SELECT 'Food Items (preserved)' as table_name, COUNT(*) as count FROM FOOD_CATALOG;
