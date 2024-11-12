package repositories

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"errors"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
)

type MusicRepo struct {
	db *sqlx.DB
}

func NewMusic(db *sqlx.DB) *MusicRepo {
	return &MusicRepo{db}
}

func (repo *MusicRepo) InsertData(data *models.Music) (int64, error) {
	q := `INSERT INTO stream.music (artis_id, slug, title, release_date, cover, source_url)
	VALUES(:artis_id, :slug, :title, :release_date, :cover, :source_url);`

	res, err := repo.db.NamedExec(q, data)
	if err != nil {
		return 0, errors.New(config.BadData)
	}

	return res.RowsAffected()
}

func (repo *MusicRepo) UpdateData(data *models.Music) (int64, error) {
	q := `
	UPDATE stream.music SET 
		artis_id=COALESCE(NULLIF(:artis_id, ''), artis_id),
		slug=COALESCE(NULLIF(:slug, ''), slug),
		title=COALESCE(NULLIF(:title, ''), title),
		release_date=COALESCE(NULLIF(:release_date, ''), release_date),
		cover=COALESCE(NULLIF(:cover, ''), cover),
		source_url=COALESCE(NULLIF(:source_url, ''), source_url),
		updated_at=now()
	WHERE music_id = :music_id;
	`

	res, err := repo.db.NamedExec(q, data)
	if err != nil {
		return 0, errors.New(config.BadData)
	}

	return res.RowsAffected()
}

func (repo *MusicRepo) DeleteData(uid string) (int64, error) {
	q := `DELETE FROM stream.music WHERE music_id = $1;`
	res := repo.db.MustExec(q, uid)
	return res.RowsAffected()
}

func (repo *MusicRepo) GetDataById(uid string) (*models.Music, error) {
	q := `SELECT music_id, artis_id, slug, title, release_date, cover, source_url, created_at, updated_at
	FROM stream.music WHERE music_id = $1`

	var data = new(models.Music)
	if err := repo.db.Get(data, q, uid); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return data, nil
}

func (repo *MusicRepo) GetDataBySlug(slug string) (*models.Music, error) {
	q := `SELECT music_id, artis_id, slug, title, release_date, cover, source_url, created_at, updated_at
	FROM stream.music WHERE slug = $1`

	var data = new(models.Music)
	if err := repo.db.Get(data, q, slug); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return data, nil
}

func (repo *MusicRepo) GetAllData(params *config.Pagination) (*config.ResultWarp, error) {
	var data = new(models.Musics)
	var metaResult = new(config.Meta)
	var filterQuery string
	var metaQuery string
	// var orderQuery string

	filterConditions := []config.FilterParams{
		{Param: params.Title, Column: "title"},
		{Param: params.Release, Column: "release_date"},
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

	countQuery := fmt.Sprintf(`SELECT COUNT(music_id) as "count" FROM stream.music WHERE true %s`, filterQuery)
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

	q := fmt.Sprintf(`SELECT music_id, artis_id, slug, title, release_date, cover, source_url, created_at, updated_at
	FROM stream.music WHERE true %s ORDER BY created_at DESC %s `, filterQuery, metaQuery)
	if err := repo.db.Select(data, repo.db.Rebind(q)); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return &config.ResultWarp{Data: data, Meta: metaResult}, nil
}
