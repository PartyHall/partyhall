INSERT INTO event (id, name, date, author, location)
VALUES
    (1, 'Anniversaire toto', datetime('2021-10-23T19:00:00+01:00'), 'Jean TOTO', 'Salle des fêtes'),
    (2, 'Crémaillère tutu', datetime('2022-05-07T19:00:00+01:00'), 'Charles TUTU', 'Chez TUTU')
;

INSERT INTO app_state(id, hwid, token, current_event)
VALUES (1, '73d05017-d8ee-4de8-8229-b8edf452202f', '73d05017-d8ee-4de8-8229-b8edf452202f', 1);