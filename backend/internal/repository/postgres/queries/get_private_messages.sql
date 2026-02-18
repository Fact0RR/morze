SELECT 
    id,
    contact_id,
    user_id,
    data,
    additionals,
    created_at
FROM private_message
WHERE contact_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2 OFFSET $3;