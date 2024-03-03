package database

func getInitialSchemaQueries() string {
	return `create table if not exists thunderbyte_settings
(
    id  serial not null
        constraint table_name_pk
            primary key,
    key   text    not null
        constraint key_unique
            unique,
    value text
);

insert into thunderbyte_settings (key, value) values ('version', '1.0.0');


create table if not exists auth_users (
    id serial not null
        constraint users_pk
            primary key,
    username text not null
        constraint users_username_unique
            unique
);

create table if not exists auth_passwords (
    id serial not null
        constraint user_passwords_pk
            primary key,
    user_id integer not null
        constraint user_passwords_user_id_fk
            references auth_users,
    password text not null
);`
}

func getDefaultRepoQueries() string {
	return `
	-- queries.sql

-- name: get-all-settings
SELECT * FROM thunderbyte_settings;

-- name: get-setting-by-key
SELECT * FROM thunderbyte_settings where key=$1;

-- name: verify-creds
SELECT 
	au.id as auth_uid, 
	au.username as username
	FROM auth_users as au
	left join auth_passwords as ap
	where au.username='$1' and ap.password='$2';
`
}