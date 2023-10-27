CREATE TABLE IF NOT EXISTS switches ( 
    id INTEGER PRIMARY KEY,
     name TEXT NOT NULL, 
     state TEXT NOT NULL);

INSERT INTO switches (name, state)
    VALUES
        ('switch-1', 'off'),
        ('switch-2', 'off'),
        ('switch-3', 'off'),
        ('switch-4', 'off');
