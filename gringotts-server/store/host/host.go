package host

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/jinlingan/gringotts/gringotts-server/store/db"

	"github.com/jinlingan/gringotts/gringotts-server/model"
	"github.com/jmoiron/sqlx"
)

// New returns a new UserStore.
func New(db *sqlx.DB) model.HostStore {
	return &hostStore{db}
}

type hostStore struct {
	db *sqlx.DB
}

func (h *hostStore) UpdateAgentInfo(host *model.Host) error {
	params := toParams(host)
	stmt, args, err := h.db.BindNamed(stmtUpdateHostInfo, params)
	if err != nil {
		return nil
	}
	_, err = h.db.Exec(stmt, args...)
	return err
}

func (h *hostStore) GetConfigVersion(agentID string) (int64, error) {
	var configVersion int64
	params := map[string]interface{}{
		"agent_id": agentID,
	}
	stmt, args, err := h.db.BindNamed(stmtGetConfigVersion, params)
	if err != nil {
		return configVersion, err
	}

	err = h.db.QueryRow(stmt, args...).Scan(&configVersion)
	if err != nil {
		return configVersion, err
	}
	return configVersion, nil
}

func (h *hostStore) UpdateConfigVersion(agentID string) error {
	params := map[string]interface{}{
		"agent_id": agentID,
	}
	stmt, args, err := h.db.BindNamed(stmtUpdateConfigVersion, params)
	if err != nil {
		return err
	}
	_, err = h.db.Exec(stmt, args...)
	return err
}

func (h *hostStore) Exist(agentID string) (bool, error) {
	count := 0
	params := map[string]interface{}{
		"agent_id": agentID,
	}
	stmt, args, err := h.db.BindNamed(queryCountByID, params)
	if err != nil {
		return false, err
	}

	err = h.db.QueryRow(stmt, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil

}

func (h *hostStore) UpdateHeartBeatTime(agentID string, time int64) error {
	params := map[string]interface{}{
		"agent_id":            agentID,
		"last_heartbeat_time": time,
	}
	stmt, args, err := h.db.BindNamed(updateHeartBeatTime, params)
	if err != nil {
		return err
	}
	_, err = h.db.Exec(stmt, args...)
	return err
}

func (h *hostStore) Find(agentID string) (*model.Host, error) {
	if agentID == "" {
		return nil, nil
	}
	out := &model.Host{AgentID: agentID}
	params := toParams(out)
	query, args, err := h.db.BindNamed(queryKey, params)
	if err != nil {
		return nil, err
	}
	row := h.db.QueryRow(query, args...)
	err = scanRow(row, out)
	if err == sql.ErrNoRows {
		return nil, nil
	}
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
		&dest.ConfigVersion,
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
		"config_version":        h.ConfigVersion,
	}
}
func (h *hostStore) Create(host *model.Host) (agentID string, err error) {

	params := toParams(host)
	nowTime := time.Now().Unix()
	params["create_time"] = nowTime
	params["update_time"] = nowTime
	//params["last_heartbeat_time"] = nowTime
	stmt, args, err := h.db.BindNamed(stmtInsert, params)
	if err != nil {
		return "", err
	}
	res, err := h.db.Exec(stmt, args...)
	if err != nil {
		return "", err
	}
	agentIDInt, err := res.LastInsertId()

	return strconv.FormatInt(agentIDInt, 10), err
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
,config_version
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
,:config_version
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
,config_version
`

const stmtUpdateHostInfo = `
UPDATE host SET
host_name = :host_name,
host_UUID = :host_UUID,
os = :os,
platform = :platform,
platform_family = :platform_family,
platform_version = :platform_version,
kernel_version = :kernel_version,
virtualization_system = :virtualization_system,
virtualization_role = :virtualization_role,
interfaces_json = :interfaces_json
WHERE agent_id = :agent_id
`

const queryKey = queryBase + `
FROM host
WHERE agent_id = :agent_id
`

const updateHeartBeatTime = `
UPDATE host SET
last_heartbeat_time = :last_heartbeat_time
WHERE agent_id = :agent_id
`

const queryCountByID = `
SELECT COUNT(*)
FROM host
WHERE agent_id = :agent_id
`

const stmtUpdateConfigVersion = `
UPDATE host SET
config_version = config_version + 1
WHERE agent_id = :agent_id
`

const stmtGetConfigVersion = `
SELECT config_version
FROM host
WHERE agent_id = :agent_id
`
