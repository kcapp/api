package queries

func QueryAllPlayers() string {
	return `SELECT
			p.id, p.first_name, p.last_name, p.vocal_name, p.nickname, p.slack_handle, p.color, p.profile_pic_url, p.smartcard_uid,
			 p.board_stream_url, p.board_stream_css, p.active, p.office_id, p.is_bot, p.is_placeholder, p.created_at, p.updated_at
		FROM player p`
}

func QueryActivePlayers() string {
	return `SELECT
			p.id, p.first_name, p.last_name, p.vocal_name, p.nickname, p.slack_handle, p.color, p.profile_pic_url, p.smartcard_uid,
			 p.board_stream_url, p.board_stream_css, p.active, p.office_id, p.is_bot, p.is_placeholder, p.created_at, p.updated_at
		FROM player p WHERE active = 1`
}

func QueryPlayer() string {
	return `
		SELECT
			p.id, p.first_name, p.last_name, p.vocal_name, p.nickname,
			p.slack_handle, p.color, p.profile_pic_url, p.smartcard_uid, p.board_stream_url, p.board_stream_css,
			p.office_id, p.active, p.is_bot, p.is_placeholder, p.created_at, p.updated_at, pe.current_elo, pe.tournament_elo
		FROM player p
		JOIN player_elo pe on pe.player_id = p.id
		WHERE p.id = ?`
}

func QueryMatchesPlayed() string {
	return `
		SELECT
			player_id,
			MAX(matches_played) AS 'matches_played',
			MAX(matches_won) AS 'matches_won',
			MAX(legs_played) AS 'legs_played',
			MAX(legs_won) AS 'legs_won'
		FROM (
			SELECT
				p2l.player_id,
				COUNT(DISTINCT p2l.match_id) AS 'matches_played',
				0 AS 'matches_won',
				COUNT(m.id)  AS 'legs_played',
				SUM(CASE WHEN p2l.player_id = m.winner_id THEN 1 ELSE 0 END) AS 'legs_won'
			FROM player2leg p2l
				JOIN leg l ON l.id = p2l.leg_id
				JOIN matches m ON m.id = p2l.match_id
			WHERE l.is_finished = 1 AND m.is_abandoned = 0
			GROUP BY p2l.player_id
			UNION ALL
			SELECT
				p2l.player_id,
				0 AS 'matches_played',
				COUNT(DISTINCT m.id) AS 'matches_won',
				0 AS 'legs_played',
				0 AS 'legs_won'
			FROM matches m
				JOIN leg l ON l.match_id = m.id
				JOIN player2leg p2l ON p2l.player_id = m.winner_id AND p2l.match_id = m.id
			WHERE l.is_finished = 1 AND m.is_abandoned = 0
			GROUP BY m.winner_id
		) matches
		GROUP BY player_id`
}

func QueryPlayersX01StatisticsComparison() string {
	return `
	SELECT
	p.id AS 'player_id',
		COUNT(DISTINCT m.id) AS 'matches_played',
		COUNT(DISTINCT m2.id) AS 'matches_won',
		COUNT(DISTINCT l.id) AS 'legs_played',
		COUNT(DISTINCT l2.id) AS 'legs_won',
		SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
		SUM(s.first_nine_ppd) / COUNT(p.id) AS 'first_nine_ppd',
		(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
		(SUM(s.first_nine_ppd) / COUNT(p.id)) * 3 AS 'first_nine_three_dart_avg',
		SUM(s.60s_plus) AS '60s_plus',
		SUM(s.100s_plus) AS '100s_plus',
		SUM(s.140s_plus) AS '140s_plus',
		SUM(s.180s) AS '180s',
		SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
		SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
		SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
		COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage'
	FROM statistics_x01 s
	JOIN player p ON p.id = s.player_id
	JOIN leg l ON l.id = s.leg_id
	JOIN matches m ON m.id = l.match_id
	LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
	LEFT JOIN matches m2 ON m2.id = l2.match_id AND l2.winner_id = p.id
	WHERE s.player_id IN (?)
	AND l.starting_score IN (?)
	AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
	AND m.match_type_id = 1
	GROUP BY s.player_id
	ORDER BY p.id`
}

func QueryPlayerProgression() string {
	return `
		SELECT
			s.player_id AS 'player_id',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(s.player_id) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
			SUM(s.first_nine_ppd) / COUNT(s.player_id) * 3 AS 'first_nine_three_dart_avg',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
			SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
			SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage',
			DATE(m.updated_at) AS 'date'
		FROM statistics_x01 s
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND m.match_type_id = 1
			AND m.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
		GROUP BY YEAR(m.updateD_at), WEEK(m.updated_at)
		ORDER BY date DESC`
}
