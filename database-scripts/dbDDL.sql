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

drop table if exists ingredient;
create table ingredient (
	ingredient_id		int not null auto_increment,
	item				varchar(255) not null,
	quantity			FLOAT not null,
	brand				varchar(255), 
	uom_id				int not null,
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
 	served_on		date not null,
 	account_id		int not null, 
 	meal_type_id	int not null,
 	recipe_id		int not null,
 	primary key(meal_id),
 		index recipeIdx (recipe_id),
 		foreign key(recipe_id) references recipe(recipe_id) ON DELETE RESTRICT,
 		foreign key(account_id) REFERENCES owner(owner_id) on DELETE RESTRICT,
 		foreign key(meal_type_id) references meal_type(meal_type_id) ON DELETE RESTRICT
) ENGINE=INNODB;

drop table if exists suggested_meal;
create table suggested_meal(
	suggested_meal_id	int not null auto_increment,	
 	meal_name			varchar(255),
 	servings			int not null,
 	edit_date			datetime default now(),
 	rating				int DEFAULT 1, 
 	picture				blob,
 	owner_id			int not null, 
 	recipe_id			int not null, 
 	primary key(suggested_meal_id),
 		index recipeIdx (recipe_id),
 		foreign key(recipe_id) references recipe(recipe_id) ON DELETE RESTRICT
) ENGINE=INNODB; 

drop table if EXISTS meal_profile;
create table meal_profile (
	meal_profile_id		int not null auto_increment,
	profile_name		varchar(255) UNIQUE,
	suggested_meal_id	int not null,
	primary key(meal_profile_id),
		foreign key(suggested_meal_id) REFERENCES suggested_meal(suggested_meal_id)
) ENGINE=INNODB;
