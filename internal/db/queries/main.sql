-- name: FindLastTutorial :one
SELECT
  array_agg (s.content)::text[] AS sheet_contents,
  array_agg (te.content)::text[] AS test_contents,
  l.name AS language_name
FROM
  tutorials tu
  JOIN sheets s ON s.tutorial_id = tu.id
  JOIN tests te ON te.id = s.test_id
  JOIN languages l ON l.id = tu.language_id
GROUP BY
  tu.id,
  te.docker_image,
  l.name
ORDER BY
  tu.updated_at DESC
LIMIT
  1;

-- name: InsertCompleteTutorial :one
WITH lang AS (
  INSERT INTO languages (name)
  VALUES (@language_name)
  RETURNING id
),
test AS (
  INSERT INTO tests (content, docker_image)
  SELECT unnest(@test_contents::text[]), @docker_image
  RETURNING id
),
tut AS (
  INSERT INTO tutorials (language_id)
  VALUES ((SELECT id FROM lang))
  RETURNING id
),
sheet AS (
  INSERT INTO sheets (content, test_id, tutorial_id)
  SELECT unnest(@sheet_contents::text[]), (SELECT id FROM test LIMIT 1), (SELECT id FROM tut)
  RETURNING id
)
SELECT (SELECT id FROM tut) AS tutorial_id;
