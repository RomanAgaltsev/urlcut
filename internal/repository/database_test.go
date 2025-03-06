package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/pkg/random"
)

func TestDBRepository(t *testing.T) {
	const (
		BaseURL = "http://localhost:8080"
	)
	longURL := fmt.Sprintf("https://%s.%s", random.String(20), random.String(3))
	urlID := random.String(8)

	uid, _ := uuid.NewRandom()

	urlS := &model.URL{
		Long:   longURL,
		Base:   BaseURL,
		ID:     urlID,
		CorrID: "",
		UID:    uid,
	}

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	rowsIns := sqlmock.NewRows([]string{"id", "long_url", "base_url", "url_id", "created_at", "uid", "is_deleted"}).
		AddRow(1, longURL, BaseURL, urlID, time.Now(), uid, false)
	rowsSel := sqlmock.NewRows([]string{"id", "long_url", "base_url", "url_id", "created_at", "uid", "is_deleted"}).
		AddRow(1, longURL, BaseURL, urlID, time.Now(), uid, false)
	rowsStats := sqlmock.NewRows([]string{"urls", "users"})

	mock.ExpectPrepare("(.*)UPDATE(.*)")
	mock.ExpectPrepare("(.*)SELECT(.*)")
	mock.ExpectPrepare("(.*)SELECT(.*)")
	mock.ExpectPrepare("(.*)SELECT(.*)")
	mock.ExpectBegin()
	mock.ExpectQuery("(.*)INSERT(.*)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rowsIns)
	mock.ExpectCommit()
	mock.ExpectQuery("(.*)SELECT(.*)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rowsSel)
	mock.ExpectQuery("(.*)SELECT(.*)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rowsSel)
	mock.ExpectBegin()
	mock.ExpectExec("(.*)UPDATE(.*)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("(.*)SELECT(.*)").
		WithArgs().
		WillReturnRows(rowsStats)
	mock.ExpectClose()

	dbRepository, err := NewDBRepository(db)
	require.NoError(t, err)

	storedURLs, err := dbRepository.Store(context.TODO(), []*model.URL{urlS})
	require.NoError(t, err)
	assert.IsType(t, urlS, storedURLs)

	urlG, err := dbRepository.Get(context.TODO(), urlID)
	require.NoError(t, err)

	assert.Equal(t, urlS.Long, urlG.Long)
	assert.Equal(t, urlS.Base, urlG.Base)
	assert.Equal(t, urlS.ID, urlG.ID)
	assert.Equal(t, urlS.CorrID, urlG.CorrID)

	userURLs, err := dbRepository.GetUserURLs(context.TODO(), uid)
	require.NoError(t, err)
	assert.IsType(t, []*model.URL{urlS}, userURLs)

	err = dbRepository.DeleteURLs(context.TODO(), []*model.URL{urlS})
	require.NoError(t, err)

	stats, err := dbRepository.GetStats(context.TODO())
	assert.Equal(t, err, sql.ErrNoRows)
	assert.Nil(t, stats)

	err = dbRepository.Close()
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
