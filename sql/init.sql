-- Always keep the comment at the end so that the init script
-- can split this file correctly

-- Do not use /* */ or // comments as this will break too

DROP TABLE IF EXISTS exported_event;
DROP TABLE IF EXISTS image;
DROP TABLE IF EXISTS event;
DROP TABLE IF EXISTS song;

DROP TABLE IF EXISTS refresh_token;
DROP TABLE IF EXISTS ph_user;
DROP TABLE IF EXISTS app_state;

CREATE TABLE event (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(512) NOT NULL,
    date datetime NOT NULL,
    author VARCHAR(512) NOT NULL,
    location VARCHAR(512) NULL,
    exporting BOOLEAN NOT NULL DEFAULT FALSE,
    last_export datetime NULL
);

CREATE TABLE image (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL REFERENCES event(id),
    unattended BOOLEAN NOT NULL DEFAULT 0 CHECK (unattended IN (0, 1)),
    created_at datetime NOT NULL
);

CREATE TABLE exported_event (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER REFERENCES event(id),
    filename VARCHAR(512) NOT NULL,
    date datetime NOT NULL
);

CREATE TABLE song (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename VARCHAR(512) NOT NULL UNIQUE,
    artist VARCHAR(512) NULL,
    title VARCHAR(512) NULL,
    format VARCHAR(32) NOT NULL DEFAULT 'cdg',
    play_count INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE app_state (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hwid VARCHAR(255) NOT NULL,
    token VARCHAR(255) NULL DEFAULT NULL,
    current_event INTEGER NULL REFERENCES event(id),
    last_applied_migration INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE ph_user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(512),
    username VARCHAR(512) NOT NULL UNIQUE,
    password VARCHAR(512) NOT NULL,
    roles JSON NOT NULL DEFAULT '["USER"]'
);

CREATE TABLE refresh_token (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token VARCHAR(128) NOT NULL UNIQUE,
    expires_at datetime NOT NULL,
    user_id INTEGER NOT NULL REFERENCES ph_user(id)
);

--