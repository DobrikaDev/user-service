-- +goose Up
-- +goose StatementBegin
CREATE TABLE reputation_groups (
    id INT PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    coefficient NUMERIC(4,2) NOT NULL,
    reputation_need INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

INSERT INTO reputation_groups (id, name, description, coefficient, reputation_need)
VALUES
    (1, 'Новичок', 'Первая ступень коммьюнити: участник только присоединился и изучает основы.', 0.80, 0),
    (2, 'Добряк', 'Активный участник, помогает другим и регулярно участвует в инициативах сообщества.', 1.00, 100),
    (3, 'Наставник', 'Опытный участник, который делится знаниями, сопровождает новичков и развивает проекты.', 1.20, 300),
    (4, 'Посол Добрики', 'Лидер сообщества, представляет Добрику вовне и вдохновляет других своим примером.', 1.30, 500);

ALTER TABLE users
    ADD CONSTRAINT users_reputation_group_id_fkey
    FOREIGN KEY (reputation_group_id) REFERENCES reputation_groups(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_reputation_group_id_fkey;
DROP TABLE reputation_groups;
-- +goose StatementEnd
