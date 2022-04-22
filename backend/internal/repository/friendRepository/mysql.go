package friendRepository

import (
	"context"

	"github.com/joineroff/social-network/backend/internal/infrastructure/database"
	"github.com/joineroff/social-network/backend/internal/repository"
)

type mysqlFriendRepository struct {
	db *database.MysqlDB
}

const stmtMysqlCountFriends = `
SELECT COUNT(*) as cnt
FROM user_friends as uf
WHERE uf.user_id=? AND uf.friend_id=?;`

const stmtMysqlInsertFriend = `
INSERT INTO user_friends
(user_id, friend_id)
VALUES (?, ?);`

const stmtMysqlRemoveFriend = `
DELETE FROM user_friends
WHERE user_id=? AND friend_id=?;`

func NewMysqlFriendRepository(
	db *database.MysqlDB,
) repository.FriendRepository {
	return &mysqlFriendRepository{db}
}

// Add friendID as friend to userID
// If already exists ignore
func (r *mysqlFriendRepository) Add(ctx context.Context, userID, friendID string) error {
	tx, err := r.db.Master().Beginx()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := r.db.Master().Rebind(stmtMysqlInsertFriend)

	_, err = tx.ExecContext(ctx, query, userID, friendID)
	if err != nil {
		// @TODO return nil duplicate error after add UNQUE INDEX (user_id, friend_id)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Remove user's friend
// If not exists ignore
func (r *mysqlFriendRepository) Remove(ctx context.Context, userID, friendID string) error {
	tx, err := r.db.Master().Beginx()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := r.db.Master().Rebind(stmtMysqlRemoveFriend)

	_, err = tx.ExecContext(ctx, query, userID, friendID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Count user's friends
func (r *mysqlFriendRepository) Count(ctx context.Context, userID string) (int, error) {
	query := r.db.Master().Rebind(stmtMysqlCountFriends)

	var cnt int

	row := r.db.Master().QueryRowxContext(ctx, query, userID)
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}

	return cnt, nil
}
