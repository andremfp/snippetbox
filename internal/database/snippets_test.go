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

	t.Run("insert snippet error", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		stmt := regexp.QuoteMeta("INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))")

		mock.ExpectExec(stmt).WithArgs("title", "content", 7).WillReturnError(database.ErrGeneric)

		gotID, gotErr := testSnippetStore.Insert("title", "content", 7)
		wantID := 0
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		if gotID != wantID {
			t.Errorf("got id %d, want %d", gotID, wantID)
		}
		if !errors.Is(gotErr, database.ErrGeneric) {
			t.Errorf("got error %v, want %v", gotErr, database.ErrGeneric)
		}

	})

	t.Run("insert snippet result error", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		stmt := regexp.QuoteMeta("INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))")

		mock.ExpectExec(stmt).WithArgs("title", "content", 7).WillReturnResult(sqlmock.NewErrorResult(database.ErrGeneric))

		gotID, gotErr := testSnippetStore.Insert("title", "content", 7)
		wantID := 0
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		if gotID != wantID {
			t.Errorf("got id %d, want %d", gotID, wantID)
		}
		if !errors.Is(gotErr, database.ErrGeneric) {
			t.Errorf("got error %v, want %v", gotErr, database.ErrGeneric)
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

		mock.ExpectQuery(stmt).WithArgs(1).WillReturnError(database.ErrGeneric)

		_, gotErr := testSnippetStore.Get(1)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		if !errors.Is(gotErr, database.ErrGeneric) {
			t.Errorf("got error %v, want %v", gotErr, database.ErrGeneric)
		}

	})

	t.Run("get latest snippets successfully", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		createdDate := time.Now().AddDate(0, 0, -1)

		expiredDate := time.Now().AddDate(0, 0, +1)

		wantSnippets := []*database.Snippet{
			{ID: 1, Title: "title1", Content: "content1", Created: createdDate, Expires: expiredDate},
			{ID: 2, Title: "title2", Content: "content2", Created: createdDate, Expires: expiredDate},
			{ID: 3, Title: "title3", Content: "content3", Created: createdDate, Expires: expiredDate},
			{ID: 4, Title: "title4", Content: "content4", Created: createdDate, Expires: expiredDate},
			{ID: 5, Title: "title5", Content: "content5", Created: createdDate, Expires: expiredDate},
			{ID: 6, Title: "title6", Content: "content6", Created: createdDate, Expires: expiredDate},
			{ID: 7, Title: "title7", Content: "content7", Created: createdDate, Expires: expiredDate},
			{ID: 8, Title: "title8", Content: "content8", Created: createdDate, Expires: expiredDate},
			{ID: 9, Title: "title9", Content: "content9", Created: createdDate, Expires: expiredDate},
			{ID: 10, Title: "title10", Content: "content10", Created: createdDate, Expires: expiredDate},
		}

		mokedDbResponse := sqlmock.NewRows([]string{"id", "title", "content", "created", "expires"}).
			AddRow(1, "title1", "content1", createdDate, expiredDate).
			AddRow(2, "title2", "content2", createdDate, expiredDate).
			AddRow(3, "title3", "content3", createdDate, expiredDate).
			AddRow(4, "title4", "content4", createdDate, expiredDate).
			AddRow(5, "title5", "content5", createdDate, expiredDate).
			AddRow(6, "title6", "content6", createdDate, expiredDate).
			AddRow(7, "title7", "content7", createdDate, expiredDate).
			AddRow(8, "title8", "content8", createdDate, expiredDate).
			AddRow(9, "title9", "content9", createdDate, expiredDate).
			AddRow(10, "title10", "content10", createdDate, expiredDate)

		stmt := regexp.QuoteMeta("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10")

		mock.ExpectQuery(stmt).WillReturnRows(mokedDbResponse)

		gotSnippets, _ := testSnippetStore.Latest()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		assertSnippetList(t, gotSnippets, wantSnippets)

	})

	t.Run("get latest snippets generic error", func(t *testing.T) {
		db, mock := setDbMock(t)
		defer db.Close()
		testSnippetStore := database.SnippetModel{DB: db}

		stmt := regexp.QuoteMeta("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10")

		mock.ExpectQuery(stmt).WillReturnError(database.ErrGeneric)

		_, gotErr := testSnippetStore.Latest()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected sql statement not met, %v", err)
		}

		if !errors.Is(gotErr, database.ErrGeneric) {
			t.Errorf("got error %v, want %v", gotErr, database.ErrGeneric)
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

func assertSnippetList(t testing.TB, got, want []*database.Snippet) {
	t.Helper()
	for i := 0; i < len(got); i++ {
		assertSnippet(t, got[i], want[i])
	}
}
