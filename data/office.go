package data

import "github.com/kcapp/api/models"

// GetOffices will return all offices
func GetOffices() (map[int]*models.Office, error) {
	rows, err := models.DB.Query("SELECT id, name FROM office")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offices := make(map[int]*models.Office, 0)
	for rows.Next() {
		office := new(models.Office)
		err := rows.Scan(&office.ID, &office.Name)
		if err != nil {
			return nil, err
		}
		offices[office.ID] = office
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return offices, nil
}

// GetOffice will return a office for the given id
func GetOffice(id int) (*models.Office, error) {
	office := new(models.Office)
	err := models.DB.QueryRow("SELECT id, name FROM office WHERE id = ?", id).Scan(&office.ID, &office.Name)
	if err != nil {
		return nil, err
	}
	return office, nil
}
