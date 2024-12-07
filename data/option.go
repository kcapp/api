package data

import (
	"database/sql"

	"github.com/kcapp/api/models"
)

// GetDefaultOptions will return default options
func GetDefaultOptions() (*models.DefaultOptions, error) {
	opts := new(models.DefaultOptions)
	opts.MatchType = new(models.MatchType)
	opts.MatchMode = new(models.MatchMode)
	opts.OutshotType = new(models.OutshotType)

	err := models.DB.QueryRow(`
		SELECT
			mt.id, mt.name, mt.description,
			mm.id, mm.name, mm.short_name, mm.wins_required, mm.legs_required,
			starting_score, max_rounds, ot.id, ot.name, ot.short_name
		FROM match_default md
		    LEFT JOIN match_type mt ON mt.id = md.match_type_id
    		LEFT JOIN match_mode mm ON mm.id = md.match_mode_id
    		LEFT JOIN outshot_type ot ON ot.id = md.outshot_type_id
		LIMIT 1`).
		Scan(&opts.MatchType.ID, &opts.MatchType.Name, &opts.MatchType.Description, &opts.MatchMode.ID, &opts.MatchMode.Name,
			&opts.MatchMode.ShortName, &opts.MatchMode.WinsRequired, &opts.MatchMode.LegsRequired, &opts.StartingScore, &opts.MaxRounds,
			&opts.OutshotType.ID, &opts.OutshotType.Name, &opts.OutshotType.ShortName)
	if err != nil {
		if err == sql.ErrNoRows {
			return opts, nil
		}
		return nil, err
	}
	return opts, nil
}
