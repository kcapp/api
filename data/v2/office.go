package data_v2

import "github.com/kcapp/api/models"

// GetOffices will return all offices
func GetOffices() ([]*models.Office, error) {
	rows, err := models.DB.Query("SELECT id, name, is_global, is_active FROM office")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offices := make([]*models.Office, 0)
	for rows.Next() {
		office := new(models.Office)
		err := rows.Scan(&office.ID, &office.Name, &office.IsGlobal, &office.IsActive)
		if err != nil {
			return nil, err
		}
		offices = append(offices, office)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return offices, nil
}
