package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	todo := &model.TODO{ID: id}

	row := s.db.QueryRowContext(ctx, confirm, id)
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	todos := []*model.TODO{}

	if prevID == 0 {
		//prevIDの指定なし
		rows, err := s.db.QueryContext(ctx, read, size)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			todo := &model.TODO{}
			err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			todos = append(todos, todo)
			//なぜか警告文が出る．静的解析ツールによる誤検出？(2023/3/11)
			// -> （解決）関数の最後に返り値としてtodosを返してやると警告文がなくなる
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
			return nil, err
		}
	} else {
		//prevIDの指定あり
		rowsWithID, err := s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rowsWithID.Close()

		for rowsWithID.Next() {
			todo := &model.TODO{}
			err := rowsWithID.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			todos = append(todos, todo)
		}
		if err := rowsWithID.Err(); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	num, err := result.RowsAffected()
	if num == 0 {
		return nil, &model.ErrNotFound{}
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	todo := &model.TODO{}
	todo.ID = id
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	idSlice := make([]interface{}, len(ids))
	for i, num := range ids {
		idSlice[i] = num
	}
	//以下のコードは誤用(makeでスライスを作った場合，空のスライスではなく0で初期化されたスライスが渡される)
	// for _,num:=range ids {
	// 	idSlice = append(idSlice, num)
	// }

	result, err := s.db.ExecContext(ctx, fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1)), idSlice...)
	if err != nil {
		log.Println(err)
		return err
	}

	num, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if num == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
