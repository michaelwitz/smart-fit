-- Initial goals data for the smart-fit-girl application
-- This file should be run after the migration 005_create_survey_and_goals_tables.up.sql

INSERT INTO GOALS (category, name, description) VALUES
-- Weight goals
('Weight', 'Lose', 'Lose body weight gradually and sustainably'),
('Weight', 'Maintain', 'Maintain current weight while improving fitness'),
('Weight', 'Gain', 'Gain body weight through muscle building'),

-- Appearance goals  
('Appearance', 'Bulk', 'Build muscle mass and increase overall size'),
('Appearance', 'Lean', 'Build lean muscle with minimal body fat'),
('Appearance', 'Definition/Cut', 'Reduce body fat to show muscle definition'),

-- Strength goals
('Strength', 'Gain', 'Increase strength and power output'),
('Strength', 'Maintain', 'Maintain current strength levels'),

-- Endurance goals
('Endurance', 'Gain', 'Improve cardiovascular endurance and stamina'),
('Endurance', 'Maintain', 'Maintain current endurance fitness levels');
