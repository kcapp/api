package data

import (
	"github.com/kcapp/api/models"
)

// GetTournamentPresets returns all tournament presets
func GetTournamentPresets() ([]*models.TournamentPreset, error) {
	rows, err := models.DB.Query(`
		SELECT
			tp.id, tp.name, tp.starting_score, tp.description,
			tp.match_type_id, mt.name,
			mml16.id, mml16.name, mml16.short_name,
			mmqf.id, mmqf.name, mmqf.short_name,
			mmsf.id, mmsf.name, mmsf.short_name,
			mmgf.id, mmgf.name, mmgf.short_name,
			tg.id, tg.name,
			tp.player_id_walkover, tp.player_id_placeholder_home, tp.player_id_placeholder_away
		FROM tournament_preset tp
			JOIN match_type mt ON mt.id = tp.match_type_id
			JOIN match_mode mml16 ON mml16.id = tp.match_mode_id_last_16
			JOIN match_mode mmqf ON mmqf.id = tp.match_mode_id_quarter_final
			JOIN match_mode mmsf ON mmsf.id = tp.match_mode_id_semi_final
			JOIN match_mode mmgf ON mmgf.id = tp.match_mode_id_grand_final
			JOIN tournament_group tg ON tg.id = tp.playoffs_tournament_group_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	presets := make([]*models.TournamentPreset, 0)
	for rows.Next() {
		tp := new(models.TournamentPreset)
		tp.MatchModeLast16 = new(models.MatchMode)
		tp.MatchModeQuarterFinal = new(models.MatchMode)
		tp.MatchModeSemiFinal = new(models.MatchMode)
		tp.MatchModeGrandFinal = new(models.MatchMode)
		tp.MatchType = new(models.MatchType)
		tp.PlayoffsTournamentGroup = new(models.TournamentGroup)

		err := rows.Scan(&tp.ID, &tp.Name, &tp.StartingScore, &tp.Description, &tp.MatchType.ID, &tp.MatchType.Name,
			&tp.MatchModeLast16.ID, &tp.MatchModeLast16.Name, &tp.MatchModeLast16.ShortName,
			&tp.MatchModeQuarterFinal.ID, &tp.MatchModeQuarterFinal.Name, &tp.MatchModeQuarterFinal.ShortName,
			&tp.MatchModeSemiFinal.ID, &tp.MatchModeSemiFinal.Name, &tp.MatchModeSemiFinal.ShortName,
			&tp.MatchModeGrandFinal.ID, &tp.MatchModeGrandFinal.Name, &tp.MatchModeGrandFinal.ShortName,
			&tp.PlayoffsTournamentGroup.ID, &tp.PlayoffsTournamentGroup.Name,
			&tp.PlayerIDWalkover, &tp.PlayerIDPlaceholderHome, &tp.PlayerIDPlaceholderAway)
		if err != nil {
			return nil, err
		}
		presets = append(presets, tp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return presets, nil
}

// GetPreset returns the preset for the given ID
func GetTournamentPreset(id int) (*models.TournamentPreset, error) {
	tp := new(models.TournamentPreset)
	tp.MatchModeLast16 = new(models.MatchMode)
	tp.MatchModeQuarterFinal = new(models.MatchMode)
	tp.MatchModeSemiFinal = new(models.MatchMode)
	tp.MatchModeGrandFinal = new(models.MatchMode)
	tp.MatchType = new(models.MatchType)
	tp.PlayoffsTournamentGroup = new(models.TournamentGroup)
	err := models.DB.QueryRow(`
		SELECT
			tp.id, tp.name, tp.starting_score, tp.description,
			tp.match_type_id, mt.name,
			mml16.id, mml16.name, mml16.short_name,
			mmqf.id, mmqf.name, mmqf.short_name,
			mmsf.id, mmsf.name, mmsf.short_name,
			mmgf.id, mmgf.name, mmgf.short_name,
			tg.id, tg.name,
			tp.player_id_walkover, tp.player_id_placeholder_home, tp.player_id_placeholder_away
		FROM tournament_preset tp
			JOIN match_type mt ON mt.id = tp.match_type_id
			JOIN match_mode mml16 ON mml16.id = tp.match_mode_id_last_16
			JOIN match_mode mmqf ON mmqf.id = tp.match_mode_id_quarter_final
			JOIN match_mode mmsf ON mmsf.id = tp.match_mode_id_semi_final
			JOIN match_mode mmgf ON mmgf.id = tp.match_mode_id_grand_final
			JOIN tournament_group tg ON tg.id = tp.playoffs_tournament_group_id
		WHERE tp.id = ?`, id).
		Scan(&tp.ID, &tp.Name, &tp.StartingScore, &tp.Description, &tp.MatchType.ID, &tp.MatchType.Name,
			&tp.MatchModeLast16.ID, &tp.MatchModeLast16.Name, &tp.MatchModeLast16.ShortName,
			&tp.MatchModeQuarterFinal.ID, &tp.MatchModeQuarterFinal.Name, &tp.MatchModeQuarterFinal.ShortName,
			&tp.MatchModeSemiFinal.ID, &tp.MatchModeSemiFinal.Name, &tp.MatchModeSemiFinal.ShortName,
			&tp.MatchModeGrandFinal.ID, &tp.MatchModeGrandFinal.Name, &tp.MatchModeGrandFinal.ShortName,
			&tp.PlayoffsTournamentGroup.ID, &tp.PlayoffsTournamentGroup.Name,
			&tp.PlayerIDWalkover, &tp.PlayerIDPlaceholderHome, &tp.PlayerIDPlaceholderAway)
	if err != nil {
		return nil, err
	}
	return tp, nil
}

// AddTournamentPreset will add a new preset to the database
/*func AddTournamentPreset(preset models.TournamentPreset) error {
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

// UpdateTournamentPreset will update the given preset
func UpdateTournamentPreset(id int, preset models.TournamentPreset) error {
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
*/
