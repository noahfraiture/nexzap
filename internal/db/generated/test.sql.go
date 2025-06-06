// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: test.sql

package db

import (
	"context"
)

const findCorrectionSheet = `-- name: FindCorrectionSheet :many
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
  tu.title, s.id, s.docker_image, s.command, s.submission_name, s.correction_content
`

type FindCorrectionSheetRow struct {
	Title             string
	Page              int32
	SubmissionName    string
	CorrectionContent string
	DockerImage       string
	Command           string
	FilesName         []string
	FilesContent      []string
}

func (q *Queries) FindCorrectionSheet(ctx context.Context) ([]FindCorrectionSheetRow, error) {
	rows, err := q.db.Query(ctx, findCorrectionSheet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindCorrectionSheetRow
	for rows.Next() {
		var i FindCorrectionSheetRow
		if err := rows.Scan(
			&i.Title,
			&i.Page,
			&i.SubmissionName,
			&i.CorrectionContent,
			&i.DockerImage,
			&i.Command,
			&i.FilesName,
			&i.FilesContent,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findSpecificCorrectionSheet = `-- name: FindSpecificCorrectionSheet :one
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
  tu.title = $1 AND s.page = $2
GROUP BY
  tu.title, s.id, s.docker_image, s.command, s.submission_name, s.correction_content
`

type FindSpecificCorrectionSheetParams struct {
	Title string
	Page  int32
}

type FindSpecificCorrectionSheetRow struct {
	SubmissionName    string
	CorrectionContent string
	DockerImage       string
	Command           string
	FilesName         []string
	FilesContent      []string
}

func (q *Queries) FindSpecificCorrectionSheet(ctx context.Context, arg FindSpecificCorrectionSheetParams) (FindSpecificCorrectionSheetRow, error) {
	row := q.db.QueryRow(ctx, findSpecificCorrectionSheet, arg.Title, arg.Page)
	var i FindSpecificCorrectionSheetRow
	err := row.Scan(
		&i.SubmissionName,
		&i.CorrectionContent,
		&i.DockerImage,
		&i.Command,
		&i.FilesName,
		&i.FilesContent,
	)
	return i, err
}
