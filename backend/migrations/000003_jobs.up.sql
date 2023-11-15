CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY,
    interval INTEGER NOT NULL,
    duration INTEGER NOT NULL,
    switch INTEGER NOT NULL UNIQUE,
    on_start TEXT NOT NULL,
    on_end TEXT NOT NULL,
    FOREIGN KEY (switch) REFERENCES switches(id) ON DELETE CASCADE
);

INSERT INTO jobs (interval, duration, switch, on_start, on_end)
VALUES 
    (24 * 60,  12 * 60, 3, 'on', 'off'),
    (3*60, 10, 4, 'on', 'off');