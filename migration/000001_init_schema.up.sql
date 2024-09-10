CREATE TABLE task
(
    ID UUID PRIMARY KEY,
    Status VARCHAR(50),
    HTTPStatusCode INT,
    Headers JSONB,
    Length INT
);