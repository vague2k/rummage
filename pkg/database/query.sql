-- name: AddItem :one
INSERT INTO rummage_items (
    entry, score, lastaccessed
) VALUES (
    ?, ?, ?
)
RETURNING * ;


-- name: SelectItem :one
SELECT * FROM rummage_items
WHERE entry = ?
LIMIT 1 ;

-- name: UpdateItem :exec
UPDATE rummage_items
SET score = ?, lastaccessed = ?
WHERE entry = ?
RETURNING * ;

-- name: EntryWithHighestScore :one
SELECT * FROM rummage_items
WHERE entry LIKE ?
ORDER BY score
DESC LIMIT 1 ;

-- name: FindTopNMatches :many
SELECT * FROM rummage_items
WHERE entry LIKE ?
ORDER BY score
DESC LIMIT ? ;

-- name: DeleteItem :exec
DELETE FROM rummage_items
WHERE entry = ? ;

-- name: DeleteAllItem :exec
DELETE FROM rummage_items ;
