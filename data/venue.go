package data

import (
	"github.com/kcapp/api/models"
)

// GetVenues will return all venues
func GetVenues() ([]*models.Venue, error) {
	rows, err := models.DB.Query("SELECT id, name, description FROM venue")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	venues := make([]*models.Venue, 0)
	for rows.Next() {
		venue := new(models.Venue)
		err := rows.Scan(&venue.ID, &venue.Name, &venue.Description)
		if err != nil {
			return nil, err
		}
		venues = append(venues, venue)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return venues, nil
}
