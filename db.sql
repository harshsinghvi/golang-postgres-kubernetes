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