-- name: InsertTutorial :exec
WITH tutorial AS (
  INSERT INTO tutorials (title, highlight, code_editor, docker_image, version)
  VALUES (@title, @highlight, @code_editor, @docker_image, @version)
  RETURNING id
), sheet AS (
  INSERT INTO sheets (guide_content, exercise_content, tutorial_id)
  SELECT unnest(@guides_content::text[]), unnest(@exercises_content::text[]), (SELECT id FROM tutorial)
  RETURNING id
), file AS (
  INSERT INTO files (name, content, sheet_id)
  SELECT unnest(@files_name::text[]), unnest(@files_content::text[]), (SELECT id FROM sheet)
)
SELECT id FROM tutorial;


-- name: FindLastTutorial :one
SELECT
  array_agg (s.guide_content)::text[] AS guide_contents,
  array_agg (s.test_content)::text[] AS test_contents,
  l.name AS language_name
FROM
  tutorials tu
  JOIN sheets s ON s.tutorial_id = tu.id
  JOIN languages l ON l.id = tu.language_id
GROUP BY
  tu.id,
  l.name
ORDER BY
  tu.updated_at DESC
LIMIT
  1;
