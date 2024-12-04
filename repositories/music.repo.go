package repositories

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"encoding/json"
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

func (repo *MusicRepo) InsertData(data *models.MusicData) (string, error) {
	var uid string
	tx := repo.db.MustBegin()
	stm, _ := tx.Preparex(`INSERT INTO stream.music (slug, title, release_date, cover, source_url)
	VALUES($1, $2, $3, $4, $5) returning music_id`)

	err := stm.Get(&uid, data.Slug, data.Title, data.Release_date, data.Cover, data.Source_url)
	if err != nil {
		if errrb := tx.Rollback(); errrb != nil {
			return uid, errrb
		}
		return uid, err
	}

	q1 := `INSERT INTO stream.music_artis (music_id, artis_id) VALUES(:music_id, :artis_id) ON CONFLICT ON CONSTRAINT music_artis_unique1 DO NOTHING;`
	artisMusicData := &models.MusicArtis{Music_id: &uid, Artis_id: data.MusicArtis.Artis_id}
	_, err = tx.NamedExec(q1, artisMusicData)
	if err != nil {
		if errrb := tx.Rollback(); errrb != nil {
			return uid, errrb
		}
	}

	q2 := `INSERT INTO stream.music_genre (music_id, genre_id) VALUES(:music_id, :genre_id) ON CONFLICT ON CONSTRAINT music_genre_unique DO NOTHING;`
	for _, v := range data.MusicGenre {
		genreMusicData := &models.GenreMusic{Music_id: &uid, Genre_id: v.Genre_id}
		_, err := tx.NamedExec(q2, genreMusicData)
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

func (repo *MusicRepo) InsertSource(url, uid string) (int64, error) {
	q := `UPDATE stream.music SET source_url = $1, updated_at = now() WHERE music_id = $2`
	res := repo.db.MustExec(q, url, uid)
	return res.RowsAffected()
}

func (repo *MusicRepo) UpdateData(data *models.Music) (int64, error) {
	tx := repo.db.MustBegin()
	q1 := `
	UPDATE stream.music SET 
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

func (repo *MusicRepo) GetDataById(uid string) (*models.MusicData, error) {
	q1 := `
	SELECT 
		m.music_id,
		m.title,
		m.slug,
		(SELECT 
			DISTINCT JSONB_BUILD_OBJECT(
				'artis_id', a.artis_id,
				'artis_name', CONCAT(a.first_name, ' ' , a.last_name) 
			)
			FROM stream.artis a
			JOIN stream.music_artis ma ON a.artis_id = ma.artis_id
			WHERE ma.music_id = m.music_id 
		) AS music_artis,
		(SELECT 
			JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
				'genre_id', g.genre_id,
				'genre_name', g.genre_name
			))
			FROM stream.genre g
			JOIN stream.music_genre mg ON g.genre_id = mg.genre_id
			WHERE mg.music_id = m.music_id 
		) AS music_genre,
		m.release_date,
		m.cover,
		m.source_url,
		m.created_at,
		m.updated_at
	FROM stream.music m
	WHERE m.music_id = $1
	GROUP BY m.music_id`

	var data = new(models.MusicData)
	var musicGenreJSON, musicArtisJSON []byte
	if rows, err := repo.db.Queryx(q1, uid); err == nil {
		for rows.Next() {
			err := rows.Scan(
				&data.Music_id,
				&data.Title,
				&data.Slug,
				&musicArtisJSON,
				&musicGenreJSON,
				&data.Release_date,
				&data.Cover,
				&data.Source_url,
				&data.CreatedAt,
				&data.UpdateAt,
			)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			if len(musicArtisJSON) > 0 {
				if err := json.Unmarshal(musicArtisJSON, &data.MusicArtis); err != nil {
					return nil, fmt.Errorf("failed to unmarshal music_artis: %w", err)
				}
			}

			if len(musicGenreJSON) > 0 {
				if err := json.Unmarshal(musicGenreJSON, &data.MusicGenre); err != nil {
					return nil, fmt.Errorf("failed to unmarshal music_genre: %w", err)
				}
			}
		}
	}

	// if err := repo.db.Get(data, q1, uid); err != nil {
	// 	if err.Error() == "sql: no rows in result set" {
	// 		return nil, errors.New(config.NotFound)
	// 	}
	// 	return nil, err
	// }

	return data, nil
}

func (repo *MusicRepo) GetDataBySlug(slug string) (*models.MusicData, error) {
	q1 := `
	SELECT 
		m.music_id,
		m.title,
		m.slug,
		(SELECT 
			DISTINCT JSONB_BUILD_OBJECT(
				'artis_id', a.artis_id,
				'artis_name', CONCAT(a.first_name, ' ' , a.last_name) 
			)
			FROM stream.artis a
			JOIN stream.music_artis ma ON a.artis_id = ma.artis_id
			WHERE ma.music_id = m.music_id 
		) AS music_artis,
		(SELECT 
			JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
				'genre_id', g.genre_id,
				'genre_name', g.genre_name
			))
			FROM stream.genre g
			JOIN stream.music_genre mg ON g.genre_id = mg.genre_id
			WHERE mg.music_id = m.music_id 
		) AS music_genre,
		m.release_date,
		m.cover,
		m.source_url,
		m.created_at,
		m.updated_at
	FROM stream.music m
	WHERE m.slug = $1
	GROUP BY m.music_id
	`

	var data = new(models.MusicData)
	var musicGenreJSON, musicArtisJSON []byte
	if rows, err := repo.db.Queryx(q1, slug); err == nil {
		for rows.Next() {
			err := rows.Scan(
				&data.Music_id,
				&data.Title,
				&data.Slug,
				&musicArtisJSON,
				&musicGenreJSON,
				&data.Release_date,
				&data.Cover,
				&data.Source_url,
				&data.CreatedAt,
				&data.UpdateAt,
			)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			if len(musicArtisJSON) > 0 {
				if err := json.Unmarshal(musicArtisJSON, &data.MusicArtis); err != nil {
					return nil, fmt.Errorf("failed to unmarshal music_artis: %w", err)
				}
			}

			if len(musicGenreJSON) > 0 {
				if err := json.Unmarshal(musicGenreJSON, &data.MusicGenre); err != nil {
					return nil, fmt.Errorf("failed to unmarshal music_genre: %w", err)
				}
			}
		}
	}

	return data, nil
}

func (repo *MusicRepo) GetAllData(params *config.Pagination) (*config.ResultWarp, error) {
	var datas = new([]models.MusicData)
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

	// q := fmt.Sprintf(`SELECT music_id, slug, title, release_date, cover, source_url, created_at, updated_at
	// FROM stream.music WHERE true %s ORDER BY created_at DESC %s `, filterQuery, metaQuery)

	q := fmt.Sprintf(`
	SELECT 
		m.music_id,
		m.title,
		m.slug,
		(SELECT 
			DISTINCT JSONB_BUILD_OBJECT(
				'artis_id', a.artis_id,
				'artis_name', CONCAT(a.first_name, ' ' , a.last_name) 
			)
			FROM stream.artis a
			JOIN stream.music_artis ma ON a.artis_id = ma.artis_id
			WHERE ma.music_id = m.music_id 
		) AS music_artis,
		(SELECT 
			JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT(
				'genre_id', g.genre_id,
				'genre_name', g.genre_name
			))
			FROM stream.genre g
			JOIN stream.music_genre mg ON g.genre_id = mg.genre_id
			WHERE mg.music_id = m.music_id 
		) AS music_genre,
		m.release_date,
		m.cover,
		m.source_url,
		m.created_at,
		m.updated_at
	FROM stream.music m
	WHERE true %s
	GROUP BY m.music_id
	ORDER BY created_at DESC %s `, filterQuery, metaQuery)

	if rows, err := repo.db.Queryx(q); err == nil {
		for rows.Next() {
			var data = new(models.MusicData)
			var musicGenreJSON, musicArtisJSON []byte
			err := rows.Scan(
				&data.Music_id,
				&data.Title,
				&data.Slug,
				&musicArtisJSON,
				&musicGenreJSON,
				&data.Release_date,
				&data.Cover,
				&data.Source_url,
				&data.CreatedAt,
				&data.UpdateAt,
			)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			if len(musicArtisJSON) > 0 {
				if err := json.Unmarshal(musicArtisJSON, &data.MusicArtis); err != nil {
					return nil, fmt.Errorf("failed to unmarshal music_artis: %w", err)
				}
			}

			if len(musicGenreJSON) > 0 {
				if err := json.Unmarshal(musicGenreJSON, &data.MusicGenre); err != nil {
					return nil, fmt.Errorf("failed to unmarshal music_genre: %w", err)
				}
			}

			*datas = append(*datas, *data)
		}
	}

	// if err := repo.db.Select(data, repo.db.Rebind(q)); err != nil {
	// 	fmt.Println(err)
	// 	if err.Error() == "sql: no rows in result set" {
	// 		return nil, errors.New(config.NotFound)
	// 	}
	// 	return nil, err
	// }

	return &config.ResultWarp{Data: datas, Meta: metaResult}, nil
}
