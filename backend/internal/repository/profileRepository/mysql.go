package profileRepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/joineroff/social-network/backend/internal/entity"
	"github.com/joineroff/social-network/backend/internal/infrastructure/database"
	"github.com/joineroff/social-network/backend/internal/repository"
)

const stmtMysqlGetProfile = `
	SELECT DISTINCT
		u.id,
		u.login,
		u.first_name,
		u.last_name,
		u.gender,
		u.city,
		u.interests,
		uf.user_id as isme
	FROM users as u

	LEFT JOIN user_friends as uf
		ON uf.user_id=? AND uf.friend_id=u.id

	WHERE u.id=?
`

// @TODO DISTINCT unnecessary after adding UNIQUE INDEX (user_id, friend_id)
const stmtMysqlFindProfiles = `
	SELECT DISTINCT
		u.id,
		u.login,
		u.first_name,
		u.last_name,
		u.gender,
		u.city,
		u.interests,
		uf.user_id as isme
	FROM users as u

	LEFT JOIN user_friends as uf
		ON uf.user_id=? AND uf.friend_id=u.id

	WHERE u.id != ?
`

// @TODO DISTINCT unnecessary after adding UNIQUE INDEX (user_id, friend_id)
const stmtMysqlCountProfiles = `
	SELECT COUNT(DISTINCT u.id)
	FROM users as u

	LEFT JOIN user_friends as uf
		ON uf.user_id=? AND uf.friend_id=u.id

	WHERE u.id != ?
`

const (
	stmtMysqlShouldBeFriendCondition = ` AND (uf.user_id IS NOT NULL)`
	stmtMysqlNameLikeCondition       = ` AND (u.first_name LIKE ? AND u.last_name LIKE ?)`
	stmtMysqlLimitOffset             = ` ORDER BY u.id ASC LIMIT ? OFFSET ? `

	mysqlMaxArgs = 6
)

type mysqlProfileRepository struct {
	db *database.MysqlDB
}

func NewMysqlProfileRepository(
	db *database.MysqlDB,
) repository.ProfileRepository {
	return &mysqlProfileRepository{
		db: db,
	}
}

// count profiles satisfied to searchQuery
func (r *mysqlProfileRepository) CountProfiles(
	ctx context.Context,
	searchQuery string,
	currentUserID string,
	requireFriend bool,
) (int, error) {
	var q strings.Builder

	args := make([]interface{}, 0, mysqlMaxArgs)

	q.WriteString(stmtMysqlCountProfiles)

	args = append(args, currentUserID, currentUserID)

	if searchQuery != "" {
		q.WriteString(stmtMysqlNameLikeCondition)

		name := fmt.Sprintf("%s%%", searchQuery)

		args = append(args, name, name)
	}

	if requireFriend {
		q.WriteString(stmtMysqlShouldBeFriendCondition)
	}

	query := r.db.Slave().Rebind(q.String())
	q.Reset()

	var cnt int

	row := r.db.Slave().QueryRowxContext(ctx, query, args...)
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}

	return cnt, nil
}

// find profiles satisfied to searchQuery
func (r *mysqlProfileRepository) GetProfile(
	ctx context.Context,
	id string,
	currentUserID string,
) (*entity.Profile, error) {
	query := r.db.Slave().Rebind(stmtMysqlGetProfile)

	profile := &entity.Profile{}
	profile.User = &entity.User{}

	row := r.db.Slave().QueryRowxContext(ctx, query, currentUserID, id)

	var (
		firstName sql.NullString
		lastName  sql.NullString
		gender    sql.NullInt64
		city      sql.NullString
		interests sql.NullString
		isFriend  sql.NullString
	)

	if err := row.Scan(
		&profile.User.ID,
		&profile.User.Login,
		&firstName,
		&lastName,
		&gender,
		&city,
		&interests,
		&isFriend,
	); err != nil {
		return nil, err
	}

	if firstName.Valid {
		profile.User.FirstName = firstName.String
	}

	if lastName.Valid {
		profile.User.LastName = lastName.String
	}

	if gender.Valid && gender.Int64 > 0 {
		profile.User.Gender = 1
	}

	if city.Valid {
		profile.User.City = city.String
	}

	if interests.Valid {
		profile.User.Interests = strings.Split(interests.String, ",")
	}

	if isFriend.Valid {
		profile.IsFriend = true
	}

	return profile, nil
}

// find profiles satisfied to searchQuery
func (r *mysqlProfileRepository) FindProfiles(
	ctx context.Context,
	searchQuery string,
	currentUserID string,
	requireFriend bool,
	limit int,
	offset int,
) ([]*entity.Profile, error) {
	var q strings.Builder

	args := make([]interface{}, 0, mysqlMaxArgs)

	q.WriteString(stmtMysqlFindProfiles)

	args = append(args, currentUserID, currentUserID)

	if searchQuery != "" {
		q.WriteString(stmtMysqlNameLikeCondition)

		name := fmt.Sprintf("%s%%", searchQuery)

		args = append(args, name, name)
	}

	if requireFriend {
		q.WriteString(stmtMysqlShouldBeFriendCondition)
	}

	q.WriteString(stmtMysqlLimitOffset)

	args = append(args, limit, offset)

	query := r.db.Slave().Rebind(q.String())
	q.Reset()

	rows, err := r.db.Slave().QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	profiles := make([]*entity.Profile, 0, limit)

	var (
		firstName sql.NullString
		lastName  sql.NullString
		gender    sql.NullInt64
		city      sql.NullString
		interests sql.NullString
		isFriend  sql.NullString
	)

	for rows.Next() {
		profile := &entity.Profile{}
		profile.User = &entity.User{}

		if err := rows.Scan(
			&profile.User.ID,
			&profile.User.Login,
			&firstName,
			&lastName,
			&gender,
			&city,
			&interests,
			&isFriend,
		); err != nil {
			return nil, err
		}

		if firstName.Valid {
			profile.User.FirstName = firstName.String
		}

		if lastName.Valid {
			profile.User.LastName = lastName.String
		}

		if gender.Valid && gender.Int64 > 0 {
			profile.User.Gender = 1
		}

		if city.Valid {
			profile.User.City = city.String
		}

		if interests.Valid {
			profile.User.Interests = strings.Split(interests.String, ",")
		}

		if isFriend.Valid {
			profile.IsFriend = true
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
