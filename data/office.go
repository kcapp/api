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

	res, err := tx.Exec("INSERT INTO office (name, is_active, is_global) VALUES (?, ?, ?)", office.Name, office.IsActive, office.IsGlobal)
	if err != nil {
		tx.Rollback()
		return err
	}
	officeID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	// Add all existing players to this office if there are any not connected to a office
	_, err = tx.Exec("UPDATE player SET office_id = ? WHERE office_id IS NULL", officeID)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("Created new office (%d) %s", officeID, office.Name)

	// Update any players without office
	_, err = tx.Exec("UPDATE player SET office_id = ? WHERE office_id IS NULL", officeID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// UpdateOffice will update the given player
func UpdateOffice(officeID int, office models.Office) error {
	stmt, err := models.DB.Prepare(`UPDATE office SET name = ?, is_active = ?, is_global = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(office.Name, office.IsActive, office.IsGlobal, officeID)
	if err != nil {
		return err
	}
	log.Printf("Updated office (%d) (%s)", officeID, office.Name)
	return nil
}

// GetOffices will return all offices
func GetOffices() (map[int]*models.Office, error) {
	rows, err := models.DB.Query("SELECT id, name, is_global, is_active FROM office")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offices := make(map[int]*models.Office)
	for rows.Next() {
		office := new(models.Office)
		err := rows.Scan(&office.ID, &office.Name, &office.IsGlobal, &office.IsActive)
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
	err := models.DB.QueryRow("SELECT id, name, is_global, is_active FROM office WHERE id = ?", id).Scan(&office.ID, &office.Name, &office.IsGlobal, &office.IsActive)
	if err != nil {
		return nil, err
	}
	return office, nil
}
