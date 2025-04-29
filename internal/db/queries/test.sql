-- name: FindCorrectionSheet :many
SELECT
  tu.title, -- debug log
  s.page, -- debug log
  s.submission_name,
  s.correction_content,
  s.docker_image,
  s.command,
  array_agg(f.name)::text[] AS files_name,
  array_agg(f.content)::text[] AS files_content
FROM
  sheets s
  JOIN files f ON f.sheet_id = s.id
  JOIN tutorials tu ON tu.id = s.tutorial_id
GROUP BY
  tu.title, s.id, s.docker_image, s.command, s.submission_name, s.correction_content;

-- name: FindSpecificCorrectionSheet :one
SELECT
  s.submission_name,
  s.correction_content,
  s.docker_image,
  s.command,
  array_agg(f.name)::text[] AS files_name,
  array_agg(f.content)::text[] AS files_content
FROM
  sheets s
  JOIN files f ON f.sheet_id = s.id
  JOIN tutorials tu ON tu.id = s.tutorial_id
WHERE
  tu.title = @title AND s.page = @page
GROUP BY
  tu.title, s.id, s.docker_image, s.command, s.submission_name, s.correction_content;

