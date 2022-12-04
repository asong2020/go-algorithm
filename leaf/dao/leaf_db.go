package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"asong.cloud/go-algorithm/leaf/model"
)

const leafTableName = "leaf_alloc"

type LeafDB struct {
	db *sql.DB
}

func NewLeafDB(db *sql.DB) *LeafDB {
	return &LeafDB{
		db: db,
	}
}

func (l *LeafDB) Create(ctx context.Context, leaf *model.Leaf) error {
	now := time.Now().Unix()
	query := fmt.Sprintf(`INSERT
								INTO
								%s (biz_tag, max_id, step, description, update_time)
								VALUES (?, ?, ?, ?, ?)
								`, leafTableName)

	res, err := l.db.ExecContext(ctx, query, leaf.BizTag, leaf.MaxID, leaf.Step, leaf.Description, uint64(now))
	if err != nil {
		fmt.Printf("insert leaf failed; leaf: %v; err: %v", leaf, err)
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (l *LeafDB) Get(ctx context.Context, bizTag string, tx *sql.Tx) (*model.Leaf, error) {
	query := fmt.Sprintf(`
SELECT 
	id, biz_tag, max_id, step, description,update_time
FROM %s WHERE biz_tag = ?
`, leafTableName)
	var leaf model.Leaf
	var err error
	// 开启事务则用事务
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, bizTag).Scan(&leaf.ID, &leaf.BizTag, &leaf.MaxID,
			&leaf.Step, &leaf.Description, &leaf.UpdateTime)
	} else {
		err = l.db.QueryRowContext(ctx, query, bizTag).Scan(&leaf.ID, &leaf.BizTag, &leaf.MaxID,
			&leaf.Step, &leaf.Description, &leaf.UpdateTime)
	}
	if err != nil {
		fmt.Printf("get leaf failed; biz_tag:%s;err: %v", bizTag, err)
		return nil, err
	}
	return &leaf, nil
}

func (l *LeafDB) UpdateMaxID(ctx context.Context, bizTag string, tx *sql.Tx) error {
	query := fmt.Sprintf(`UPDATE %v SET max_id = max_id + step, update_time = ? WHERE biz_tag = ?`, leafTableName)
	var err error
	var res sql.Result
	now := uint64(time.Now().Unix())
	if tx != nil {
		res, err = tx.ExecContext(ctx, query, now, bizTag)
	} else {
		res, err = l.db.ExecContext(ctx, query, now, bizTag)
	}
	if err != nil {
		fmt.Printf("update max_id failed; bizTag: %s; err: %v", bizTag, err)
		return err
	}
	rowsID, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsID == 0 {
		return errors.New("no update")
	}
	return nil
}

func (l *LeafDB) UpdateMaxIdByCustomStep(ctx context.Context, step int32, bizTag string, tx *sql.Tx) error {
	query := fmt.Sprintf(`UPDATE %v SET max_id = max_id + ?, update_time = ? WHERE biz_tag = ?`, leafTableName)
	now := uint64(time.Now().Unix())
	var err error
	var res sql.Result
	if tx != nil {
		res, err = tx.ExecContext(ctx, query, step, now, bizTag)
	} else {
		res, err = l.db.ExecContext(ctx, query, step, now, bizTag)
	}
	if err != nil {
		fmt.Printf("update max_id by custom step failed; bizTag: %s; err: %v", bizTag, err)
		return err
	}
	rowsID, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsID == 0 {
		return errors.New("no update")
	}
	return nil
}

func (l *LeafDB) GetAll(ctx context.Context) ([]*model.Leaf, error) {
	query := fmt.Sprintf(`
SELECT 
	id, biz_tag, max_id, step, description,update_time
FROM %s `, leafTableName)
	rows, err := l.db.QueryContext(ctx, query)
	if err != nil {
		fmt.Printf("Get all bia_tags failed;err:%v", err)
		return nil, err
	}
	defer func() {
		if rows != nil {
			err = rows.Close()
			if err != nil {
				fmt.Println("close rows failed")
			}
		}
	}()
	list := make([]*model.Leaf, 0)
	for rows.Next() {
		var leaf model.Leaf
		err = rows.Scan(&leaf.ID, &leaf.BizTag, &leaf.MaxID,
			&leaf.Step, &leaf.Description, &leaf.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &leaf)
	}
	return list, nil
}

func (l *LeafDB) UpdateStep(ctx context.Context, step int32, bizTag string) error {
	query := fmt.Sprintf(`UPDATE %v SET step = ?, update_time = ? WHERE biz_tag = ?`, leafTableName)
	now := uint64(time.Now().Unix())
	var err error
	var res sql.Result
	res, err = l.db.ExecContext(ctx, query, step, now, bizTag)
	if err != nil {
		fmt.Printf("update max_id by custom step failed; bizTag: %s; err: %v", bizTag, err)
		return err
	}
	rowsID, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsID == 0 {
		return errors.New("no update")
	}
	return nil
}
