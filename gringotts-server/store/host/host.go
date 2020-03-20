package host

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jinlingan/gringotts/gringotts-server/store/db"

	"github.com/jinlingan/gringotts/gringotts-server/model"
	"github.com/jmoiron/sqlx"
)

// New returns a new UserStore.
func New(db *sqlx.DB) model.HostService {
	return &hostStore{db}
}

type hostStore struct {
	db *sqlx.DB
}

func (h *hostStore) Find(ctx context.Context, agentID string) (*model.Host, error) {
	if agentID == "" {
		return nil, nil
	}
	out := &model.Host{}
	params := toParams(out)
	query, args, err := h.db.BindNamed(queryKey, params)
	if err != nil {
		return nil, err
	}
	row := h.db.QueryRow(query, args...)
	err = scanRow(row, out)
	return out, err
}
func scanRow(scanner db.Scanner, dest *model.Host) error {
	err := scanner.Scan(
		&dest.AgentID,
		&dest.HostName,
		&dest.HostUUID,
		&dest.Os,
		&dest.Platform,
		&dest.PlatformFamily,
		&dest.PlatformVersion,
		&dest.KernelVersion,
		&dest.VirtualizationSystem,
		&dest.VirtualizationRole,
		&dest.InterfacesJSON,
		&dest.CreateTime,
		&dest.UpdateTime,
		&dest.LastHeartBeatTime,
	)

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(dest.InterfacesJSON), &dest.Interfaces)

}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*model.Host, error) {
	defer rows.Close()

	users := []*model.Host{}
	for rows.Next() {
		user := new(model.Host)
		err := scanRow(rows, user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func toParams(h *model.Host) map[string]interface{} {
	return map[string]interface{}{
		"agent_id":              h.AgentID,
		"host_name":             h.HostName,
		"host_UUID":             h.HostUUID,
		"os":                    h.Os,
		"platform":              h.Platform,
		"platform_family":       h.PlatformFamily,
		"platform_version":      h.PlatformVersion,
		"kernel_version":        h.KernelVersion,
		"virtualization_system": h.VirtualizationSystem,
		"virtualization_role":   h.VirtualizationRole,
		"interfaces_json":       h.InterfacesJSON,
		"create_time":           h.CreateTime,
		"update_time":           h.UpdateTime,
		"last_heartbeat_time":   h.LastHeartBeatTime,
	}
}
func (s *hostStore) Create(ctx context.Context, host *model.Host) error {

	params := toParams(host)
	nowTime := time.Now().Unix()
	params["create_time"] = nowTime
	params["update_time"] = nowTime
	params["last_heartbeat_time"] = nowTime
	stmt, args, err := s.db.BindNamed(stmtInsert, params)
	if err != nil {
		return err
	}
	res, err := s.db.Exec(stmt, args...)
	if err != nil {
		return err
	}
	agentIDInt, err := res.LastInsertId()

	host.AgentID = string(agentIDInt)

	return err
}

const stmtInsert = `
INSERT INTO host (
host_name
,host_UUID
,os
,platform
,platform_family
,platform_version
,kernel_version
,virtualization_system
,virtualization_role
,interfaces_json
,create_time
,update_time
,last_heartbeat_time
) VALUES (
:host_name
,:host_UUID
,:os
,:platform
,:platform_family
,:platform_version
,:kernel_version
,:virtualization_system
,:virtualization_role
,:interfaces_json
,:create_time
,:update_time
,:last_heartbeat_time
)
`
const queryBase = `
SELECT
agent_id
,host_name
,host_UUID
,os
,platform
,platform_family
,platform_version
,kernel_version
,virtualization_system
,virtualization_role
,interfaces_json
,create_time
,update_time
,last_heartbeat_time
`

const queryKey = queryBase + `
FROM host
WHERE agent_id = :agent_id
`
