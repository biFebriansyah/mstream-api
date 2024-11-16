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

type MusicRepo struct {
	db *sqlx.DB
}

func NewMusic(db *sqlx.DB) *MusicRepo {
	return &MusicRepo{db}
}

func (repo *MusicRepo) InsertData(data *models.Music) (string, error) {
	var uid string
	tx := repo.db.MustBegin()
	stm, _ := tx.Preparex(`INSERT INTO stream.music (artis_id, slug, title, release_date, cover, source_url)
	VALUES($1, $2, $3, $4, $5, $6) returning music_id`)

	err := stm.Get(&uid, data.Artis_id, data.Slug, data.Title, data.Release_date, data.Cover, data.Source_url)
	if err != nil {
		if errrb := tx.Rollback(); errrb != nil {
			return uid, errrb
		}
		return uid, err
	}

	q := `INSERT INTO stream.music_genre (music_id, genre_id) VALUES(:music_id, :genre_id) ON CONFLICT ON CONSTRAINT music_genre_musics_unique DO NOTHING;`
	for _, v := range data.Genre {
		genreMusicData := &models.GenreMusic{Music_id: &uid, Genre_id: &v}
		_, err := tx.NamedExec(q, genreMusicData)
		if err != nil {
			if errrb := tx.Rollback(); errrb != nil {
				return uid, errrb
			}
			return uid, err
		}
	}

	if err := tx.Commit(); err != nil {
		return uid, err
	}

	return uid, nil
}

func (repo *MusicRepo) UpdateData(data *models.Music) (int64, error) {
	tx := repo.db.MustBegin()
	q1 := `
	UPDATE stream.music SET 
		artis_id=COALESCE(NULLIF(:artis_id, artis_id), artis_id),
		slug=COALESCE(NULLIF(:slug, ''), slug),
		title=COALESCE(NULLIF(:title, ''), title),
		release_date=COALESCE(NULLIF(:release_date, release_date), release_date),
		cover=COALESCE(NULLIF(:cover, ''), cover),
		source_url=COALESCE(NULLIF(:source_url, ''), source_url),
		updated_at=now()
	WHERE music_id = :music_id;
	`

	res, err := tx.NamedExec(q1, data)
	if err != nil {
		if errrb := tx.Rollback(); errrb != nil {
			return 0, errrb
		}
		return 0, err
	}

	q2 := `UPDATE stream.music_genre SET genre_id=COALESCE(NULLIF(:genre_id, genre_id), genre_id) WHERE music_id = :music_id`
	if len(data.Genre) > 0 {
		for _, v := range data.Genre {
			genreMusicData := &models.GenreMusic{Music_id: &data.Music_id, Genre_id: &v}
			_, err := tx.NamedExec(q2, genreMusicData)
			if err != nil {
				if errrb := tx.Rollback(); errrb != nil {
					return 0, errrb
				}
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (repo *MusicRepo) DeleteData(uid string) (int64, error) {
	q := `DELETE FROM stream.music WHERE music_id = $1;`
	res := repo.db.MustExec(q, uid)
	return res.RowsAffected()
}

func (repo *MusicRepo) GetDataById(uid string) (*models.Music, error) {
	q2 := `
	select 
		m.music_id,
		m.title,
		m.slug,
		a."name" as artis,
		m.release_date,
		m.cover,
		m.source_url,
		ARRAY_AGG(gn.genre_name) AS genres,
		m.created_at,
		m.updated_at
	from stream.music m
	join stream.artis a on m.artis_id = a.artis_id
	join stream.music_genre mg on m.music_id = mg.music_id
	join stream.genre gn on mg.genre_id = gn.genre_id
	where m.music_id = $1
	group by m.music_id, a."name"
	`

	var data = new(models.Music)
	if err := repo.db.Get(data, q2, uid); err != nil {
		log.Println(err)
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return data, nil
}

func (repo *MusicRepo) GetDataBySlug(slug string) (*models.Music, error) {
	q := `
	select 
		m.music_id,
		m.title,
		m.,
		a."name" as artis,
		m.release_date,
		m.cover,
		m.source_url,
		ARRAY_AGG(gn.genre_name) AS genres,
		m.created_at,
		m.updated_at
	from stream.music m
	join stream.artis a on m.artis_id = a.artis_id
	join stream.music_genre mg on m.music_id = mg.music_id
	join stream.genre gn on mg.genre_id = gn.genre_id
	where m.slug = $1
	group by m.music_id, a."name"
	`

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
		cek := int32(math.Ceil(float64(metaResult.Total) / float64(params.Limit)))
		if params.Page >= cek {
			metaResult.Next = 0
		} else {
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
