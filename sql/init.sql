-- Always keep the comment at the end so that the init script
-- can split this file correctly

-- Do not use /* */ or // comments as this will break too

DROP TABLE IF EXISTS exported_event;
DROP TABLE IF EXISTS image;
DROP TABLE IF EXISTS event;
DROP TABLE IF EXISTS karaoke_image;
DROP TABLE IF EXISTS karaoke_song_session;
DROP TABLE IF EXISTS karaoke_song;

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

CREATE TABLE karaoke_song (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid VARCHAR(255) NOT NULL UNIQUE, -- @TODO: Make it correct
    spotify_id VARCHAR(255),
    artist VARCHAR(512) NULL,
    title VARCHAR(512) NULL,
    hotspot VARCHAR(512) NULL,
    format VARCHAR(32) NOT NULL DEFAULT 'cdg',

    has_cover BOOL NOT NULL DEFAULT FALSE,
    has_vocals BOOL NOT NULL DEFAULT FALSE,
    has_full BOOL NOT NULL DEFAULT FALSE,

    filename VARCHAR(512) NOT NULL UNIQUE
);

CREATE TABLE karaoke_song_session (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    song_id INTEGER NOT NULL REFERENCES karaoke_song(id),
    event_id INTEGER NOT NULL REFERENCES event(id),
    sung_by VARCHAR(255) NOT NULL,
    added_at datetime NULL DEFAULT NULL,
    started_at datetime NULL DEFAULT NULL,
    ended_at datetime NULL DEFAULT NULL,
    cancelled_at datetime NULL DEFAULT NULL
);

CREATE TABLE karaoke_image (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    song_session_id INTEGER NOT NULL REFERENCES karaoke_song_session(id),
    created_at datetime NOT NULL
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