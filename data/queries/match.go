package queries

func QueryMatchTypes() string {
	return "SELECT id, `name`, description FROM match_type"
}

func QueryMatchModes() string {
	return "SELECT id, wins_required, legs_required, tiebreak_match_type_id, `name`, short_name FROM match_mode ORDER BY wins_required"
}
