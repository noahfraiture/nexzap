-- name: InsertTutorial :many
WITH tutorial AS (
  INSERT INTO tutorials (title, highlight, code_editor, version, unlock)
  VALUES (@title, @highlight, @code_editor, @version, @unlock)
  RETURNING id
), sheet AS (
  INSERT INTO sheets (
    tutorial_id,
    page,
    guide_content,
    exercise_content,
    submission_name,
    submission_content,
    docker_image,
    command
  )
  SELECT
    (SELECT id FROM tutorial),
    unnest(@pages::integer[]),
    unnest(@guides_content::text[]),
    unnest(@exercises_content::text[]),
    unnest(@submissions_name::text[]),
    unnest(@submissions_content::text[]),
    unnest(@docker_images::text[]),
    unnest(@commands::text[])
  RETURNING id
)
SELECT id FROM sheet;

-- name: InsertFiles :exec
INSERT INTO files (name, content, sheet_id)
SELECT unnest(@names::text[]), unnest(@contents::text[]), @sheet_id;

-- name: FindLastTutorialFirstSheet :one
SELECT
  tu.title,
  tu.highlight,
  tu.code_editor,
  s.id,
  s.guide_content,
  s.exercise_content,
  s.page,
  s.submission_content,
  (SELECT COUNT(page) FROM sheets sh WHERE sh.tutorial_id = tu.id) as total_pages
FROM
  tutorials tu
  JOIN sheets s ON s.tutorial_id = tu.id
WHERE
  s.page = 1
  AND tu.unlock < NOW ()
ORDER BY
  tu.updated_at DESC,
  tu.version DESC
LIMIT
  1;

-- name: FindLastTutorialSheet :one
SELECT
  tu.title,
  tu.highlight,
  tu.code_editor,
  s.id,
  s.guide_content,
  s.exercise_content,
  s.page,
  s.submission_content,
  (SELECT COUNT(page) FROM sheets sh WHERE sh.tutorial_id = tu.id) as total_pages
FROM
  tutorials tu
  JOIN sheets s ON s.tutorial_id = tu.id
WHERE
  s.page = @page
  AND tu.unlock < NOW ()
ORDER BY
  tu.updated_at DESC,
  tu.version DESC
LIMIT
  1;

-- name: FindSubmissionData :one
SELECT
  s.docker_image,
  s.command,
  s.submission_name,
  array_agg(f.name)::text[] AS files_name,
  array_agg(f.content)::text[] AS files_content
FROM
  sheets s
  JOIN files f ON f.sheet_id = s.id
WHERE
  s.id = @sheet_id
GROUP BY
  s.id, s.docker_image, s.command, s.submission_name;
