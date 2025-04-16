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

-- name: InsertCompleteTutorial :one
WITH lang_ins AS (
  INSERT INTO languages (name)
  VALUES (@language_name)
  ON CONFLICT (name) DO NOTHING
  RETURNING id
), lang_sel AS (
  SELECT id
  FROM languages
  WHERE name = @language_name
), tut AS (
  INSERT INTO tutorials (language_id)
  VALUES ((SELECT id FROM lang_ins UNION SELECT id FROM lang_sel))
  RETURNING id
), sheet AS (
  INSERT INTO sheets (guide_content, test_content, docker_image, tutorial_id)
  SELECT unnest(@guide_contents::text[]), unnest(@test_contents::text[]), unnest(@docker_images::text[]), (SELECT id FROM tut)
  RETURNING id
)
SELECT (SELECT id FROM tut) AS tutorial_id;
