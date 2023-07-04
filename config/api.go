package config

type API struct {
	HttpAPI        HttpAPI
	HttpJsonRpcAPI HttpJsonRpcAPI
}

func DefaultAPI() API {
	return API{
		HttpAPI:        defaultHttpAPI(),
		HttpJsonRpcAPI: defaultHttpJsonRpcApi(),
	}
}

var HttpApiVersion = "v1"

type HttpAPI struct {
	// ListenAddress is a string indicates listening address, protocol and port
	// of a http api server in multi-address manner
	ListenAddress string `yaml:"ListenAddress"`
	Version       string `yaml:"Version"`
}

func defaultHttpAPI() HttpAPI {
	return HttpAPI{
		ListenAddress: "/ip4/127.0.0.1/tcp/9000",
		Version:       HttpApiVersion,
	}
}

var HttpJsonRpcApiVersion = "v1"

type HttpJsonRpcAPI struct {
	ListenAddress string `yaml:"ListenAddress"`
	Version       string `yaml:"Version"`
}

func defaultHttpJsonRpcApi() HttpJsonRpcAPI {
	return HttpJsonRpcAPI{
		ListenAddress: "/ip4/127.0.0.1/tcp/9010",
		Version:       HttpJsonRpcApiVersion,
	}
}
