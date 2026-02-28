package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	sqlcdb "gorm-sqlc/sqlc"
)

var sqlcLogger = log.New(os.Stdout, "[SQLC] ", log.LstdFlags)

// LoggingQueries wraps sqlc.Queries to log all queries
type LoggingQueries struct {
	*sqlcdb.Queries
}

// NewLoggingQueries creates a new LoggingQueries wrapper
func NewLoggingQueries(queries *sqlcdb.Queries) *LoggingQueries {
	return &LoggingQueries{Queries: queries}
}

// logQuery logs the SQLC query with parameters
func logQuery(queryName string, args ...interface{}) {
	var params []string
	for i, arg := range args {
		if arg == nil {
			params = append(params, "NULL")
		} else {
			val := reflect.ValueOf(arg)
			switch val.Kind() {
			case reflect.Struct:
				// Handle sql.NullString, sql.NullInt32, sql.NullBool
				if val.Type().String() == "sql.NullString" {
					if val.FieldByName("Valid").Bool() {
						params = append(params, fmt.Sprintf("'%v'", val.FieldByName("String").String()))
					} else {
						params = append(params, "NULL")
					}
				} else if val.Type().String() == "sql.NullInt32" {
					if val.FieldByName("Valid").Bool() {
						params = append(params, fmt.Sprintf("%v", val.FieldByName("Int32").Int()))
					} else {
						params = append(params, "NULL")
					}
				} else if val.Type().String() == "sql.NullBool" {
					if val.FieldByName("Valid").Bool() {
						params = append(params, fmt.Sprintf("%v", val.FieldByName("Bool").Bool()))
					} else {
						params = append(params, "NULL")
					}
				} else {
					params = append(params, fmt.Sprintf("%v", arg))
				}
			case reflect.String:
				params = append(params, fmt.Sprintf("'%v'", arg))
			default:
				params = append(params, fmt.Sprintf("%v", arg))
			}
		}
		if i < len(args)-1 {
			params = append(params, ", ")
		}
	}
	
	query := getQuerySQL(queryName)
	if query != "" {
		// Replace $1, $2, etc. with actual values
		for i, param := range params {
			if i < len(params) && params[i] != ", " {
				placeholder := fmt.Sprintf("$%d", i+1)
				query = strings.Replace(query, placeholder, param, 1)
			}
		}
		sqlcLogger.Printf("%s\n", query)
	} else {
		sqlcLogger.Printf("Executing: %s with params: %s\n", queryName, strings.Join(params, ""))
	}
}

// getQuerySQL returns the SQL for a given query name
func getQuerySQL(queryName string) string {
	queries := map[string]string{
		"CreateUser": `INSERT INTO users (email, username, full_name, age, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		"GetUserByID": `SELECT * FROM users WHERE id = $1`,
		"GetUserByEmail": `SELECT * FROM users WHERE email = $1`,
		"ListUsers": `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		"ListActiveUsers": `SELECT * FROM users WHERE is_active = true ORDER BY created_at DESC`,
		"UpdateUser": `UPDATE users SET email = $2, username = $3, full_name = $4, age = $5, updated_at = NOW() WHERE id = $1 RETURNING *`,
		"UpdateUserActiveStatus": `UPDATE users SET is_active = $2, updated_at = NOW() WHERE id = $1`,
		"DeleteUser": `DELETE FROM users WHERE id = $1`,
		"CountUsers": `SELECT COUNT(*) FROM users`,
		"SearchUsersByName": `SELECT * FROM users WHERE full_name ILIKE '%' || $1 || '%' ORDER BY full_name`,
		"GetUsersByAgeRange": `SELECT * FROM users WHERE age BETWEEN $1 AND $2 ORDER BY age`,
	}
	return queries[queryName]
}

// Wrapper methods that log before calling the actual query
func (q *LoggingQueries) CreateUser(ctx context.Context, arg sqlcdb.CreateUserParams) (sqlcdb.User, error) {
	logQuery("CreateUser", arg.Email, arg.Username, arg.FullName, arg.Age, arg.IsActive)
	return q.Queries.CreateUser(ctx, arg)
}

func (q *LoggingQueries) GetUserByID(ctx context.Context, id int32) (sqlcdb.User, error) {
	logQuery("GetUserByID", id)
	return q.Queries.GetUserByID(ctx, id)
}

func (q *LoggingQueries) GetUserByEmail(ctx context.Context, email string) (sqlcdb.User, error) {
	logQuery("GetUserByEmail", email)
	return q.Queries.GetUserByEmail(ctx, email)
}

func (q *LoggingQueries) ListUsers(ctx context.Context, arg sqlcdb.ListUsersParams) ([]sqlcdb.User, error) {
	logQuery("ListUsers", arg.Limit, arg.Offset)
	return q.Queries.ListUsers(ctx, arg)
}

func (q *LoggingQueries) ListActiveUsers(ctx context.Context) ([]sqlcdb.User, error) {
	logQuery("ListActiveUsers")
	return q.Queries.ListActiveUsers(ctx)
}

func (q *LoggingQueries) UpdateUser(ctx context.Context, arg sqlcdb.UpdateUserParams) (sqlcdb.User, error) {
	logQuery("UpdateUser", arg.ID, arg.Email, arg.Username, arg.FullName, arg.Age)
	return q.Queries.UpdateUser(ctx, arg)
}

func (q *LoggingQueries) UpdateUserActiveStatus(ctx context.Context, arg sqlcdb.UpdateUserActiveStatusParams) error {
	logQuery("UpdateUserActiveStatus", arg.ID, arg.IsActive)
	return q.Queries.UpdateUserActiveStatus(ctx, arg)
}

func (q *LoggingQueries) DeleteUser(ctx context.Context, id int32) error {
	logQuery("DeleteUser", id)
	return q.Queries.DeleteUser(ctx, id)
}

func (q *LoggingQueries) CountUsers(ctx context.Context) (int64, error) {
	logQuery("CountUsers")
	return q.Queries.CountUsers(ctx)
}

func (q *LoggingQueries) SearchUsersByName(ctx context.Context, query sql.NullString) ([]sqlcdb.User, error) {
	logQuery("SearchUsersByName", query)
	return q.Queries.SearchUsersByName(ctx, query)
}

func (q *LoggingQueries) GetUsersByAgeRange(ctx context.Context, arg sqlcdb.GetUsersByAgeRangeParams) ([]sqlcdb.User, error) {
	logQuery("GetUsersByAgeRange", arg.Age, arg.Age_2)
	return q.Queries.GetUsersByAgeRange(ctx, arg)
}

func (q *LoggingQueries) WithTx(tx *sql.Tx) *sqlcdb.Queries {
	return q.Queries.WithTx(tx)
}
