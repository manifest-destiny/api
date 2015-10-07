CREATE TABLE app_user ( id serial PRIMARY KEY, google_id varchar(32) NOT NULL UNIQUE, email varchar(254) NOT NULL, full_name varchar(512) NOT NULL, alias varchar(256), picture varchar(512), show_picture boolean, locale varchar(5), country varchar(2));
