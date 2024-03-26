package store

var SubscribedAddresses = make(map[string]bool)

func Subscribe(address string) {
	SubscribedAddresses[address] = true
}

func Unsubscribe(address string) {
	delete(SubscribedAddresses, address)
}
