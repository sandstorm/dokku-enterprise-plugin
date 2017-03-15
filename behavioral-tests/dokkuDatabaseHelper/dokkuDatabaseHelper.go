package dokkuDatabaseHelper

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"github.com/sandstorm/dokku-enterprise-plugin/core/utility"
	"regexp"
)

func Execute(databaseName, query string) (sql.Result, error) {
	db, err := initializeDatabase(databaseName)
	if err != nil {
		return nil, fmt.Errorf("dokku: Could not create connection to database %s: %v", databaseName, err)
	}
	defer db.Close()

	// Prepare statement for inserting data
	res, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("dokku: Could not execute query on database %s: %v", databaseName, err)
	}

	return res, nil
}

func QueryString(databaseName, query string) (string, error) {
	db, err := initializeDatabase(databaseName)
	if err != nil {
		return "", fmt.Errorf("dokku: Could not create connection to database %s: %v", databaseName, err)
	}
	defer db.Close()

	// Prepare statement for retrieving data
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("dokku: Could not execute query on database %s: %v", databaseName, err)
	}
	defer rows.Close()

	var queryResult string
	if rows.Next() {
		err = rows.Scan(&queryResult)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("dokku: Query did not return any results for database %s", databaseName)
	}

	return queryResult, nil
}

func initializeDatabase(databaseName string) (*sql.DB, error) {
	// First, we need to expose the database to the outside
	utility.ExecCommand("ssh", "dokku@dokku.me", "mariadb:expose", databaseName)

	// Then, we need to figure out the username and password, as well as the port the database is listening to;
	userWithPasswordRegex := regexp.MustCompile(`Dsn:\s*mysql://(\w*:\w*)`)
	portRegex := regexp.MustCompile(`Exposed ports:\s*\d*->(\d*)`)

	databaseInfo := utility.ExecCommand("ssh", "dokku@dokku.me", "mariadb:info", databaseName)

	userWithPassword := userWithPasswordRegex.FindStringSubmatch(databaseInfo)[1]
	if len(userWithPassword) == 0 {
		return nil, fmt.Errorf("dokku: Could not find username and password for database: %s", databaseName)
	}

	port := portRegex.FindStringSubmatch(databaseInfo)[1]
	if len(port) == 0 {
		return nil, fmt.Errorf("dokku: Could not find port for database: %s", databaseName)
	}

	return sql.Open("mysql", fmt.Sprintf("%s@tcp(dokku.me:%s)/%s", userWithPassword, port, databaseName))
}