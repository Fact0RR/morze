INSERT INTO private_message (contact_id, user_id, data, additionals, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING id