package builtin

import (
	"encoding/json"
	"github.com/Jumpscale/agent2/agent/lib/pm"
	"github.com/shirou/gopsutil/net"
)

const (
	CmdGetNicInfo = "get_nic_info"
)

func init() {
	pm.CMD_MAP[CmdGetNicInfo] = InternalProcessFactory(getNicInfo)
}

func getNicInfo(cmd *pm.Cmd, cfg pm.RunCfg) *pm.JobResult {
	result := pm.NewBasicJobResult(cmd)
	result.Level = pm.L_RESULT_JSON

	info, err := net.NetInterfaces()

	if err != nil {
		result.State = pm.S_ERROR
		m, _ := json.Marshal(err)
		result.Data = string(m)
	} else {
		result.State = pm.S_SUCCESS
		m, _ := json.Marshal(info)

		result.Data = string(m)
	}

	return result
}
