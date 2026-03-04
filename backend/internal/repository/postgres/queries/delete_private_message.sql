UPDATE private_message 
SET deleted = TRUE 
WHERE id = $1;
