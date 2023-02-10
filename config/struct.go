package config

type Model struct {
	BrandName string        `ini:"brand_name"`
	Prod      bool          `ini:"prod"`
	Server    ServerModel   `ini:"server"`
	Session   SessionModel  `ini:"session"`
	Database  DatabaseModel `ini:"database"`
	Log       LogModel      `ini:"log"`
	Security  SecurityModel `ini:"security"`
}

type ServerModel struct {
	HTTPAddr    string `ini:"http_addr"`
	HTTPPort    int    `ini:"http_port"`
	EnableHTTPS bool   `ini:"enable_https"`
	SSLCert     string `ini:"ssl_cert"`
	SSLKey      string `ini:"ssl_key"`
}

type SessionModel struct {
	Domain string `ini:"domain"`
	Path   string `ini:"path"`
	Name   string `ini:"name"`
	MaxAge int    `ini:"max_age"`
}

type DatabaseModel struct {
	Path string `ini:"path"`
}

type LogModel struct {
	MaskUnknown bool   `ini:"mask_unknown"`
	MaskDebug   bool   `ini:"mask_debug"`
	MaskTrace   bool   `ini:"mask_trace"`
	MaskInfo    bool   `ini:"mask_info"`
	MaskWarning bool   `ini:"mask_warning"`
	MaskError   bool   `ini:"mask_error"`
	MaskFatal   bool   `ini:"mask_fatal"`
	LogToFile   bool   `ini:"log_to_file"`
	FilePath    string `ini:"file_path"`
}

type SecurityModel struct {
	Username   string `ini:"username"`
	Password   string `ini:"password"`
	EncryptKey string `ini:"encrypt_key"`
}
