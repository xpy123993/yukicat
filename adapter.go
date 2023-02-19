package main

import (
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/xpy123993/corenet"
)

func createDialer(dialerURL *url.URL) (func() (net.Conn, error), func() error, error) {
	switch dialerURL.Scheme {
	case "tcp":
		return func() (net.Conn, error) {
			return net.Dial("tcp", dialerURL.Host)
		}, nil, nil
	default:
		dialer := corenet.NewDialer([]string{dialerURL.String()}, corenet.WithDialerRelayTLSConfig(tunnelTLSConfig))
		return func() (net.Conn, error) {
			return dialer.Dial(dialerURL.Path)
		}, dialer.Close, nil
	}
}

func createCorenetListener(listenerURL *url.URL) (net.Listener, error) {
	opts := corenet.CreateDefaultFallbackOptions()
	opts.TLSConfig = tunnelTLSConfig
	opts.KCPConfig = corenet.DefaultKCPConfig()
	opts.QuicConfig.KeepAlivePeriod = 5 * time.Second
	adapter, err := corenet.CreateListenerFallbackURLAdapter(listenerURL.String(), listenerURL.Path, opts)
	if err != nil {
		return nil, err
	}
	adapters := []corenet.ListenerAdapter{}
	if port := listenerURL.Query().Get("port"); len(port) > 0 {
		iPort, err := strconv.ParseInt(port, 10, 32)
		if err != nil {
			return nil, err
		}
		localAdapter, err := corenet.CreateListenerTCPPortAdapter(int(iPort))
		if err != nil {
			return nil, err
		}
		adapters = append(adapters, localAdapter)
	}
	adapters = append(adapters, adapter)
	return corenet.NewMultiListener(adapters...), nil
}

func createListener(listenerURL *url.URL) (net.Listener, error) {
	switch listenerURL.Scheme {
	case "tcp":
		return net.Listen("tcp", listenerURL.Host)
	default:
		return createCorenetListener(listenerURL)
	}
}
