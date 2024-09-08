-- Enable the uuid-ossp extension to generate UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/*
Create users table with the following columns:
id - uuid primary key
*/

CREATE TABLE users (
  id UUID PRIMARY KEY
);

/*
Create reminders table with the following columns:
id - uuid primary key
date - bigint
send_alert - boolean default false
*/

CREATE TABLE reminders (
  id UUID PRIMARY KEY,
  date BIGINT,
  send_alert BOOLEAN DEFAULT FALSE
);


/*
Create task table with the following columns:
id - uuid primary key
user_id - uuid foreign key
body - text
completed - boolean default false
parent - uuid foreign key
reminder - uuid foreign key
*/

CREATE TABLE tasks (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  body TEXT,
  completed BOOLEAN DEFAULT FALSE,
  parent UUID REFERENCES tasks(id),
  reminder UUID REFERENCES reminders(id)
);

/*
populate the users table with one user
*/

INSERT INTO users (id) VALUES (uuid_generate_v4());

/*
populate the tasks table 3 tasks, all belonging to the user created above
one of the tasks should have a parent task
*/

INSERT INTO tasks (id, user_id, body, completed) VALUES (uuid_generate_v4(), (SELECT id FROM users), 'task 1', false);
INSERT INTO tasks (id, user_id, body, completed) VALUES (uuid_generate_v4(), (SELECT id FROM users), 'task 2', false);
INSERT INTO tasks (id, user_id, body, completed, parent) VALUES (uuid_generate_v4(), (SELECT id FROM users), 'task 3', false, (SELECT id FROM tasks WHERE body = 'task 1'));


