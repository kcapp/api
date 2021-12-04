package data

import (
	"log"

	"github.com/kcapp/api/models"
)

// GetPresets returns all match presets
func GetPresets() ([]*models.MatchPreset, error) {
	rows, err := models.DB.Query(`
		SELECT
			mp.id, mp.name, match_type_id, mt.name, match_mode_id, mm.name, mm.short_name,
			starting_score, smartcard_uid, mp.description
		FROM match_preset mp
			JOIN match_type mt ON mt.id = mp.match_type_id
			JOIN match_mode mm ON mm.id = mp.match_mode_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	presets := make([]*models.MatchPreset, 0)
	for rows.Next() {
		p := new(models.MatchPreset)
		mt := new(models.MatchType)
		mm := new(models.MatchMode)
		err := rows.Scan(&p.ID, &p.Name, &mt.ID, &mt.Name, &mm.ID, &mm.Name, &mm.ShortName, &p.StartingScore, &p.SmartcardUID, &p.Description)
		if err != nil {
			return nil, err
		}
		p.MatchMode = mm
		p.MatchType = mt
		presets = append(presets, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return presets, nil
}

// GetPreset returns the preset for the given ID
func GetPreset(id int) (*models.MatchPreset, error) {
	p := new(models.MatchPreset)
	mt := new(models.MatchType)
	mm := new(models.MatchMode)
	err := models.DB.QueryRow(`
		SELECT
			mp.id, mp.name, match_type_id, mt.name, match_mode_id, mm.name, mm.short_name,
			starting_score, smartcard_uid, mp.description
		FROM match_preset mp
			JOIN match_type mt ON mt.id = mp.match_type_id
			JOIN match_mode mm ON mm.id = mp.match_mode_id
		WHERE mp.id = ?`, id).
		Scan(&p.ID, &p.Name, &mt.ID, &mt.Name, &mm.ID, &mm.Name, &mm.ShortName, &p.StartingScore, &p.SmartcardUID, &p.Description)
	if err != nil {
		return nil, err
	}
	p.MatchMode = mm
	p.MatchType = mt

	return p, nil
}

// AddPreset will add a new preset to the database
func AddPreset(preset models.MatchPreset) error {
	stmt, err := models.DB.Prepare(`
		INSERT INTO match_preset(name, match_type_id, match_mode_id, starting_score, smartcard_uid, description) VALUES(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(preset.Name, preset.MatchType.ID, preset.MatchMode.ID, preset.StartingScore, preset.SmartcardUID, preset.Description)
	if err != nil {
		return err
	}
	log.Printf("Created preset %s (%v)", preset.Name, preset)
	return nil
}

// UpdatePreset will update the given preset
func UpdatePreset(id int, preset models.MatchPreset) error {
	stmt, err := models.DB.Prepare(`
		UPDATE match_preset SET
			name = ?, match_type_id = ?, match_mode_id = ?, starting_score = ?, smartcard_uid = ?, description = ?
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(preset.Name, preset.MatchType.ID, preset.MatchMode.ID, preset.StartingScore, preset.SmartcardUID, preset.Description, id)
	if err != nil {
		return err
	}
	log.Printf("Updated preset %s (%v)", preset.Name, preset)
	return nil
}

// DeletePreset will delete the given preset
func DeletePreset(id int) error {
	stmt, err := models.DB.Prepare(`DELETE FROM match_preset WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	log.Printf("Deleted preset %d", id)
	return nil
}
