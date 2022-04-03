package userRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joineroff/social-network/backend/internal/entity"
	"github.com/joineroff/social-network/backend/internal/infrastructure/database"
	"github.com/joineroff/social-network/backend/internal/repository"
)

var _ repository.UserRepository = &MysqlUserRepository{}

type MysqlUserRepository struct {
	db *database.MysqlDB
}

const stmtMysqlCountUsers = `SELECT COUNT(u.id) as cnt FROM users as u`

const stmtMysqlInsertUser = `
INSERT INTO users (
	login,
	password,
	first_name,
	last_name,
	interests,
	city,
	gender
) VALUES (?, ?, ?, ?, ?, ?, ?)
`

const stmtMysqlDeleteUser = `
DELETE FROM users as u (
WHERE u.id = ?
`

const stmtMysqlFindUserByID = `
SELECT
	u.id,
	u.login,
	u.password,
	u.first_name,
	u.last_name,
	u.interests,
	u.city,
	u.gender
FROM users as u
WHERE u.id=?
LIMIT 1;
`

const stmtMysqlFindUserByLogin = `
SELECT
	u.id,
	u.login,
	u.password,
	u.first_name,
	u.last_name,
	u.interests,
	u.city,
	u.gender
FROM users as u
WHERE u.login=?
LIMIT 1;
`

const stmtMysqlFindAllUsers = `
SELECT
	u.id,
	u.login,
	u.password,
	u.first_name,
	u.last_name,
	u.interests,
	u.city,
	u.gender
FROM users as u
LIMIT ? OFFSET ?;
`

func NewMysqlUserRepository(
	db *database.MysqlDB,
) *MysqlUserRepository {
	return &MysqlUserRepository{
		db: db,
	}
}

func (r *MysqlUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := r.db.Rebind(stmtMysqlFindUserByID)

	row := r.db.QueryRowxContext(ctx, query, id)

	user, err := r.scanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *MysqlUserRepository) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	query := r.db.Rebind(stmtMysqlFindUserByLogin)

	row := r.db.QueryRowxContext(ctx, query, login)

	user, err := r.scanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *MysqlUserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := r.db.Rebind(stmtMysqlFindAllUsers)

	rows, err := r.db.QueryxContext(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	defer rows.Close()

	users, err := r.scanUserRows(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *MysqlUserRepository) Count(ctx context.Context) (int, error) {
	var cnt int

	row := r.db.QueryRowxContext(ctx, stmtMysqlCountUsers)
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}

	return cnt, nil
}

func (r *MysqlUserRepository) Create(ctx context.Context, model *entity.User) (*entity.User, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := r.db.Rebind(stmtMysqlInsertUser)

	interests := strings.Join(model.Interests, ",")

	result, err := tx.ExecContext(
		ctx,
		query,
		&model.Login,
		&model.Password,
		&model.FirstName,
		&model.LastName,
		&interests,
		&model.City,
		&model.Gender,
	)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	model.ID = fmt.Sprint(lastID)

	return model, nil
}

func (r *MysqlUserRepository) Update(ctx context.Context, model *entity.User) (*entity.User, error) {
	panic("not implemented") // TODO: Implement
}

func (r *MysqlUserRepository) Delete(ctx context.Context, model *entity.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := r.db.Rebind(stmtMysqlDeleteUser)

	if _, err := tx.ExecContext(
		ctx,
		query,
		&model.ID,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *MysqlUserRepository) scanUserRow(row *sqlx.Row) (*entity.User, error) {
	user := entity.User{}

	var (
		firstName sql.NullString
		lastName  sql.NullString
		city      sql.NullString
		interests sql.NullString
		gender    sql.NullInt64
	)

	if err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&firstName,
		&lastName,
		&interests,
		&city,
		&gender,
	); err != nil {
		return nil, err
	}

	if firstName.Valid {
		user.FirstName = firstName.String
	}

	if lastName.Valid {
		user.LastName = lastName.String
	}

	if interests.Valid {
		user.Interests = strings.Split(interests.String, ",")
	}

	if city.Valid {
		user.City = city.String
	}

	if gender.Valid && gender.Int64 > 0 {
		user.Gender = 1
	}

	return &user, nil
}

func (r *MysqlUserRepository) scanUserRows(rows *sqlx.Rows) ([]*entity.User, error) {
	users := make([]*entity.User, 0)

	for rows.Next() {
		user := entity.User{}

		var (
			firstName sql.NullString
			lastName  sql.NullString
			city      sql.NullString
			interests sql.NullString
			gender    sql.NullInt64
		)

		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Password,
			&firstName,
			&lastName,
			&interests,
			&city,
			&gender,
		); err != nil {
			return nil, err
		}

		if firstName.Valid {
			user.FirstName = firstName.String
		}

		if lastName.Valid {
			user.LastName = lastName.String
		}

		if interests.Valid {
			user.Interests = strings.Split(interests.String, ",")
		}

		if city.Valid {
			user.City = city.String
		}

		if gender.Valid && gender.Int64 > 0 {
			user.Gender = 1
		}

		users = append(users, &user)
	}

	return users, nil
}
