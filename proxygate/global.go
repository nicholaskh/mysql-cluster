package proxygate

var (
	proxyGate *ProxyGate
)

func InitGlobal() {
	proxyGate = NewProxyGate()
}
