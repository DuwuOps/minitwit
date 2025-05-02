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

        // Special handling for `author_id`, `follower_id`, `following_id`
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

    rowsAffected, err := r.executeQuery(ctx, query, values...)
    if err == nil {
        log.Printf("✅ Inserted %d row(s) into %s", rowsAffected, r.tableName)
    } else {
        log.Printf("❌ SQL ERROR: Query: %s | Values: %v | Err: %v", query, values, err)
    }

    return err
}

func (r *Repository[T]) GetByField(ctx context.Context, field string, value any) (*T, error) {
    return r.queryRow(ctx, field, value)
}

func (r *Repository[T]) GetByID(ctx context.Context, id int) (*T, error) {
	// Detect the correct primary key dynamically
	primaryKey := detectPrimaryKey(r.tableName)
    return r.queryRow(ctx, primaryKey, id)
}

func (r *Repository[T]) CountAll(ctx context.Context) (int, error) {
    query := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM %s
    `, r.tableName)
    
    row := r.db.QueryRowContext(ctx, query)
    var count int
    if err := row.Scan(&count); err != nil {
        return 0, err
    }

    log.Printf("📝 Executed Query: %s | Values: %v", query, count)
    
    return count, nil
}

func detectPrimaryKey(tableName string) string {
	switch tableName {
	case "users":
		return "user_id"
	case "message":
		return "message_id"
	case "follower":
		return "follower_id" // Follower table has `follower_id` and `following_id`, adjust logic as needed.
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

func (r *Repository[T]) CountRowsWhenGroupedByFieldInRange(ctx context.Context, field string, lower, upper int) (int, error) {
    query := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM (
            SELECT %s, COUNT(*) AS amount
            FROM %s
            GROUP BY %s
        ) sub
        WHERE amount BETWEEN $1 AND $2;
    `, field, r.tableName, field)
    
    row := r.db.QueryRowContext(ctx, query, lower, upper)
    var count int
    if err := row.Scan(&count); err != nil {
        return 0, err
    }
    return count, nil
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
    rowsAffected, err := r.executeQuery(ctx, query, values...)
    if err == nil {
        log.Printf("✅ Deleted %d rows from %s where %v", rowsAffected, r.tableName, conditions)
    }
    
    return err
}

func (r *Repository[T]) SetAllFields(ctx context.Context, updates map[string]any) error {
    var update_fields []string
    var update_values []any

    for update_field, update_value := range updates {
        update_fields = append(update_fields, fmt.Sprintf("%s = ?", update_field))
        update_values = append(update_values, update_value)
    }

    query := fmt.Sprintf("UPDATE %s SET %s", r.tableName, strings.Join(update_fields, ", "))
    rowsAffected, err := r.executeQuery(ctx, query, update_values...)
    if err == nil {
        log.Printf("✅ Updated %d row(s) in %s with %s", rowsAffected, r.tableName, updates)
    }
    
    return err
}

// Query Utils

func (r *Repository[T]) queryRow(ctx context.Context, field string, values ...any) (*T, error) {
    
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", r.tableName, field)
    log.Printf("📝 Executing Query: %s | Values: %v", query, values)

	row := r.db.QueryRowContext(ctx, query, values...)

    var entity T
	val := reflect.ValueOf(&entity).Elem()
	fields := make([]any, val.NumField())

    for i := 0; i < val.NumField(); i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}

	err := row.Scan(fields...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("⚠️ Warning: Query: %s | Values: %v | Err: %v", query, values, err)
			return nil, ErrRecordNotFound
		}
		log.Printf("❌ SQL ERROR: Query: %s | Values: %v | Err: %v", query, values, err)
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) executeQuery(ctx context.Context, query string, values ...any) (int64, error) {
    log.Printf("📝 Executing Query: %s | Values: %v", query, values)

    result, err := r.db.ExecContext(ctx, query, values...)
    if err != nil {
        log.Printf("❌ SQL ERROR: Query: %s | Values: %v | Err: %v", query, values, err)
        return 0, err
    }
    
    rowsAffected, _ := result.RowsAffected()
    
    return rowsAffected, nil
}

func formatParamIndex(i int) (string) {
    return fmt.Sprintf("$%d", i)
}