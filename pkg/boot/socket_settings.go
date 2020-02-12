package boot

// Settings for the websocket

type SocketSettings struct {
	*CommonSettings
	withoutApp         bool
	createAppifMissing bool
}

func (n *SocketSettings) WebSocketEngineURL() string {
	engineUrl := n.EngineURL()
	// Only modify the URL path if there is no path set
	if engineUrl.Path == "" || engineUrl.Path == "/" {
		engineUrl.Path = "/app/engineData/ttl/" + n.Ttl()
	}
	return engineUrl.String()
}

func (n *SocketSettings) Ttl() string {
	return n.GetString("ttl")
}

func (n *SocketSettings) NoData() bool {
	return n.GetBool("no-data")
}
func (n *SocketSettings) NoSave() bool {
	return n.GetBoolAllowNoFlag("no-save")
}
func (n *SocketSettings) WithoutApp() bool {
	return n.withoutApp
}
func (n *SocketSettings) CreateAppIfMissing() bool {
	return n.createAppifMissing
}
