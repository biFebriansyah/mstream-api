package repositories

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/jmoiron/sqlx"
)

type GenreRepo struct {
	db *sqlx.DB
}

func NewGenre(db *sqlx.DB) *GenreRepo {
	return &GenreRepo{db}
}

func (repo *GenreRepo) InsertData(data *models.Genre) (int64, error) {
	q := `INSERT INTO stream.genre (name, slug) VALUES(:name, :slug)`

	res, err := repo.db.NamedExec(q, data)
	if err != nil {
		log.Println(err)
		return 0, errors.New(config.BadData)
	}

	return res.RowsAffected()
}

func (repo *GenreRepo) UpdateData(data *models.Genre) (int64, error) {
	q := `
	UPDATE stream.genre SET 
		name=COALESCE(NULLIF(:name, ''), name),
		slug=COALESCE(NULLIF(:slug, ''), slug),
		updated_at=now()
	WHERE genre_id = :genre_id;
	`

	res, err := repo.db.NamedExec(q, data)
	if err != nil {
		return 0, errors.New(config.BadData)
	}

	return res.RowsAffected()
}

func (repo *GenreRepo) DeleteData(uid string) (int64, error) {
	q := `DELETE FROM stream.genre WHERE genre_id = $1;`
	res := repo.db.MustExec(q, uid)
	return res.RowsAffected()
}

func (repo *GenreRepo) GetDataById(uid string) (*models.Genre, error) {
	q := `SELECT genre_id, name, slug, created_at, updated_at
	FROM stream.genre WHERE genre_id = $1`

	var data = new(models.Genre)
	if err := repo.db.Get(data, q, uid); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return data, nil
}

func (repo *GenreRepo) GetDataBySlug(slug string) (*models.Genre, error) {
	q := `SELECT genre_id, name, slug, created_at, updated_at
	FROM stream.genre WHERE slug = $1`

	var data = new(models.Genre)
	if err := repo.db.Get(data, q, slug); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return data, nil
}

func (repo *GenreRepo) GetAllData(params *config.Pagination) (*config.ResultWarp, error) {
	var data = new(models.Genres)
	var metaResult = new(config.Meta)
	var filterQuery string
	var metaQuery string
	// var orderQuery string

	filterConditions := []config.FilterParams{
		{Param: params.Name, Column: "name"},
	}

	for _, v := range filterConditions {
		if v.Param != "" {
			filterQuery += fmt.Sprintf(`AND %s = '%s' `, v.Column, v.Param)
		}
	}

	if params.Page != 0 && params.Limit != 0 {
		offset := (params.Page - 1) * params.Limit
		metaQuery = fmt.Sprintf("LIMIT %d OFFSET %d", params.Limit, offset)
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(genre_id) as "count" FROM stream.genre WHERE true %s`, filterQuery)
	err := repo.db.Get(&metaResult.Total, repo.db.Rebind(countQuery))
	if err != nil {
		return nil, errors.New(config.BadData)
	}

	if metaResult.Total > 0 {
		if params.Page != int32(math.Ceil(float64(metaResult.Total)/float64(params.Limit))) {
			metaResult.Next = params.Page + 1
		}
	}
	if params.Page > 1 {
		metaResult.Prev = params.Page - 1
	}

	q := fmt.Sprintf(`SELECT genre_id, name, slug, created_at, updated_at
	FROM stream.genre WHERE true %s ORDER BY created_at DESC %s `, filterQuery, metaQuery)
	if err := repo.db.Select(data, repo.db.Rebind(q)); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return &config.ResultWarp{Data: data, Meta: metaResult}, nil
}
