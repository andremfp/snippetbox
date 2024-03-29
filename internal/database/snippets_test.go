package database_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

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

	t.Run("get snippet successfully", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		createdDate := time.Now().AddDate(0, 0, -1)
		expiresDate := time.Now().AddDate(0, 0, +1)

		wantSnippet := &database.Snippet{
			ID:      1,
			Title:   "title",
			Content: "content",
			Created: createdDate,
			Expires: expiresDate,
		}

		mokedDbResponse := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"}).AddRow(1, "title", "content", createdDate, expiresDate)

		stmt := regexp.QuoteMeta("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?")

		mock.ExpectQuery(stmt).WithArgs(1).WillReturnRows(mokedDbResponse)

		gotSnippet, _ := testSnippetStore.Get(1)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		assertSnippet(t, gotSnippet, wantSnippet)

	})

	t.Run("snippet not found", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		stmt := regexp.QuoteMeta("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?")

		mock.ExpectQuery(stmt).WithArgs(10).WillReturnError(sql.ErrNoRows)

		_, getErr := testSnippetStore.Get(10)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		if !errors.Is(getErr, database.ErrNoRecord) {
			t.Errorf("got error %v, want %v", getErr, database.ErrNoRecord)
		}

	})

	t.Run("get snippet generic error", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		stmt := regexp.QuoteMeta("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?")

		mock.ExpectQuery(stmt).WithArgs(1).WillReturnError(errors.New("generic error"))

		_, getErr := testSnippetStore.Get(1)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		if getErr.Error() != "generic error" {
			t.Errorf("got error %s, want %s", getErr.Error(), "generic error")
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

func assertSnippet(t testing.TB, got, want *database.Snippet) {
	t.Helper()
	if got.ID != want.ID || got.Content != want.Content || got.Title != want.Title || got.Created != want.Created || got.Expires != want.Expires {
		t.Errorf("got snippet %v, want %v", got, want)
	}
}
