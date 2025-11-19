package exitcodes

const (
	NoError                       = 0
	UnknownError                  = 1
	DockerCommandNotFound         = 2
	DockerInsufficientPermissions = 3
	ConfigFileNotFound            = 4
	ServiceNotFound               = 5
	NoServicesFound               = 6
	MasterKeyStatError            = 7
	MasterKeyNotFound             = 8
)
