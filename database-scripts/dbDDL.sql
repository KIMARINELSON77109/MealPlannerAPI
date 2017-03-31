-- Written by Fred T. Dunaway
-- Feb 2017
-- Made available via the Apache 2.0 license.
-- Please give credit where credit is due.

use `mooseware_meal_planner`;

SET foreign_key_checks = 0;

drop table if exists uom;
create table uom (
	uom_id			int not null auto_increment,
	uom				varchar(255) UNIQUE not null,
	primary key(uom_id)
) ENGINE=INNODB;

drop table if exists meal_type;
create table meal_type (
	meal_type_id	int not null auto_increment,
	meal_type		varchar(255) UNIQUE not null,
	primary key(meal_type_id)	
) ENGINE=INNODB;

-- the key value assigned to this ingredient by the USDA.  Appears to be a number, but can have leading zeros.
drop table if exists ingredient;
create table ingredient (
	ingredient_id		int not null auto_increment,
	item				varchar(255) not null,
	quantity			FLOAT not null,
	brand				varchar(255), 
	uom_id				int not null,
	usda_ndbno			varchar(255),
	picture				MEDIUMBLOB,
	primary key(ingredient_id),
		index uomIdx (uom_id),
		foreign key(uom_id) references uom(uom_id) ON DELETE RESTRICT	
) ENGINE=INNODB;

drop table if exists owner;
create table owner (
	owner_id		int not null auto_increment UNIQUE,
	owner_name		varchar(255) not null,
	owner_email		varchar(3000) not null UNIQUE,
	primary key (owner_id)
) ENGINE=INNODB;

drop table if exists recipe;
create table recipe (
	recipe_id					int not null auto_increment,		
	recipe_name					varchar(255) not null,
	directions					text not null,
	servings					int not null,
	owner_id					int not null,
	picture						MEDIUMBLOB,
	primary key(recipe_id),
		index ownerIdx (owner_id),
		foreign key(owner_id) references owner(owner_id) ON DELETE RESTRICT
) ENGINE=INNODB;

drop table if exists ingredient_list;
create table ingredient_list (
	ingredient_list_id	int not null auto_increment,	
	recipe_id			int not null,
	ingredient_id		int not null,
	primary key(ingredient_list_id),
		index recipeIdx (recipe_id),
		foreign key(recipe_id) references recipe(recipe_id) ON DELETE RESTRICT,
		foreign key(ingredient_id) references ingredient(ingredient_id) ON DELETE RESTRICT	
) ENGINE=INNODB;

drop table if exists meal;
create table meal(
 	meal_id			int not null auto_increment,
 	meal_name		varchar(255) not null,
 	servings		int not null, 
 	account_id		int not null, 
 	meal_type_id	int not null,
 	primary key(meal_id),
 		foreign key(account_id) REFERENCES owner(owner_id) on DELETE RESTRICT,
 		foreign key(meal_type_id) references meal_type(meal_type_id) ON DELETE RESTRICT
) ENGINE=INNODB;

drop table if EXISTS meal_profile;
create table meal_profile (
	meal_profile_id		int not null auto_increment,
	profile_name		varchar(255) UNIQUE,
	primary key(meal_profile_id)
) ENGINE=INNODB;

-- one meal: 1..* recipes
drop table if exists meal_recipes;
create table meal_recipes (
	meal_recipes_id		int not null auto_increment,
	meal_id				int not null,
	recipe_id			int not null,
 	meal_profile_id		int not null,
	PRIMARY key(meal_recipes_id),
	FOREIGN key(meal_id) REFERENCES meal(meal_id) on DELETE RESTRICT,
	FOREIGN key(recipe_id) REFERENCES recipe(recipe_id) on DELETE RESTRICT,
	FOREIGN KEY(meal_profile_id) REFERENCES meal_profile(meal_profile_id)
) ENGINE=INNODB;

drop table if exists suggested_meal;
create table suggested_meal(
	suggested_meal_id	int not null auto_increment,	
 	edit_date			datetime default now(),
 	rating				int DEFAULT 1, 
 	owner_id			int not null, 
 	meal_id				int not null,
 	serve_on			datetime,
 	primary key(suggested_meal_id),
 		foreign key(meal_id) references meal(meal_id) ON DELETE RESTRICT,
 		FOREIGN key (owner_id) REFERENCES owner(owner_id) on DELETE RESTRICT
) ENGINE=INNODB; 

