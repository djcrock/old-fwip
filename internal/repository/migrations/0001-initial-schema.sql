CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    username text NOT NULL
) STRICT;

CREATE UNIQUE INDEX uix_user__username ON user(username);

INSERT INTO user (username) VALUES ('david'), ('hannah');

CREATE TABLE title (
    id           INTEGER PRIMARY KEY,
    imdb_id      TEXT    NOT NULL,
    type         TEXT    NOT NULL,
    name         TEXT    NOT NULL,
    year         INTEGER NOT NULL,
    release_date TEXT NOT NULL,
    runtime      INTEGER NOT NULL
) STRICT;

CREATE UNIQUE INDEX uix_title__imdb_id ON title (imdb_id);

CREATE TABLE service (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
) STRICT;

INSERT INTO service (id, name)
VALUES
(1, 'Amazon Prime'),
(2, 'Disney+'),
(3, 'Hulu'),
(4, 'Crunchyroll'),
(5, 'Max'),
(6, 'Paramount+'),
(7, 'Netflix');

CREATE TABLE service_title (
    service_id INTEGER NOT NULL,
    title_id INTEGER NOT NULL,
    FOREIGN KEY (service_id) REFERENCES service(id),
    FOREIGN KEY (title_id) REFERENCES title(id),
    PRIMARY KEY (service_id, title_id)
) STRICT, WITHOUT ROWID;

CREATE INDEX ix_provider_title__title_id__service_id ON service_title(title_id, service_id);

CREATE TABLE watch_history (
    user_id INTEGER NOT NULL,
    title_id INTEGER NOT NULL,
    watched INTEGER NOT NULL,
    want_to_watch INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (title_id) REFERENCES title(id),
    PRIMARY KEY (user_id, title_id)
) STRICT, WITHOUT ROWID;
