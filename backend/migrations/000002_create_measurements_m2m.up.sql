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
    measurement_id INTEGER NOT NULL UNIQUE ,
    max_value NUMERIC NOT NULL,
    min_value NUMERIC NOT NULL,
    display_value TEXT NOT NULL,
    active BOOLEAN NOT NULL CHECK (active IN (0, 1) DEFAULT 1, 
    FOREIGN KEY (measurement_id) REFERENCES measurements (id)

);

INSERT INTO measurements (abbreviation, display_value) 
VALUES 
    ('hum', 'Humidity'),
    ('temp', 'Temperature'),
    ('co', 'CO2');


INSERT INTO switches_measurements (switch_id, measurement_id) 
VALUES
    (2, 2), 
    (1,1);
	
INSERT INTO targets ( measurement_id, max_value, min_value,display_value)
    VALUES
        (2, 25, 23, 'C'),
        (1, 70, 60, '%');


-- SELECT measurements.abbreviation, measurements.display_value, switches.name, targets.max_value, targets.min_value
-- 	FROM measurements
-- 	INNER JOIN switches_measurements  ON switches_measurements.measurement_id = measurements.id
-- 	INNER JOIN switches ON switches_measurements.switch_id = switches.id
--     INNER JOIN targets ON measurements.id = targets.measurement_id
-- 	WHERE measurements.abbreviation = 'temp'    


-- SELECT switches.name, switches.state, CASE WHEN measurements.abbreviation IS NULL THEN  "false" ELSE "true" END AS automation, ifnull(measurements.abbreviation, 'N/A') AS abbreviation, ifnull(measurements.display_value,'N/A') AS display_value, CASE WHEN jobs.id IS NULL THEN  "false" ELSE "true" END AS schedule, ifnull(targets.max_value, 0) AS max_value, ifnull(targets.min_value, 0) AS min_value,  ifnull(targets.display_value, 'N/A') AS display_value,  ifnull(jobs.interval, 0) AS interval,  ifnull(jobs.duration, 0) AS duration
-- FROM switches
-- LEFT JOIN switches_measurements ON switches_measurements.switch_id = switches.id
-- LEFT JOIN measurements ON switches_measurements.measurement_id = measurements.id
-- LEFT JOIN targets ON measurements.id = targets.measurement_id
-- LEFT JOIN jobs ON jobs.switch = switches.id