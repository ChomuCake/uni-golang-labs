-- migration/000002_initial_user_expenses.up

-- Створення таблиці користувачів
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- Створення таблиці витра
CREATE TABLE expenses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    date TIMESTAMP NOT NULL,
    category VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- migrate -path migration -database "mysql://root:12345@tcp(localhost:3306)/test" up 