drop MATERIALIZED VIEW if exists rough_totals;
drop table if exists votes;

drop table if exists participants;

CREATE TABLE
    participants (
        participant_id SERIAL PRIMARY KEY,
        participant_name CHAR(100) NOT NULL
    );

INSERT INTO
    participants (participant_name)
VALUES
    ('Isaac Newton'),
    ('Albert Einstein'),
    ('Marie Curie');

CREATE TABLE
    votes (
        vote_id SERIAL PRIMARY KEY,
        participant_id INT NOT NULL,
        timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (participant_id) REFERENCES participants (participant_id)
    );
