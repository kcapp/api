package data

import "github.com/kcapp/api/models"

func GetBadges() ([]*models.Badge, error) {
	rows, err := models.DB.Query(`
		SELECT
			b.id,
			b.name,
			b.description,
			b.filename
		FROM badge b`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make([]*models.Badge, 0)
	for rows.Next() {
		badge := new(models.Badge)
		err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.Description,
			&badge.Filename,
		)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}
	return badges, nil
}
