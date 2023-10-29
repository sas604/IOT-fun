CREATE TABLE IF NOT EXISTS switches ( 
    id INTEGER PRIMARY KEY,
     name TEXT NOT NULL, 
     state TEXT NOT NULL,
     topic_base TEXT NOT NULL
     );

INSERT INTO switches (name, state,topic_base)
    VALUES
        ('switch-1', 'off', 'mush/switch-group'),
        ('switch-2', 'off', 'mush/switch-group'),
        ('switch-3', 'off', 'mush/switch-group'),
        ('switch-4', 'off', 'mush/switch-group');
