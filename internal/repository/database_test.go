package repository

import (
	"context"
	"testing"
	"time"

	"github.com/RomanAgaltsev/urlcut/internal/database/queries"
	"github.com/RomanAgaltsev/urlcut/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBRepository(t *testing.T) {
	const (
		longURL = "https://app.pachca.com"
		BaseURL = "http://localhost:8080"
		urlID   = "1q2w3e4r"
	)

	urlS := &model.URL{
		Long:   longURL,
		Base:   BaseURL,
		ID:     urlID,
		CorrID: "",
	}

	var dbRepository *DBRepository

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	rowsIns := sqlmock.NewRows([]string{"id", "long_url", "base_url", "url_id", "created_at"}).
		AddRow(1, longURL, BaseURL, urlID, time.Now())
	rowsSel := sqlmock.NewRows([]string{"id", "long_url", "base_url", "url_id", "created_at"}).
		AddRow(1, longURL, BaseURL, urlID, time.Now())

	mock.ExpectBegin()
	mock.ExpectQuery("(.*)INSERT(.*)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rowsIns)
	mock.ExpectCommit()
	mock.ExpectQuery("(.*)SELECT(.*)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rowsSel)
	mock.ExpectClose()

	var q *queries.Queries

	q, err = queries.Prepare(context.Background(), db)
	if err != nil {
		q = queries.New(db)
	}

	dbRepository = &DBRepository{
		db: db,
		q:  q,
	}

	_, err = dbRepository.Store(context.TODO(), []*model.URL{urlS})
	require.NoError(t, err)

	urlG, err := dbRepository.Get(context.TODO(), urlID)
	require.NoError(t, err)

	assert.Equal(t, urlS, urlG)

	err = dbRepository.Close()
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
