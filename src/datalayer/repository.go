package datalayer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Repository[T any] struct {
	db        *sql.DB
	tableName string
}

var ErrRecordNotFound = errors.New("record not found")

var _ IRepository[any] = (*Repository[any])(nil)

func NewRepository[T any](db *sql.DB, tableName string) *Repository[T] {
	return &Repository[T]{db: db, tableName: tableName}
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
    if entity == nil {
        return errors.New("entity is nil")
    }

    val := reflect.ValueOf(entity)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    if !val.IsValid() {
        return errors.New("invalid entity value")
    }

    typeOfEntity := val.Type()
    var columns []string
    var placeholders []string
    var values []interface{}
    paramCount := 0

    for i := 0; i < val.NumField(); i++ {
        field := typeOfEntity.Field(i)
        fieldName := field.Tag.Get("db")
        if fieldName == "" {
            fieldName = field.Name
        }

        // Ensure that only exportable fields are included
        if !val.Field(i).CanInterface() {
            continue
        }

        // Special handling for `author_id`, `who_id`, `whom_id`
        if fieldName == "id" || strings.Contains(fieldName, "_id") {
            if val.Field(i).IsZero() {
                log.Printf("⚠️ Skipping %s (auto-generated or missing)", fieldName)
                continue
            }
        }

        columns = append(columns, fieldName)
        paramCount++
        placeholders = append(placeholders, formatParamIndex(paramCount))
        values = append(values, val.Field(i).Interface())
    }

    if len(columns) == 0 {
        return errors.New("no valid fields to insert")
    }

    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", r.tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))

    log.Printf("📝 Executing INSERT Query: %s | Values: %v", query, values)

    result, err := r.db.ExecContext(ctx, query, values...)
    if err != nil {
        log.Printf("❌ SQL ERROR: Query: %s | Values: %v | Err: %v", query, values, err)
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("✅ Inserted %d row(s) into %s", rowsAffected, r.tableName)

    return nil
}

func (r *Repository[T]) GetByField(ctx context.Context, field string, value any) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", r.tableName, field)
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

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", r.tableName, primaryKey)
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
	case "users":
		return "user_id"
	case "message":
		return "message_id"
	case "follower":
		return "who_id" // Follower table has `who_id` and `whom_id`, adjust logic as needed.
	default:
		return "id" // Default to `id`, but this should never happen.
	}
}

func (r *Repository[T]) GetFiltered(ctx context.Context, conditions map[string]any, limit int, orderBy string) ([]T, error) {
    var whereClauses []string
    var values []any
    paramCount := 0

    for key, value := range conditions {
        if slice, ok := value.([]int); ok {  
            if len(slice) > 0 {
                placeholders := make([]string, len(slice))
                for i, v := range slice {
                    paramCount++
                    placeholders[i] = formatParamIndex(paramCount)
                    values = append(values, v)
                }
                whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ",")))
            }
        } else {
            paramCount++
            whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", key, formatParamIndex(paramCount)))
            values = append(values, value)
        }
    }

    query := fmt.Sprintf("SELECT DISTINCT * FROM %s", r.tableName)
    if len(whereClauses) > 0 {
        query += " WHERE " + strings.Join(whereClauses, " AND ")
    }
    if orderBy != "" {
        query += " ORDER BY " + orderBy
    }
    if limit > 0 {
        paramCount++
        query += fmt.Sprintf(" LIMIT %s", formatParamIndex(paramCount))
        values = append(values, limit)
    }

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

func (r *Repository[T]) DeleteByFields(ctx context.Context, conditions map[string]any) error {
    var whereClauses []string
    var values []any
    paramCount := 0

    for field, value := range conditions {
        paramCount++
        whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", field, formatParamIndex(paramCount)))
        values = append(values, value)
    }

    query := fmt.Sprintf("DELETE FROM %s WHERE %s", r.tableName, strings.Join(whereClauses, " AND "))
    result, err := r.db.ExecContext(ctx, query, values...)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Deleted %d rows from %s where %v", rowsAffected, r.tableName, conditions)
    return nil
}

func formatParamIndex(i int) (string) {
    return fmt.Sprintf("$%d", i)
}