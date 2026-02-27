UPDATE private_message 
SET
    data = $2, 
    updated = TRUE 
WHERE id = $1;