package repositories

import (
	"biFebriansyah/gostream/config"
	"biFebriansyah/gostream/models"
	"errors"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
)

type ArtiRepo struct {
	db *sqlx.DB
}

func NewArtis(db *sqlx.DB) *ArtiRepo {
	return &ArtiRepo{db}
}

func (repo *ArtiRepo) InsertData(data *models.Artis) (int64, error) {
	q := `INSERT INTO stream.artis (slug, first_name, last_name, email, birth_date, picture, nationality, street_address, city, province, zip_code) VALUES(:slug, :first_name, :last_name, :email, :birth_date, :picture, :nationality, :street_address, :city, :province, :zip_code)`

	res, err := repo.db.NamedExec(q, data)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New(config.BadData)
	}

	return res.RowsAffected()
}

func (repo *ArtiRepo) UpdateData(data *models.Artis) (int64, error) {
	q := `
	UPDATE stream.artis SET 
		first_name=COALESCE(NULLIF(:first_name, ''), first_name),
		last_name=COALESCE(NULLIF(:last_name, ''), last_name),
		slug=COALESCE(NULLIF(:slug, ''), slug),
		email=COALESCE(NULLIF(:email, ''), email),
		birth_date=COALESCE(NULLIF(:birth_date, ''), birth_date),
		picture=COALESCE(NULLIF(:picture, ''), picture),
		nationality=COALESCE(NULLIF(:nationality, ''), nationality),
		street_address=COALESCE(NULLIF(:street_address, ''), street_address),
		city=COALESCE(NULLIF(:city, ''), city),
		province=COALESCE(NULLIF(:province, ''), province),
		zip_code=COALESCE(NULLIF(:zip_code, 0), zip_code),
		updated_at=now()
	WHERE artis_id = :artis_id;
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

func (repo *ArtiRepo) GetDataById(uid string) (*models.Artis, error) {
	q := `SELECT artis_id, slug, first_name, last_name, email, birth_date, picture, nationality, street_address, city, province, zip_code, created_at, updated_at FROM stream.artis WHERE artis_id = $1`

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
	q := `SELECT artis_id, slug, first_name, last_name, email, birth_date, picture, nationality, street_address, city, province, zip_code, created_at, updated_at FROM stream.artis WHERE slug = $1`

	var data models.Artis
	if err := repo.db.Get(&data, q, slug); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return &data, nil
}

func (repo *ArtiRepo) GetAllData(params *config.Pagination) (*config.ResultWarp, error) {
	var data = new(models.Arties)
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

	countQuery := fmt.Sprintf(`SELECT COUNT(artis_id) as "count" FROM stream.artis WHERE true %s`, filterQuery)
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

	q := fmt.Sprintf(`SELECT artis_id, slug, first_name, last_name, email, birth_date, picture, nationality, street_address, city, province, zip_code, created_at, updated_at FROM stream.artis WHERE true %s ORDER BY created_at DESC %s `, filterQuery, metaQuery)

	if err := repo.db.Select(data, repo.db.Rebind(q)); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(config.NotFound)
		}
		return nil, err
	}

	return &config.ResultWarp{Data: data, Meta: metaResult}, nil
}
