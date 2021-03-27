package app

type Message interface {
	Read()
	Write()
	Delete()
}

type SQSMessage struct {
	SubmissionID string `json:"submission_id"`
	DeviceID     string `json:"device_id"`
	TimeCreated  string `json:"time_created"`
	Events       Event  `json:"events"`
}

type Event struct {
	Processes          []Process           `json:"new_process"`
	NetworkConnections []NetworkConnection `json:"network_connection"`
}

type Process struct {
	Cmdl string `json:"cmdl"`
	User string `json:"user"`
}

type NetworkConnection struct {
	SourceIP        string `json:"source_ip"`
	DestinationIP   string `json:"destination_ip"`
	DestinationPort int    `json:"destination_port"`
}

type KinesisRecord struct {
	RecordID           string              `json:"id"`
	DeviceID           string              `json:"device_id"`
	Processes          []Process           `json:"new_process"`
	NetworkConnections []NetworkConnection `json:"network_connection"`
	Created            string              `json:"created"`
}
