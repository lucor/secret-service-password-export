# Secret Service Password Export

A Go CLI tool to export passwords from a [Secret Service](https://specifications.freedesktop.org/secret-service/latest/) keychain collection using D-Bus

It uses under the hood the following libraries that are the only dependencies:
- [go-dbus-keyring](https://github.com/ppacher/go-dbus-keyring) used to query the keyring application via D-Bus
- [godbus/dbus](https://github.com/godbus/dbus) library that implements the D-Bus message protocol

## Installation

```
go install github.com/lucor/secret-service-password-export@latest
```

## Usage

```
Usage:
	secret-service-export [collection]
	
Options:
	-c, --collection    Collection to export. Leave empty to list the available collections
	-f, --format	    Output format for the export. Allowed values: [paw, csv]. Default to Paw JSON format
	-o, --output	    Write the output to the specified file. If omitted, writes to stdout
	-h, --help	    Displays the help and exit

Export the Secrect Service collection using the specified format to stdout.
```

## License

This software is available under the BSD 3-Clause License; see the [LICENSE](/LICENSE) file for the full text.
