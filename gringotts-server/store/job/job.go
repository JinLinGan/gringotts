package job

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/jinlingan/gringotts/gringotts-server/store/db"

	"github.com/jinlingan/gringotts/gringotts-server/model"
	"github.com/jmoiron/sqlx"
)

// New returns a new UserStore.
func New(db *sqlx.DB) model.JobStore {
	return &jobStore{db}
}

type jobStore struct {
	db *sqlx.DB
}

func (j *jobStore) CreateJob(job *model.Job) (jobID string, err error) {

	params := toParams(job)
	nowTime := time.Now().Unix()
	params["create_time"] = nowTime
	params["update_time"] = nowTime
	//params["last_heartbeat_time"] = nowTime
	stmt, args, err := j.db.BindNamed(stmtInsert, params)
	if err != nil {
		return "", err
	}
	res, err := j.db.Exec(stmt, args...)
	if err != nil {
		return "", err
	}
	jobIDInt, err := res.LastInsertId()

	return strconv.FormatInt(jobIDInt, 10), err
}

func (j *jobStore) DeleteJob(jobID string) error {
	params := map[string]interface{}{
		"job_id": jobID,
	}
	stmt, args, err := j.db.BindNamed(stmtDelete, params)
	if err != nil {
		return err
	}
	_, err = j.db.Exec(stmt, args...)

	return err
}

func (j *jobStore) UpdateJobConfig(job *model.Job) error {
	params := map[string]interface{}{
		"running_interval": job.RunningInterval,
		"config":           job.Config,
		"update_time":      time.Now().Unix(),
		"job_id":           job.JobID,
	}

	stmt, args, err := j.db.BindNamed(stmtUpdateConfig, params)
	if err != nil {
		return err
	}
	_, err = j.db.Exec(stmt, args...)
	return err
}

func (j *jobStore) GetJobs(agentID string) ([]*model.Job, error) {

	params := map[string]interface{}{
		"agent_id": agentID,
	}
	stmt, args, err := j.db.BindNamed(queryByAgent, params)
	if err != nil {
		return nil, err
	}
	rows, err := j.db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	out, err := scanRows(rows)
	return out, err

}

func (j *jobStore) UpdateJobRunningState(job *model.Job) error {
	params := map[string]interface{}{
		"running_state":     job.RunningState,
		"error_msg":         job.ErrorMsg,
		"last_running_time": job.LastRunningTime,
		"last_report_time":  time.Now().Unix(),
	}

	stmt, args, err := j.db.BindNamed(stmtUpdateRunningState, params)
	if err != nil {
		return err
	}
	_, err = j.db.Exec(stmt, args...)
	return err
}

func scanRow(scanner db.Scanner, dest *model.Job) error {
	err := scanner.Scan(
		&dest.JobID,
		&dest.AgentID,
		&dest.RunnerType,
		&dest.RunnerModule,
		&dest.ModuleVersion,
		&dest.RunningInterval,
		&dest.Config,
		&dest.CreateTime,
		&dest.UpdateTime,
		&dest.RunningState,
		&dest.ErrorMsg,
		&dest.LastRunningTime,
		&dest.LastReportTime,
	)
	return err
}
func scanRows(rows *sql.Rows) ([]*model.Job, error) {
	defer rows.Close()

	builds := []*model.Job{}
	for rows.Next() {
		build := new(model.Job)
		err := scanRow(rows, build)
		if err != nil {
			return nil, err
		}
		builds = append(builds, build)
	}
	return builds, nil
}

func toParams(j *model.Job) map[string]interface{} {
	return map[string]interface{}{
		"job_id":            j.JobID,
		"agent_id":          j.AgentID,
		"runner_type":       j.RunnerType,
		"runner_module":     j.RunnerModule,
		"module_version":    j.ModuleVersion,
		"running_interval":  j.RunningInterval,
		"config":            j.Config,
		"create_time":       j.CreateTime,
		"update_time":       j.UpdateTime,
		"running_state":     j.RunningState,
		"error_msg":         j.ErrorMsg,
		"last_running_time": j.LastRunningTime,
		"last_report_time":  j.LastReportTime,
	}
}

const stmtUpdateConfig = `
UPDATE job SET
update_time = :update_time,
running_interval = :running_interval,
config = :interval
WHERE job_id = :job_id
`

const stmtUpdateRunningState = `
UPDATE job SET
running_state = :running_state     
error_msg = :error_msg         
last_running_time = :last_running_time 
last_report_time = :last_report_time  
WHERE job_id = :job_id
`

const stmtInsert = `
INSERT INTO job (
agent_id
,runner_type
,runner_module
,module_version
,running_interval
,config
,create_time
,update_time
,running_state
,error_msg
,last_running_time
,last_report_time
) VALUES (
:agent_id
,:runner_type
,:runner_module
,:module_version
,:running_interval
,:config
,:create_time
,:update_time
,:running_state
,:error_msg
,:last_running_time
,:last_report_time
)
`
const queryBase = `
SELECT
job_id
,agent_id
,runner_type
,runner_module
,module_version
,running_interval
,config
,create_time
,update_time
,running_state
,error_msg
,last_running_time
,last_report_time
`
const stmtDelete = `
DELETE FROM job
WHERE job_id = :job_id
`
const queryByAgent = queryBase + `
FROM job
WHERE agent_id = :agent_id
`
