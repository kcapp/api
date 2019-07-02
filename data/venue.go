package data

import (
	"log"

	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
)

// GetVenues will return all venues
func GetVenues() ([]*models.Venue, error) {
	rows, err := models.DB.Query("SELECT id, name, office_id, description FROM venue")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	venues := make([]*models.Venue, 0)
	for rows.Next() {
		venue := new(models.Venue)
		err := rows.Scan(&venue.ID, &venue.Name, &venue.OfficeID, &venue.Description)
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

// GetVenue will return a venue for the given id
func GetVenue(id int) (*models.Venue, error) {
	venue := new(models.Venue)
	err := models.DB.QueryRow("SELECT id, name, office_id, description FROM venue WHERE id = ?", id).Scan(&venue.ID, &venue.Name, &venue.OfficeID, &venue.Description)
	if err != nil {
		return nil, err
	}
	venue.Config, err = GetVenueConfiguration(id)
	if err != nil {
		log.Printf("Unable to get venue configuration for %d", id)
	}
	return venue, nil
}

// GetVenueConfiguration will return the configuration for a venue with the given id
func GetVenueConfiguration(id int) (*models.VenueConfig, error) {
	config := new(models.VenueConfig)
	err := models.DB.QueryRow("SELECT venue_id, has_dual_monitor, has_leg_lights, has_smartboard, smartboard_uuid, smartboard_button_number FROM venue_configuration WHERE venue_id = ?",
		id).Scan(&config.VenueID, &config.HasDualMonitor, &config.HasLEDLights, &config.HasSmartboard, &config.SmartboardUUID, &config.SmartboardButtonNumber)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// SpectateVenue will return the current active match at a venue
func SpectateVenue(venueID int) ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.current_leg_id, m.winner_id, m.created_at, m.updated_at, m.owe_type_id, m.venue_id,
			mt.id, mt.name, mt.description, mm.id, mm.name, mm.short_name, mm.wins_required, mm.legs_required,
			ot.id, ot.item, v.id, v.name, v.description,
			l.updated_at as 'last_throw', GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players'
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			LEFT JOIN leg l ON l.id = m.current_leg_id
			LEFT JOIN owe_type ot ON ot.id = m.owe_type_id
			LEFT JOIN venue v on v.id = m.venue_id
			LEFT JOIN player2leg p2l ON p2l.match_id = m.id
		WHERE m.venue_id = ? AND m.is_finished = 0
		GROUP BY m.id
		ORDER BY m.id
		LIMIT 1`, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		m := new(models.Match)
		m.MatchType = new(models.MatchType)
		m.MatchMode = new(models.MatchMode)
		ot := new(models.OweType)
		venue := new(models.Venue)
		var players string
		err := rows.Scan(&m.ID, &m.IsFinished, &m.CurrentLegID, &m.WinnerID, &m.CreatedAt, &m.UpdatedAt, &m.OweTypeID, &m.VenueID,
			&m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
			&ot.ID, &ot.Item, &venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players)
		if err != nil {
			return nil, err
		}
		if m.OweTypeID.Valid {
			m.OweType = ot
		}
		if m.VenueID.Valid {
			m.Venue = venue
		}
		m.Players = util.StringToIntArray(players)
		matches = append(matches, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}
