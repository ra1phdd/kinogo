package services

import (
	"encoding/json"
	"fmt"
	"kinogo/internal/models"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	"time"

	"go.uber.org/zap"
)

func GetCommentsFromDB(id int) ([]models.Comment, error) {
	var commentsSlice []models.Comment
	comments, err := cache.Rdb.Get(cache.Ctx, "QueryComments_"+fmt.Sprint(id)).Result()
	if err == nil && comments != "" {
		// Данные найдены в Redis, возвращаем их
		err = json.Unmarshal([]byte(comments), &commentsSlice)
		if err != nil {
			return nil, err
		}
		return commentsSlice, nil
	}

	// Данных нет в Redis, получаем их из базы данных
	rows, err := db.Conn.Queryx(`SELECT * FROM comments WHERE movieid = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logger.Debug("Получены комментарии из БД", zap.Any("rows", rows))

	for rows.Next() {
		var CommentsData models.Comment
		err := rows.Scan(
			&CommentsData.ID,
			&CommentsData.Text,
			&CommentsData.ParentID,
			&CommentsData.User.ID,
			&CommentsData.MovieID,
		)
		if err != nil {
			return nil, err
		}

		rowsUsers, err := db.Conn.Queryx(`SELECT * FROM users WHERE id = $1`, &CommentsData.User.ID)
		if err != nil {
			return nil, err
		}
		defer rowsUsers.Close()

		for rowsUsers.Next() {
			var lastName *string
			err := rowsUsers.Scan(
				&CommentsData.User.ID,
				&CommentsData.User.FirstName,
				&lastName,
				&CommentsData.User.Username,
				&CommentsData.User.PhotoURL,
			)
			if err != nil {
				return nil, err
			}
			if lastName != nil {
				CommentsData.User.LastName = *lastName
			}
		}

		commentsSlice = append(commentsSlice, CommentsData)
	}

	// Сохраняем данные в Redis
	commentsJSON, err := json.Marshal(commentsSlice)
	if err != nil {
		return nil, err
	}
	err = cache.Rdb.Set(cache.Ctx, "QueryComments_"+fmt.Sprint(id), commentsJSON, 15*time.Second).Err()
	if err != nil {
		return nil, err
	}

	return commentsSlice, nil
}

func BuildCommentTree(comments []models.Comment) []models.Comment {
	commentMap := make(map[int]*models.Comment)

	// Создаем карту комментариев
	for i := range comments {
		comment := &comments[i]
		commentMap[comment.ID] = comment
	}

	// Строим древовидную структуру
	for i := range comments {
		comment := &comments[i]
		parentComment, exists := commentMap[comment.ParentID]
		if !exists {
			// Корневой комментарий
			continue
		} else {
			parentComment.Children = append(parentComment.Children, comment)
		}
	}

	// Извлекаем корневые комментарии из карты
	var allComm []models.Comment
	for _, comment := range commentMap {
		if comment.ParentID == 0 {
			allComm = append(allComm, *comment)
		}
	}

	return allComm
}
