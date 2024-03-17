package database_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/andremfp/snippetbox/internal/database"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSnippetModel(t *testing.T) {
	t.Run("insert snippet successfully, returning id", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		stmt := regexp.QuoteMeta("INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))")

		mock.ExpectExec(stmt).WithArgs("title", "content", 7).WillReturnResult(sqlmock.NewResult(1, 0))

		gotID, _ := testSnippetStore.Insert("title", "content", 7)
		wantID := 1
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		if gotID != wantID {
			t.Errorf("got id %d, want %d", gotID, wantID)
		}

	})
}

func setDbMock(t testing.TB) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not open connection to mock db, %v", err)
	}

	return db, mock
}
