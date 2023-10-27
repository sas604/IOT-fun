CREATE TABLE IF NOT EXISTS measurements (
    id INTEGER PRIMARY KEY,
    abbreviation TEXT NOT NULL,
    display_value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS switches_measurements (
    switch_id INTEGER NOT NULL REFERENCES switches ON DELETE CASCADE, 
    measurement_id INTEGER NOT NULL REFERENCES measurements ON DELETE CASCADE,
    PRIMARY KEY (switch_id, measurement_id)
);

INSERT INTO measurements (abbreviation, display_value) 
VALUES 
    ('hum', 'Humidity'),
    ('temp', 'Temperature');

-- INSERT INTO switches_measurements (abbreviation, display_value) 
-- VALUES 
--     ('hum', 'Humidity'),
--     ('temp', 'Temperature');    