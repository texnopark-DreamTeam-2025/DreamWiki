package repository

import (
	"strconv"
	"strings"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func decodeIntegrationLogsCursor(cursor *string) (timeFrom time.Time, idFrom int64) {
	minTime := time.Unix(0, 0)

	if cursor == nil {
		return minTime, 0
	}

	splittedCursor := strings.Split(*cursor, "\n")
	if len(splittedCursor) != 2 {
		return minTime, 0
	}

	timeFrom, err := time.Parse(time.RFC3339, splittedCursor[0])
	if err != nil {
		return minTime, 0
	}

	idFrom, err = strconv.ParseInt(splittedCursor[1], 10, 64)
	if err != nil {
		return minTime, 0
	}

	return timeFrom, idFrom
}

func encodeIntegrationLogsCursor(timeFrom time.Time, idFrom int64) string {
	return timeFrom.Format(time.RFC3339) + "\n" + strconv.FormatInt(idFrom, 10)
}

func encodeIntegrationLogsNextInfo(timeFrom time.Time, idFrom int64, numRows int) *api.NextInfo {
	return &api.NextInfo{
		Cursor:  encodeIntegrationLogsCursor(timeFrom, idFrom),
		HasMore: numRows > 0,
	}
}

func (r *appRepositoryImpl) WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error {
	yql := `INSERT INTO IntegrationLogField (integration_id, log_text, created_at)
	VALUES ($integrationID, $logText, CurrentUtcDatetime())`

	result, err := r.ydbClient.OutsideTX().Execute(yql,
		table.ValueParam("$integrationID", types.TextValue(string(integrationID))),
		table.ValueParam("$logText", types.TextValue(logText)),
	)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) GetIntegrationLogFields(integrationID api.IntegrationID, cursor *api.Cursor, limit uint64) ([]api.IntegrationLogField, *api.NextInfo, error) {
	yql := `
		SELECT field_id, log_text, created_at
		FROM IntegrationLogField
		WHERE integration_id=$integrationID
			AND (
				created_at > $timeFrom
				OR (created_at = $timeFrom AND field_id > $idFrom)
			)
		ORDER BY created_at DESC
		LIMIT $limit
	`

	timeFrom, idFrom := decodeIntegrationLogsCursor(cursor)

	result, err := r.ydbClient.InTX().Execute(yql,
		table.ValueParam("$integrationID", types.TextValue(string(integrationID))),
		table.ValueParam("$limit", types.Uint64Value(limit)),
		table.ValueParam("$timeFrom", types.TimestampValueFromTime(timeFrom)),
		table.ValueParam("$idFrom", types.Uint64Value(uint64(idFrom))),
	)
	if err != nil {
		return nil, nil, err
	}
	defer result.Close()

	fields := make([]api.IntegrationLogField, 0, limit)

	newIDFrom := int64(0)
	newTimeFrom := time.Time{}
	for result.NextRow() {
		var content string
		var createdAt time.Time
		err := result.FetchRow(&newIDFrom, &content, &createdAt)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, api.IntegrationLogField{
			Content:   content,
			CreatedAt: createdAt,
		})
		if newTimeFrom.Before(createdAt) {
			newTimeFrom = createdAt
		}
	}

	return fields, encodeIntegrationLogsNextInfo(newTimeFrom, newIDFrom, len(fields)), nil
}
