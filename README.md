## jch-onvif

`jch-onvif` is a CLI tool to interact with security cameras that support ONVIF protocol.

### Build

Run the following command:

```
go build ./cmd/jch-onvif
```

It'll generate `jch-onvif` executable file in current directory.

### Getting Started

Perform ONVIF device discovery through multicast address (WS-Discovery) by running the following command:

```
jch-onvif -i eth0
```

It'll print response from all ONVIF devices in the same network segment.  Copy the `XAddrs` value for a device.  This is
the WSDL URL that should be entered when performing ONVIF actions in the device. 

Run the following command to get existing media profiles in an ONVIF device (use `XAddrs` output from the previous command for `serviceUrl`):

```
jch-onvif -i eth0 -serviceUrl http://192.168.1.10:1234/onvif/device_service -a get-profiles
```

Copy the value of token for one of the profile and use it as argument for the following command:

```
jch-onvif -i eth0 -serviceUrl http://192.168.1.10:1234/onvif/device_service -a get-stream-uri token0
```

If VLC media player is installed, this command will launch it and play video stream from the profile.  Otherwise, copy the value
of `Uri` and open it in any stream player.

To get a snapshot image from the stream, run the following command:

```
jch-onvif -i eth0 -serviceUrl http://192.168.1.10:1234/onvif/device_service -a get-snapshot-uri token0
```

It will launch browser and display an image that represents the snapshot of video stream at the time when the command is executed.

### Use Cases

#### Searching For Exposed ONVIF Devices

Multicast discovery only works in local network segment.  To scan for ONVIF devices one by one in IP address range, run the following command:

```
jch-onvif -i eth0 -a discover-range -ipRange 10.20.30.0/24
```

To determine if an IP address responds to ONVIF probe (either local or public), run the following command:

```
jch-onvif -i eth0 -a discover-direct -ip 10.20.30.40
```

#### Device Remote Administration From CLI

For example, to get the date time in an ONVIF device, run the following command:

```
jch-onvif -i eth0 -serviceUrl http://192.168.1.10:1234/onvif/device_service -a get-system-date-time 
```

To update the time, run the following command (note that the time value must be in UTC):

```
jch-onvif -i eth0 -serviceUrl http://192.168.1.10:1234/onvif/device_service -a set-system-date-time Manual false WIB-7 2023-08-18T09:00:00
```

To use current date time as the date time for ONVIF device, run the following command:

```
jch-onvif -i eth0 -serviceUrl http://192.168.1.10:1234/onvif/device_service -a set-system-date-time Manual false WIB-7 now
```

To change the static IP address for an ONVIF device, run the following command (note that this command will disable DHCP for the device):

```
jch-onvif -i eth0 -serviceUrl http://192.168.1.10:1234/onvif/device_service -a set-network-ipv4 eth0 192.168.1.11
```

