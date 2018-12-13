package meta

type MessageSpec interface {
	ToByte() ([]byte, error)
	ToJson() (string, error)
}

type ChannelSpec interface {
	HasMessage(msg MessageSpec) (bool, error)
	Append(msg MessageSpec) error
	Send(msg MessageSpec) error
}
