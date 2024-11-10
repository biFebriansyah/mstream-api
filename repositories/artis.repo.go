package repositories

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"errors"

	"github.com/jmoiron/sqlx"
)

type ArtiRepo struct {
	db *sqlx.DB
}

func NewArtis(db *sqlx.DB) *ArtiRepo {
	return &ArtiRepo{db}
}

func (repo *ArtiRepo) InsertData(data *models.Artis) (int64, error) {
	q := `INSERT INTO stream.artis ("name", slug, nationality) VALUES(:name, :slug, :nationality)`

	res, err := repo.db.NamedExec(q, data)
	if err != nil {
		return 0, errors.New(config.BadData)
	}

	return res.RowsAffected()
}

func (repo *ArtiRepo) UpdateData(data *models.Artis) (int64, error) {
	q := `
	UPDATE stream.artis SET 
		name=COALESCE(NULLIF(:name, ''), name),
		slug=COALESCE(NULLIF(:slug, ''), slug),
		nationality=COALESCE(NULLIF(:nationality, ''), nationality),
		updated_at=now()
	WHERE artis_id=artis_id;
	`

	res, err := repo.db.NamedExec(q, data)
	if err != nil {
		return 0, errors.New(config.BadData)
	}

	return res.RowsAffected()
}

func (repo *ArtiRepo) DeleteData(uid string) (int64, error) {
	q := `DELETE FROM stream.artis WHERE artis_id = $1;`
	res := repo.db.MustExec(q, uid)
	return res.RowsAffected()
}

func (repo *ArtiRepo) GetAllData() (*models.Arties, error) {
	q := `SELECT artis_id, "name", slug, nationality, created_at, updated_at
	FROM stream.artis ORDER BY created_at DESC`

	var data models.Arties
	if err := repo.db.Select(&data, q); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return &data, nil
}

func (repo *ArtiRepo) GetDataById(uid string) (*models.Artis, error) {
	q := `SELECT artis_id, "name", slug, nationality, created_at, updated_at
	FROM stream.artis WHERE artis_id = $1`

	var data models.Artis
	if err := repo.db.Get(&data, q, uid); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return &data, nil
}

func (repo *ArtiRepo) GetDataBySlug(slug string) (*models.Artis, error) {
	q := `SELECT artis_id, "name", slug, nationality, created_at, updated_at
	FROM stream.artis WHERE slug = $1`

	var data models.Artis
	if err := repo.db.Get(&data, q, slug); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return &data, nil
}
