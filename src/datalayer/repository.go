package datalayer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"errors"
)

type Repository[T any] struct {
	db        *sql.DB
	tableName string
}

var dbInstance *sql.DB
var ErrRecordNotFound = errors.New("record not found")

var _ IRepository[any] = (*Repository[any])(nil)

func NewRepository[T any](db *sql.DB, tableName string) *Repository[T] {
	return &Repository[T]{db: db, tableName: tableName}
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	val := reflect.ValueOf(entity).Elem()
	typeOfEntity := val.Type()

	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		fieldName := typeOfEntity.Field(i).Tag.Get("db") 
		if fieldName == "" {
			fieldName = typeOfEntity.Field(i).Name
		}

		if fieldName != "user_id" { 
			columns = append(columns, fieldName)
			placeholders = append(placeholders, "?")
			values = append(values, val.Field(i).Interface())
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", r.tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))
	_, err := r.db.ExecContext(ctx, query, values...)
	return err
}


func (r *Repository[T]) GetByField(ctx context.Context, field string, value any) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", r.tableName, field)
	row := r.db.QueryRowContext(ctx, query, value)

	var entity T
	val := reflect.ValueOf(&entity).Elem()
	fields := make([]any, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}

	err := row.Scan(fields...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound 
		}
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) GetByID(ctx context.Context, id int) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", r.tableName)
	row := r.db.QueryRowContext(ctx, query, id)

	var entity T
	val := reflect.ValueOf(&entity).Elem()

	fields := make([]any, val.NumField()) 
	for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}

	err := row.Scan(fields...)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) GetAll(ctx context.Context) ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var entity T
		val := reflect.ValueOf(&entity).Elem()

		fields := make([]any, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			fields[i] = val.Field(i).Addr().Interface()
		}

		if err := rows.Scan(fields...); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		results = append(results, entity)
	}
	return results, nil
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	val := reflect.ValueOf(entity).Elem()
	typeOfEntity := val.Type()

	var setClauses []string
	var values []interface{}
	var idValue interface{}

	for i := 0; i < val.NumField(); i++ {
		fieldName := typeOfEntity.Field(i).Name
		fieldValue := val.Field(i).Interface()

		if fieldName == "ID" { 
			idValue = fieldValue
		} else {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", fieldName))
			values = append(values, fieldValue)
		}
	}

	if idValue == nil {
		return fmt.Errorf("entity must have an ID field")
	}

	values = append(values, idValue) 

	query := fmt.Sprintf("UPDATE %s SET %s WHERE ID = ?", r.tableName, strings.Join(setClauses, ","))
	_, err := r.db.ExecContext(ctx, query, values...)
	return err
}

func (r *Repository[T]) Remove(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.tableName)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
