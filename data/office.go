package data

import (
	"log"

	"github.com/kcapp/api/models"
)

// AddOffice will add a new office
func AddOffice(office models.Office) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}

	// Prepare statement for inserting data
	res, err := tx.Exec("INSERT INTO office (name, is_active) VALUES (?, ?)", office.Name, office.IsActive)
	if err != nil {
		tx.Rollback()
		return err
	}
	officeID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("Created new office (%d) %s", officeID, office.Name)
	tx.Commit()
	return nil
}

// UpdateOffice will update the given player
func UpdateOffice(officeID int, office models.Office) error {
	stmt, err := models.DB.Prepare(`UPDATE office SET name = ?, is_active = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(office.Name, office.IsActive, officeID)
	if err != nil {
		return err
	}
	log.Printf("Updated office (%d) (%s)", officeID, office.Name)
	return nil
}

// GetOffices will return all offices
func GetOffices() (map[int]*models.Office, error) {
	rows, err := models.DB.Query("SELECT id, name, is_active FROM office")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offices := make(map[int]*models.Office, 0)
	for rows.Next() {
		office := new(models.Office)
		err := rows.Scan(&office.ID, &office.Name, &office.IsActive)
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
	err := models.DB.QueryRow("SELECT id, name, is_active FROM office WHERE id = ?", id).Scan(&office.ID, &office.Name, &office.IsActive)
	if err != nil {
		return nil, err
	}
	return office, nil
}
