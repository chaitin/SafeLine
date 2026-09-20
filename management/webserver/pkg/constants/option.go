package constants

const (
	SecretKey         = "secret_key"
	MachineID         = "machine_id"
	PolicyGroupGlobal = "policy_group_global"
	SrcIPConfig       = "src_ip_config"

	// BootstrapOpenedAt records when the anonymous TFA bootstrap of the
	// account was opened, which is what bounds how long a caller may claim the
	// TFA secret without a credential (see api.GetOTPUrl). It is written when
	// the installation is set up and by `mgt -reset_user`, and never by a
	// restart.
	BootstrapOpenedAt = "tfa_bootstrap_opened_at"
)
