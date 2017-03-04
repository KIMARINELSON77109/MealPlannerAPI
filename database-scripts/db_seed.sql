-- Written by Fred T. Dunaway
-- Feb 2017
-- Made available via the Apache 2.0 license.
-- Please give credit where credit is due.
-- ---------------------------------------------
-- Sample data useful for very basic testing of the API.

use `mooseware_meal_planner`;

INSERT INTO uom (uom) values('oz');
INSERT INTO uom (uom) values('g');
INSERT INTO uom (uom) values('cup');
INSERT INTO uom (uom) values('tablespoon');
INSERT INTO uom (uom) values('teaspoon');
INSERT INTO uom (uom) values('lb');
INSERT INTO uom (uom) values('kg');
INSERT INTO uom (uom) values('l');
INSERT INTO uom (uom) values('n/a');

INSERT INTO meal_type (meal_type) values('breakfast');
INSERT INTO meal_type (meal_type) values('lunch');
INSERT INTO meal_type (meal_type) values('dinner');
INSERT INTO meal_type (meal_type) values('snack');
INSERT INTO meal_type (meal_type) values('brunch');
INSERT INTO meal_type (meal_type) values('dessert');

INSERT INTO owner (owner_name, owner_email) VALUES('Mr. Moose', 'mrncmoose@gmail.com');

INSERT INTO ingredient (item, quantity, uom_id) VALUES('chicken', 1, (select uom_id from uom where uom='lb'));
insert into ingredient (item, quantity, uom_id) VALUES('carrots', 1, (select uom_id from uom where uom='lb'));
insert into ingredient (item, quantity, uom_id) VALUES('potato', 1, (select uom_id from uom where uom='lb'));
insert into ingredient (item, quantity, uom_id) VALUES('whole wheat flour', 1, (select uom_id from uom where uom='cup'));
insert into ingredient (item, quantity, uom_id) VALUES('vanilla yogurt', 1, (select uom_id from uom where uom='cup'));
insert into ingredient (item, quantity, uom_id) VALUES('coco powder', 6, (select uom_id from uom where uom='tablespoon'));
insert into ingredient (item, quantity, uom_id) VALUES('brow sugar', 0.5, (select uom_id from uom where uom='cup'));
insert into ingredient (item, quantity, uom_id) VALUES('eggs', 2, (select uom_id from uom where uom='n/a'));

INSERT INTO recipe (recipe_name, directions, servings, owner_id) VALUES('Brownies', 'mix ingredients in large mixing bowl and bake at 350F for 20-25 minutes.', 
10, 
(select owner_id from owner where owner_email='mrncmoose@gmail.com'));

insert into ingredient_list (recipe_id, ingredient_id) VALUES (
	(select recipe_id from recipe where recipe_name='Brownies'),
	(select ingredient_id from ingredient where item='whole wheat flour'));
insert into ingredient_list (recipe_id, ingredient_id) VALUES (
	(select recipe_id from recipe where recipe_name='Brownies'),
	(select ingredient_id from ingredient where item='vanilla yogurt'));
insert into ingredient_list (recipe_id, ingredient_id) VALUES (
	(select recipe_id from recipe where recipe_name='Brownies'),
	(select ingredient_id from ingredient where item='coco powder'));
insert into ingredient_list (recipe_id, ingredient_id) VALUES (
	(select recipe_id from recipe where recipe_name='Brownies'),
	(select ingredient_id from ingredient where item='brow sugar'));
insert into ingredient_list (recipe_id, ingredient_id) VALUES (
	(select recipe_id from recipe where recipe_name='Brownies'),
	(select ingredient_id from ingredient where item='eggs'));
	
insert into suggested_meal (meal_name, servings, rating, owner_id, recipe_id) values (
'Evening snack', 1, 5, 
(select owner_id from owner where owner_email='mrncmoose@gmail.com'),
(select recipe_id from recipe where recipe_name='Brownies'));

insert into meal_profile (profile_name, suggested_meal_id) VALUES (
"Active", (select suggested_meal_id from suggested_meal where meal_name='Evening snack'));

insert into meal (meal_name, servings, served_on, account_id, meal_type_id, recipe_id) VALUES (
	'Evening snack', 1, '2017-04-01',
	(select owner_id from owner where owner_email='mrncmoose@gmail.com'),
	(select meal_type_id from meal_type where meal_type='snack'),
	(select recipe_id from recipe where recipe_name='Brownies')
);

	