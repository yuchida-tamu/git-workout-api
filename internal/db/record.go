package db

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/yuchida-tamu/git-workout-api/internal/record"
)

type RecordRow struct {
	ID          string
	DateCreated string
	MessageBody string
	Author      string
}

func convertRecordRowToRecord(row RecordRow) record.Record {
	return record.Record{
		ID:          row.ID,
		DateCreated: row.DateCreated,
		MessageBody: row.MessageBody,
		Author:      row.Author,
	}
}

func (d *Database) GetRecordsByAuthor(ctx context.Context, ID string) ([]record.Record, error) {
	var records []record.Record
	rows, err := d.Client.QueryContext(
		ctx,
		`SELECT *
		FROM records
		WHERE author = $1`,
		ID,
	)
	if err != nil {
		return []record.Record{}, fmt.Errorf("error fetching records by author id: %w", err)
	}

	for rows.Next() {
		var ID, dateCreated, messageBody, author string
		err := rows.Scan(&ID, &dateCreated, &messageBody, &author)
		if err != nil {
			return []record.Record{}, fmt.Errorf("error fetching records by author id: %w", err)
		}

		records = append(records,
			record.Record{
				ID:          ID,
				DateCreated: dateCreated,
				MessageBody: messageBody,
				Author:      author,
			})
	}

	return records, nil
}

func (d *Database) GetRecordById(ctx context.Context, ID string) (record.Record, error) {
	var recordRow RecordRow

	row := d.Client.QueryRowContext(
		ctx,
		`SELECT * FROM records WHERE id = $1`,
		ID,
	)

	err := row.Scan(&recordRow.ID, &recordRow.DateCreated, &recordRow.MessageBody, &recordRow.Author)
	if err != nil {
		return record.Record{}, fmt.Errorf("error fetching the record by id: %w", err)
	}

	return convertRecordRowToRecord(recordRow), nil
}

func (d *Database) PostRecord(ctx context.Context, rcd record.Record) (record.Record, error) {
	rcd.ID = uuid.NewV4().String()
	postRow := RecordRow{
		ID:          rcd.ID,
		DateCreated: rcd.DateCreated,
		MessageBody: rcd.MessageBody,
		Author:      rcd.Author,
	}

	row, err := d.Client.NamedQueryContext(
		ctx,
		`INSERT INTO records
		(id, date_created, message_body, author)
		VALUES
		(:id, :datecreated, :messagebody, :author)`,
		postRow,
	)

	if err != nil {
		return record.Record{}, fmt.Errorf("failed to insert record: %w", err)
	}

	if err := row.Close(); err != nil {
		return record.Record{}, fmt.Errorf("failed to insert record: %w", err)
	}

	return rcd, nil
}

func (d *Database) UpdateRecord(ctx context.Context, ID string, rcd record.Record) (record.Record, error) {
	recordRow := RecordRow{
		ID:          ID,
		DateCreated: rcd.DateCreated,
		MessageBody: rcd.MessageBody,
		Author:      rcd.Author,
	}

	row, err := d.Client.NamedQueryContext(
		ctx,
		`UPDATE records SET
		date_created = :datecreated,
		message_body = :messagebody,
		WHERE id = :id`,
		recordRow,
	)
	if err != nil {
		return record.Record{}, fmt.Errorf("failed to update record: %w", err)
	}
	if err := row.Close(); err != nil {
		return record.Record{}, fmt.Errorf("failed to update record: %w", err)
	}

	return convertRecordRowToRecord(recordRow), nil
}

func (d *Database) DeleteRecord(ctx context.Context, ID string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM records WHERE id = $1`,
		ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete record from database: %w", err)
	}

	return nil
}
