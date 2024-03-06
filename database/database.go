package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql/v2"
	goyesqlx "github.com/knadh/goyesql/v2/sqlx"
	"github.com/lib/pq"
	"github.com/zerodha/logf"
	"os"
	"reflect"
)

type DBType string

const (
	Postgres DBType = "postgres"
)

type SQLFilePaths struct {
	// QueryFilePath .sql file containing all the queries to be executed. Useful
	// for separating SQL from code logic
	QueryFilePath *string
	// SchemaFilePath .sql file containing the SQL commands. This can have SQL commands
	// to create / alter tables. Use CREATE TABLE IF NOT EXISTS for all create
	// table queries. For any subsequent modifications,
	// add alter queries inside the schema file and those changes
	// will be synced during the server boot.
	SchemaFilePath *string
}

type DBConfig struct {
	SQLFilePaths
	// Queries must be a pointer kind ( a pointer to a struct )
	Queries
	Type            DBType
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	Params          string
	db              *sqlx.DB
	l               *logf.Logger
	defaultQuerySet ThunderbyteQueries
}

// ForRoot Sets up a connection to the database, sets up the schema and initializes a repository
// struct to be used throughout the application. Panics if it fails to achieve any of the condition
func ForRoot(c *DBConfig, l *logf.Logger) {
	l.Info("connecting to db", "host", c.Host, "port", c.Port, "database", c.Database)
	db, err := sqlx.Connect(string(c.Type),
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s %s", c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode, c.Params))

	if err != nil {
		panic(err)
	}

	c.db = db
	c.l = l
	c.install()

	defaultQueryMap, _ := goyesql.ParseBytes([]byte(getDefaultRepoQueries()))
	var tbDQ ThunderbyteQueries
	goyesqlx.ScanToStruct(&tbDQ, defaultQueryMap, c.db)
	c.defaultQuerySet = tbDQ
	fmt.Println(tbDQ.CreateAuthProfile, ">>>>")
	if c.Queries != nil && reflect.TypeOf(c.Queries).Kind() == reflect.Pointer {
		if c.QueryFilePath != nil {
			queries := goyesql.MustParseFile(*c.QueryFilePath)
			err = goyesqlx.ScanToStruct(c.Queries, queries, c.db)
			if err != nil {
				c.l.Fatal("Error scanning queries to struct: ", err)
			}
		}
	}

	if c.SchemaFilePath != nil {
		if _, err := c.db.Exec(string(readQueries(*c.SchemaFilePath))); err != nil {
			c.l.Fatal("Failed while creating schema", "error", err)
		}
		c.l.Info("Applied the schema defined", "path", *c.SchemaFilePath)
	}
}

func (dbc *DBConfig) GetDB() *sqlx.DB {
	return dbc.db
}

// GetDefaultQueries returns the queries for the inbuilt repos
// within thunderbyte
func (dbc *DBConfig) GetDefaultQueries() ThunderbyteQueries {
	return dbc.defaultQuerySet
}

// Install runs the first time setup of creating and
// migrating the database
func (dbc *DBConfig) install() {
	tbd := dbc.db
	if _, err := tbd.Exec(fmt.Sprintf("select count(*) from %s", SETTINGS_REPO)); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code != "42P01" {
			panic("Error checking existing DB schema: " + err.Error())
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P01" {
			if _, err := tbd.Exec(getInitialSchemaQueries()); err != nil {
				panic("Error executing the schema file." + err.Error())
			}
		}
	}
}

// readQueries simply reads the file from the filepath
// and returns the file bytes
func readQueries(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		panic("Unable to open the queries file. Exiting..." + err.Error())
	}
	defer file.Close()
	stat, _ := file.Stat()
	queryBytes := make([]byte, stat.Size())
	_, err = file.Read(queryBytes)
	if err != nil {
		panic("Unable to read the queries file. Exiting...")
	}
	return queryBytes

}
