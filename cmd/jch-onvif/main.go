package main

import (
	"flag"
	"fmt"
	"jch-onvif/internal/action"
	"jch-onvif/internal/util"
	"strings"
)

var netInterface string
var targetAction string
var targetIP string
var targetRange string
var serviceUrl string
var authType string
var authUsername string
var authPassword string

func main() {
	handlers := []action.Handler{
		action.DiscoverHandler, action.DiscoverDirectHandler, action.DiscoverRangeHandler,
		action.GetServicesHandler, action.GetServiceCapabilitiesHandler,
		action.GetHostnameHandler, action.SetHostnameHandler,
		action.GetDNSHandler,
		action.GetNTPHandler, action.SetNTPHandler,
		action.GetNetworkInterfacesHander, action.SetNetworkInterfacesHandler,
		action.GetNetworkProtocolsHandler, action.GetNetworkDefaultGatewayHandler,
		action.GetDeviceInformationHandler,
		action.GetSystemDateTimeHandler, action.SetSystemDateTimeHandler,
		action.SystemRebootHandler,
		action.GetDiscoveryModeHandler, action.SetDiscoveryModeHandler,
		action.GetRemoteDiscoveryModeHandler, action.SetRemoteDiscoveryModeHandler,
		action.GetUsersHandler, action.GetRemoteUserHandler, action.CreateUsersHandler, action.DeleteUsersHandler,
		action.GetProfilesHandler, action.GetStreamUriHandler,
		action.GetSnapshotUriHandler,
	}
	var actions []string
	for _, handler := range handlers {
		actions = append(actions, handler.ActionKey)
	}
	flag.StringVar(&netInterface, "i", util.GetDefaultNetworkInterface(), "Network interface to use")
	flag.StringVar(&targetAction, "a", "discover", "Action to perform: "+strings.Join(actions, ", "))
	flag.StringVar(&targetIP, "ip", "127.0.0.1", "Target IP address")
	flag.StringVar(&targetRange, "ipRange", "192.168.0.0/24", "Target IP range")
	flag.StringVar(&serviceUrl, "serviceUrl", "http://127.0.0.0/device_service", "SOAP service URL for target device")
	flag.StringVar(&authType, "authType", "", "Set to http-digest for HTTP based authentication")
	flag.StringVar(&authUsername, "authUsername", "", "Username for authentication")
	flag.StringVar(&authPassword, "authPassword", "", "Password for authentication")
	flag.Parse()
	arguments := action.Arguments{
		NetInterface: netInterface,
		TargetIP:     targetIP,
		TargetRange:  targetRange,
		ServiceUrl:   serviceUrl,
		Args:         flag.Args(),
	}
	executed := false
	if authType == "http-digest" {
		err := util.InitHttpDigestAuthentication(authUsername, authPassword, serviceUrl)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	for _, handler := range handlers {
		if handler.ActionKey == targetAction {
			err := handler.Action(arguments)
			if err != nil {
				fmt.Println("Error:", err)
			}
			executed = true
		}
	}
	if !executed {
		fmt.Println("Invalid action:", targetAction)
	}
}
