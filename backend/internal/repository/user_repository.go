package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"md2html/internal/model"
)

type UserRepository interface {
	Create(username, passwordHash string) (int64, error)
	FindByUsername(username string) (*model.User, error)
	FindByID(id int64) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(username, passwordHash string) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (username, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?)",
		username, passwordHash, time.Now(), time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by username: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByID(id int64) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by id: %w", err)
	}

	return user, nil
}

// RunMigrations executes database schema migrations
func RunMigrations(db *sql.DB) {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id            BIGINT AUTO_INCREMENT PRIMARY KEY,
			username      VARCHAR(128)  NOT NULL UNIQUE COMMENT '用户名',
			password_hash VARCHAR(255)  NOT NULL COMMENT 'bcrypt 加密密码',
			created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_username (username)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表'`,

		`CREATE TABLE IF NOT EXISTS md2html_conversions (
			id         BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id    BIGINT          NOT NULL COMMENT '用户 ID',
			markdown   MEDIUMTEXT      NOT NULL COMMENT 'Markdown 原文',
			html       MEDIUMTEXT      COMMENT '转换后的 HTML',
			theme      VARCHAR(64)     DEFAULT 'default' COMMENT '使用的主题 ID',
			created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_user_id (user_id),
			INDEX idx_created_at (created_at),
			CONSTRAINT fk_md2html_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='md2html 转换历史记录'`,

		`CREATE TABLE IF NOT EXISTS tool_configs (
			id         BIGINT AUTO_INCREMENT PRIMARY KEY,
			tool_id    VARCHAR(64)     NOT NULL COMMENT '工具标识',
			config_key VARCHAR(128)    NOT NULL COMMENT '配置键',
			config_val TEXT            COMMENT '配置值',
			updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uk_tool_key (tool_id, config_key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工具配置表'`,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			log.Fatalf("Migration %d failed: %v", i+1, err)
		}
	}

	log.Println("Database migrations completed successfully")
}
