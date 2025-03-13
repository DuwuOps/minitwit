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
	// Detect the correct primary key dynamically
	primaryKey := detectPrimaryKey(r.tableName)

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", r.tableName, primaryKey)
	row := r.db.QueryRowContext(ctx, query, id)

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

func detectPrimaryKey(tableName string) string {
	switch tableName {
	case "user":
		return "user_id"
	case "message":
		return "message_id"
	case "follower":
		return "who_id" // Follower table has `who_id` and `whom_id`, adjust logic as needed.
	default:
		return "id" // Default to `id`, but this should never happen.
	}
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

func (r *Repository[T]) GetFiltered(ctx context.Context, conditions map[string]any, limit int, orderBy string) ([]T, error) {
	var whereClauses []string
	var values []any

	for key, value := range conditions {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			sliceVal := reflect.ValueOf(value)
			placeholders := make([]string, sliceVal.Len())

			for i := 0; i < sliceVal.Len(); i++ {
				placeholders[i] = "?"
				values = append(values, sliceVal.Index(i).Interface())
			}

			whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ",")))
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", r.tableName, strings.Join(whereClauses, " AND "))
	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}
	if limit > 0 {
		query += " LIMIT ?"
		values = append(values, limit)
	}

	log.Printf("Executing Query: %s | Values: %v", query, values)

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		log.Printf("Query failed: %v", err)
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
			log.Printf("Error scanning row: %v", err)
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

func (r *Repository[T]) DeleteByFields(ctx context.Context, conditions map[string]any) error {
	var whereClauses []string
	var values []any

	for field, value := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
		values = append(values, value)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", r.tableName, strings.Join(whereClauses, " AND "))

	_, err := r.db.ExecContext(ctx, query, values...)
	return err
}


func (r *Repository[T]) Remove(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.tableName)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
