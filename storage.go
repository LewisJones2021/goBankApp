package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// storage is an interface defining methods for
// basic CRUD operations on Account objects.
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	getAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type postgresStore struct {
	db *sql.DB
}

// NewPostgresStore initializes and returns a new instance of the postgresStore.
// It establishes a connection to a PostgreSQL database using the provided connection string.
// Returns the initialized postgresStore and an error if any connection issues occur.
func NewPostgresStore() (*postgresStore, error) {
	// Connection string for PostgreSQL database.
	connStr := "user=postgres dbname=postgres password=password sslmode=disable"
	// Open a connection to the PostgreSQL database.
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	// Ping the database to ensure the connection is valid.
	if err := db.Ping(); err != nil {
		// Return nil and the error if there's an issue pinging the database.
		return nil, err
	}
	// Return a pointer to the initialized postgresStore with the opened database connection.
	return &postgresStore{
		db: db,
	}, nil
}

// Init initializes the PostgreSQL store by creating the account table.
func (s *postgresStore) Init() error {
	return s.creeateAccountTable()
}

// createAccountTable creates the 'account' table in the PostgreSQL database if it doesn't exist.
func (s *postgresStore) creeateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`
	// Execute the SQL query to create the 'account' table.
	_, err := s.db.Exec(query)
	return err
}

// CreateAccount creates a new account in the PostgreSQL database.
// It takes an Account struct as an input and inserts its values into the 'account'
func (s *postgresStore) CreateAccount(acc *Account) error {
	// SQL query to insert account details into the 'account' table
	query := `insert into account
	(first_name, last_name, number, balance, created_at)
	values($1, $2, $3, $4, $5)`
	// Execute the SQL query with account details as parameters
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance, acc.CreatedAt,
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	// Return nil if the operation is successful
	return nil
}

func (s *postgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *postgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *postgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

// getAccounts is a method of the postgresStore type responsible for retrieving all accounts from the 'account' table.
// It returns a slice of Account pointers and an error if the database
func (s *postgresStore) getAccounts() ([]*Account, error) {
	// Execute a SQL query to select all rows from the 'account' table
	rows, err := s.db.Query("select * from account")
	if err != nil {
		// Return an empty slice and the error if there's a problem with the query execution
		return nil, err
	}

	// Initialize an empty slice to store Account pointers
	accounts := []*Account{}
	// Iterate through the result set obtained from the query
	for rows.Next() {
		// Create a new Account instance to store the current row's data
		account := new(Account)
		// Scan the values from the current row into the Account instance
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt)

		if err != nil {
			return nil, err
		}

		// Append the current Account instance to the accounts slice
		accounts = append(accounts, account)
	}
	// Return the populated accounts slice and nil to indicate a successful operation
	return accounts, nil
}
