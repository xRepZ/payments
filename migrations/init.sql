CREATE USER payments WITH PASSWORD '123321';

CREATE TABLE status_types(
    id varchar(15)  PRIMARY KEY,
    description     TEXT
);

alter TABLE status_types owner to payments;

INSERT INTO status_types(id, description) VALUES 
('new',       'новый'),
('error',     'ошибка'),
('success',   'успех'),
('unsuccess', 'неуспех');


CREATE TABLE transactions
(
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER NOT NULL,
    email               TEXT NOT NULL,
    amount              NUMERIC,
    currency            TEXT,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, 
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, 
    deleted_at          TIMESTAMP,
    status              varchar(15) references status_types(id) 
);

alter TABLE transactions owner to payments;