package model

// RegisterResp 注册回复
type RegisterResp struct {
	//TODO:AgentID可能需要是个 UUID

	AgentID       int64
	ConfigVersion int64
}
