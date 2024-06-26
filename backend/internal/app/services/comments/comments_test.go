package comments_v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"kinogo/internal/app/models"
	"kinogo/internal/app/services/testutil"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"testing"
	"time"
)

func TestService_GetCommentsByIdService(t *testing.T) {
	conn, mockDB, mock, mr, rdb := testutil.SetupMocks()
	defer mockDB.Close()
	defer mr.Close()

	db.Conn = conn
	cache.Rdb = rdb

	s := Service{}

	query := `SELECT c.id, c.\"parentId\", c.text, c.\"createdAt\", c.\"updatedAt\", u.username, u.photourl, u.first_name, u.last_name\
			FROM comments c
			JOIN users u ON c.\"userId\" = u.id
			WHERE c.\"movieId\" = \$1 LIMIT 10 OFFSET 0`

	t.Run("Data from Redis", func(t *testing.T) {
		// Подготавливаем данные для Redis
		comments := []models.Comments{
			{
				ID:        1,
				UserID:    1,
				ParentID:  0,
				Text:      "Comment 1",
				CreatedAt: time.Now().Round(time.Second),
				UpdatedAt: time.Now().Round(time.Second),
				User: models.User{
					Username:  "user1",
					PhotoUrl:  "url1",
					FirstName: "First1",
					LastName:  "Last1",
				},
			},
			{
				ID:        2,
				UserID:    2,
				ParentID:  1,
				Text:      "Comment 2",
				CreatedAt: time.Now().Round(time.Second),
				UpdatedAt: time.Now().Round(time.Second),
				User: models.User{
					Username:  "user2",
					PhotoUrl:  "url2",
					FirstName: "First2",
					LastName:  "Last2",
				},
			},
		}
		commentsJSON, _ := json.Marshal(comments)
		err := mr.Set("comments_10_10_1", string(commentsJSON))
		if err != nil {
			t.Fatalf("error mr set")
		}

		result, err := s.GetCommentsByIdService(10, 10, 1)
		assert.NoError(t, err)
		assert.Equal(t, comments, result)
	})

	t.Run("Data from Database", func(t *testing.T) {
		mr.FlushAll()

		rows := sqlmock.NewRows([]string{"id", "parentId", "text", "createdAt", "updatedAt", "username", "photourl", "first_name", "last_name"}).
			AddRow(1, sql.NullInt32{Int32: 0, Valid: true}, "Comment 1", time.Now().Round(time.Second), time.Now().Round(time.Second), "user1", "url1", "First1", "Last1").
			AddRow(2, sql.NullInt32{Int32: 1, Valid: true}, "Comment 2", time.Now().Round(time.Second), time.Now().Round(time.Second), "user2", "url2", "First2", "Last2")

		mock.ExpectQuery(query).
			WithArgs(10).
			WillReturnRows(rows)

		result, err := s.GetCommentsByIdService(10, 10, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Comment 1", result[0].Text)
		assert.Equal(t, "Comment 2", result[0].Children[0].Text)
		assert.Equal(t, "user1", result[0].User.Username)
		assert.Equal(t, "user2", result[0].Children[0].User.Username)
	})

	t.Run("No Data", func(t *testing.T) {
		mr.FlushAll()

		mock.ExpectQuery(query).
			WithArgs(10).
			WillReturnRows(sqlmock.NewRows([]string{"id", "parentId", "text", "createdAt", "updatedAt", "username", "photourl", "first_name", "last_name"}))

		result, err := s.GetCommentsByIdService(10, 10, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "нет значений в БД")
	})

	t.Run("Database Error", func(t *testing.T) {
		mr.FlushAll()

		mock.ExpectQuery(query).
			WithArgs(10).
			WillReturnError(errors.New("database error"))

		result, err := s.GetCommentsByIdService(10, 10, 1)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("Redis Error", func(t *testing.T) {
		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})

		rows := sqlmock.NewRows([]string{"id", "parentId", "text", "createdAt", "updatedAt", "username", "photourl", "first_name", "last_name"}).
			AddRow(1, sql.NullInt32{Int32: 0, Valid: true}, "Comment 1", time.Now().Round(time.Second), time.Now().Round(time.Second), "user1", "url1", "First1", "Last1")

		mock.ExpectQuery(query).
			WithArgs(10).
			WillReturnRows(rows)

		result, err := s.GetCommentsByIdService(10, 10, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Comment 1", result[0].Text)
		assert.Equal(t, "user1", result[0].User.Username)

		cache.Rdb = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
	})
}

func TestBuildCommentTree(t *testing.T) {
	t.Run("No Comments", func(t *testing.T) {
		comments := []models.Comments{}
		allComments := map[int32]models.Comments{}
		result := buildCommentTree(comments, allComments)
		assert.Equal(t, 0, len(result))
	})

	t.Run("Single Comment", func(t *testing.T) {
		comments := []models.Comments{
			{ID: 1, ParentID: 0, Text: "Root Comment"},
		}
		allComments := map[int32]models.Comments{
			1: {ID: 1, ParentID: 0, Text: "Root Comment"},
		}
		result := buildCommentTree(comments, allComments)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, int32(1), result[0].ID)
		assert.Equal(t, 0, len(result[0].Children))
	})

	t.Run("Single Level Children", func(t *testing.T) {
		comments := []models.Comments{
			{ID: 1, ParentID: 0, Text: "Root Comment"},
		}
		allComments := map[int32]models.Comments{
			1: {ID: 1, ParentID: 0, Text: "Root Comment"},
			2: {ID: 2, ParentID: 1, Text: "Child Comment 1"},
			3: {ID: 3, ParentID: 1, Text: "Child Comment 2"},
		}
		result := buildCommentTree(comments, allComments)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, int32(1), result[0].ID)
		assert.Equal(t, 2, len(result[0].Children))
		assert.Equal(t, int32(2), result[0].Children[0].ID)
		assert.Equal(t, int32(3), result[0].Children[1].ID)
	})

	t.Run("Multiple Level Children", func(t *testing.T) {
		comments := []models.Comments{
			{ID: 1, ParentID: 0, Text: "Root Comment"},
		}
		allComments := map[int32]models.Comments{
			1: {ID: 1, ParentID: 0, Text: "Root Comment"},
			2: {ID: 2, ParentID: 1, Text: "Child Comment 1"},
			3: {ID: 3, ParentID: 1, Text: "Child Comment 2"},
			4: {ID: 4, ParentID: 2, Text: "Grandchild Comment 1"},
			5: {ID: 5, ParentID: 3, Text: "Grandchild Comment 2"},
		}
		result := buildCommentTree(comments, allComments)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, int32(1), result[0].ID)
		assert.Equal(t, 2, len(result[0].Children))
		assert.Equal(t, int32(2), result[0].Children[0].ID)
		assert.Equal(t, int32(4), result[0].Children[0].Children[0].ID)
		assert.Equal(t, int32(3), result[0].Children[1].ID)
		assert.Equal(t, int32(5), result[0].Children[1].Children[0].ID)
	})
}
