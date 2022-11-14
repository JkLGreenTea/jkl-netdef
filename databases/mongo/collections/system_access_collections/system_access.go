package system_access_collections

type SystemAccess struct {
	HttpRequests *SysHttpRequests
	Modules      *SysModules
	Tokens       *SysTokens
	Roles        *SysRoles
}
