CREATE TABLE IF NOT EXISTS measurements (
    id INTEGER PRIMARY KEY,
    abbreviation TEXT UNIQUE NOT NULL,
    display_value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS switches_measurements (
    switch_id INTEGER NOT NULL REFERENCES switches ON DELETE CASCADE, 
    measurement_id INTEGER NOT NULL REFERENCES measurements ON DELETE CASCADE,
    PRIMARY KEY (switch_id, measurement_id)
);

CREATE TABLE IF NOT EXISTS targets (
    id INTEGER PRIMARY KEY,
    measurement_id INTEGER NOT NULL UNIQUE,
    max_value NUMERIC NOT NULL,
    min_value NUMERIC NOT NULL,
    display_value TEXT NOT NULL,
    FOREIGN KEY (measurement_id) REFERENCES measurements (id)

);

INSERT INTO measurements (abbreviation, display_value) 
VALUES 
    ('hum', 'Humidity'),
    ('temp', 'Temperature');


INSERT INTO switches_measurements (switch_id, measurement_id) 
VALUES
    (2, 2), 
    (1,1);
	
INSERT INTO targets ( measurement_id, max_value, min_value,display_value)
    VALUES
        (2, 28, 27, 'C');


-- SELECT measurements.abbreviation, measurements.display_value, switches.name, targets.max_value, targets.min_value
-- 	FROM measurements
-- 	INNER JOIN switches_measurements  ON switches_measurements.measurement_id = measurements.id
-- 	INNER JOIN switches ON switches_measurements.switch_id = switches.id
--     INNER JOIN targets ON measurements.id = targets.measurement_id
-- 	WHERE measurements.abbreviation = 'temp'    