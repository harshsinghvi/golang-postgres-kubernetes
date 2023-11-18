CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    text VARCHAR ( 50 ) NOT NULL,
    complete BOOLEAN NOT NULL
);

SELECT id,text, complete FROM todos;

UPDATE todos SET text='nae56', WHERE id=1;

SELECT * FROM todos;

SELECT * FROM todos WHERE id=1;

DELETE FROM todos WHERE id=3;

INSERT INTO todos(text,complete) VALUES('task xyz', FALSE) ON CONFLICT DO NOTHING;

SELECT created_at, expiry FROM access_tokens WHERE token='fc19728ee6ee29ccd923379577bf34c2';
UPDATE access_tokens SET expiry='2023-11-10 19:51:38.781473+00' where token='fc19728ee6ee29ccd923379577bf34c2';
