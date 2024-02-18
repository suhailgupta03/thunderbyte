package common

// GetService It returns the service from the injected service map
func GetService(serviceName ServiceName, serviceMap *InjectedServicesMap) (interface{}, bool) {
	service, ok := (*serviceMap)[serviceName]
	return service, ok
}
