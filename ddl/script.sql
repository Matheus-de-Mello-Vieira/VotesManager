CREATE TABLE participants (
    participant_id SERIAL PRIMARY KEY,
    participant_name CHAR(100) NOT NULL
);

CREATE TABLE votes {
    vote_id SERIAL PRIMARY KEY,
    participant_id INT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (participant_id) REFERENCES participants(participant_id)
}
