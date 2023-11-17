package engine

type Engine string

const (
	InnoDB    Engine = "InnoDB"
	MyISAM    Engine = "MyISAM"
	Memory    Engine = "Memory"
	Archive   Engine = "Archive"
	CSV       Engine = "CSV"
	Blackhole Engine = "Blackhole"
	Federated Engine = "Federated"
)
