package sqlutil

import (
	"fmt"
	"os"
)

type TestDB struct {
	DriverName string
	ConnStr    string
	Host       string
	Port       string
	Database   string
	Cleanup    func()
}

func GetTestDBType() string {
	dbType := "sqlite3"

	// environment variable present for test db?
	if db, present := os.LookupEnv("GRAFANA_TEST_DB"); present {
		dbType = db
	}
	return dbType
}

func GetTestDB(dbType string) (*TestDB, error) {
	switch dbType {
	case "mysql":
		db := MySQLTestDB()
		return &db, nil
	case "postgres":
		db := PostgresTestDB()
		return &db, nil
	case "sqlite3":
		return sqLite3TestDB()
	case "ydb":
		return ydbTestDB()
	}

	return nil, fmt.Errorf("unknown test db type: %s", dbType)
}

func sqLite3TestDB() (*TestDB, error) {
	if os.Getenv("SQLITE_INMEMORY") == "true" {
		return &TestDB{
			DriverName: "sqlite3",
			ConnStr:    "file::memory:",
			Cleanup:    func() {},
		}, nil
	}

	ret := &TestDB{
		DriverName: "sqlite3",
		// ConnStr specifies an In-memory database shared between connections.
		ConnStr: "file::memory:?cache=shared",
	}
	return ret, nil
}

func SQLite3TestDB() TestDB {
	// To run all tests in a local test database, set ConnStr to "grafana_test.db"
	return TestDB{
		DriverName: "sqlite3",
		// ConnStr specifies an In-memory database shared between connections.
		ConnStr: "file::memory:?cache=shared",
	}
}

func MySQLTestDB() TestDB {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	conn_str := fmt.Sprintf("grafana:password@tcp(%s:%s)/grafana_tests?collation=utf8mb4_unicode_ci&sql_mode='ANSI_QUOTES'&parseTime=true", host, port)
	return TestDB{
		DriverName: "mysql",
		ConnStr:    conn_str,
	}
}

func PostgresTestDB() TestDB {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	connStr := fmt.Sprintf("user=grafanatest password=grafanatest host=%s port=%s dbname=grafanatest sslmode=disable",
		host, port)
	return TestDB{
		DriverName: "postgres",
		ConnStr:    connStr,
	}
}

func MSSQLTestDB() TestDB {
	host := os.Getenv("MSSQL_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MSSQL_PORT")
	if port == "" {
		port = "1433"
	}
	return TestDB{
		DriverName: "mssql",
		ConnStr:    fmt.Sprintf("server=%s;port=%s;database=grafanatest;user id=grafana;password=Password!", host, port),
	}
}

func ydbTestDB() (*TestDB, error) {
	return &TestDB{
		DriverName: "ydb",
		ConnStr:    "grpc://127.0.0.1:2136/local?go_query_mode=query&go_fake_tx=query&go_query_bind=numeric",
		Host:       "127.0.0.1",
		Port:       "2136",
		Database:   "/local",
		Cleanup:    func() {},
	}, nil
}
