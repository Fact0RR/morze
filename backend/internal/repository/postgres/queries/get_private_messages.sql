SELECT 
    id,
    contact_id,
    user_id,
    data,
    additionals,
    created_at,
    updated
FROM private_message
WHERE contact_id = $1 and deleted = FALSE
ORDER BY created_at DESC, id DESC
LIMIT $2 OFFSET $3;