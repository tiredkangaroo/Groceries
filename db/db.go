package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

var DATABASE *sql.DB

func loge(qr string) {
	fmt.Println("DATABASE:", qr)
}
func getType(val interface{}) (string, error) {
	switch reflect.TypeOf(val).Kind() {
	case reflect.String:
		return "varchar(255)", nil
	case reflect.Int:
		return "integer", nil
	default:
		return "", errors.New(fmt.Sprintf("Type %s is not a valid type.", val))
	}
}
func ConnectToDB(user string, password string, db_name string) (*sql.DB, error) {
	// connection_info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db_name)
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", user, password, db_name))
	if err != nil {
		return nil, errors.New("The connection to the database failed.")
	}
	DATABASE = db
	return db, err
}
func CreateTableIfNotExists(entity interface{}) (sql.Result, error) {
	structType := reflect.TypeOf(entity)
	if DATABASE == nil {
		return nil, errors.New("The database has not yet been connected to.")
	}
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (`, structType.Name())
	for j := 0; j < structType.NumField(); j++ {
		field := structType.Field(j)
		_type, err := getType(reflect.Zero(field.Type).Interface())
		if err != nil {
			return nil, err
		}
		query += fmt.Sprintf("%s %s", field.Name, _type)
		if j != structType.NumField()-1 {
			query += ","
		}
	}
	query += ");"
	loge(query)
	return DATABASE.Exec(query)
}
func InsertIntoTable(entity interface{}) (sql.Result, error) {
	structType := reflect.TypeOf(entity)
	if DATABASE == nil {
		return nil, errors.New("The database has not yet been connected to.")
	}
	statement := fmt.Sprintf("INSERT INTO %s (", structType.Name())
	secondValues := " VALUES ("
	for j := 0; j < structType.NumField(); j++ {
		field := structType.Field(j)
		name := field.Name
		statement += name
		if reflect.TypeOf(reflect.Zero(field.Type).Interface()).Kind() == reflect.String {
			secondValues += `'`
			secondValues += reflect.ValueOf(entity).Field(j).String()
			secondValues += `'`
		} else {
			secondValues += reflect.ValueOf(entity).Field(j).String()
		}
		if j != structType.NumField()-1 {
			statement += ","
			secondValues += ","
		}
	}
	statement += ")"
	secondValues += ")"
	statement += secondValues
	loge(statement)
	return DATABASE.Exec(statement)
}

func SelectAllFrom(entity interface{}) ([]interface{}, error) {
	entityType := reflect.TypeOf(entity)
	statement := fmt.Sprintf("SELECT * FROM %s", entityType.Name())
	loge(statement)

	rows, err := DATABASE.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []interface{}

	for rows.Next() {
		// Create a new instance of the entity type
		ntity := reflect.New(entityType).Elem()

		// Prepare fields slice with pointers
		fields := make([]interface{}, 0, ntity.NumField())
		for k := 0; k < ntity.NumField(); k++ {
			fields = append(fields, ntity.Field(k).Addr().Interface())
		}

		// Scan the row into the fields
		err := rows.Scan(fields...)
		if err != nil {
			return nil, err
		}

		// Append the populated entity to the result slice
		result = append(result, ntity.Interface())
	}

	return result, nil
}
