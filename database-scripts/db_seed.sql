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
	
insert into meal_profile (profile_name) VALUES ("Active");
insert into meal_profile (profile_name) VALUES ("Weight loss");
insert into meal_profile (profile_name) VALUES ("Anti-inflammatory ");
insert into meal_profile (profile_name) VALUES ("Normal");

insert into meal (meal_name, servings, account_id, meal_type_id) VALUES (
	'Evening snack', 1,
	(select owner_id from owner where owner_email='mrncmoose@gmail.com'),
	(select meal_type_id from meal_type where meal_type='snack')
);

insert into suggested_meal (edit_date, rating, owner_id, meal_id, serve_on) values (
'2017-03-28', 3, 
(select owner_id from owner where owner_email='mrncmoose@gmail.com'),
(select meal_id from meal where meal_name='Evening snack'),
'2017-03-28 21:30:00'
);

insert into meal_recipes (meal_id, recipe_id, meal_profile_id) VALUES (
	(select meal_id from meal where meal_name='Evening snack'),
	(select recipe_id from recipe where recipe_name='Brownies'),
	(SELECT meal_profile_id from meal_profile where profile_name='Active')
);

