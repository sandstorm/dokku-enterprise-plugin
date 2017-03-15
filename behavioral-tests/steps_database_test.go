package main

import (
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/sandstorm/dokku-enterprise-plugin/behavioral-tests/dokkuDatabaseHelper"
	"fmt"
	"strings"
)

func iExecuteTheFollowingSQLStatementsOnDatabase(databaseName string, queryString *gherkin.DocString) error {
	queries := strings.Split(queryString.Content, ";")

	for _, query := range queries {
		if len(strings.TrimSpace(query)) > 0 {
			_, err := dokkuDatabaseHelper.Execute(databaseName, query)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func theSQLStatementOnDatabaseMustReturn(query, databaseName, result string) error {
	queryResult, err := dokkuDatabaseHelper.QueryString(databaseName, query)
	if err != nil {
		return err
	}

	if queryResult != result {
		return fmt.Errorf("Expected query result is not equal to actual result: %v", queryResult)
	}

	return nil
}