package config

import "database/sql"

func RunMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			full_name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			contact_info VARCHAR(100) NOT NULL,
			password_hash TEXT NOT NULL,
			role VARCHAR(20) DEFAULT 'user',
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(150) NOT NULL,
			category VARCHAR(50) NOT NULL,
			price VARCHAR(50) NOT NULL,
			location VARCHAR(100) NOT NULL,
			contact VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			image_path TEXT,
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS contact_messages (
			id SERIAL PRIMARY KEY,
			full_name VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL,
			subject VARCHAR(150) NOT NULL,
			message TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active';`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS image_path TEXT;`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active';`,

		`UPDATE users SET status = 'active' WHERE status IS NULL;`,
		`UPDATE products SET status = 'active' WHERE status IS NULL;`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}