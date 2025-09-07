-- name: GetMovie :one
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1 LIMIT 1;

-- name: InsertMovie :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version;

-- name: UpdateMovie :one
UPDATE movies 
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;

-- uh oh sqlc cant handle this
---- name GetAllMovies :many
--fmt.Sprintf(`SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
--FROM movies
--WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
--AND (genres @> $2 OR $2 = '{}')
--ORDER BY %s %s, id ASC
--LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
