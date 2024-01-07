package conf

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

const (
	EnvZoneName        = "ZONE_NAME"
	EnvApiToken        = "API_TOKEN"
	EnvRecordType      = "RECORD_TYPE"
	EnvRecordName      = "RECORD_NAME"
	EnvCronExpression  = "CRON_EXPRESSION"
	EnvTimeToLive      = "TTL"
	EnvFritzboxAddress = "FRITZ_BOX_ADDRESS"

	DescZoneName        = "The DNS zone that DDNS updates should be applied to."
	DescApiToken        = "Your Hetzner API token."
	DescRecordType      = "The record type of your zone. If your zone uses an IPv4 address use `A`. Use `AAAA` if it uses an IPv6 address."
	DescRecordName      = "The name of the DNS-record that DDNS updates should be applied to. This could be `sub` if you like to update the subdomain `sub.example.com` of `example.com`. The default value is `@`"
	DescCronExpression  = "The cron expression of the DDNS update interval. The default is every 5 minutes - `*/5 * * * *`"
	DescTimeToLive      = "Time to live of the record"
	DescFritzboxAddress = "If set it will determine the IP from this fritzbox. Can be a domain name (like `fritz.box`) or an IP address (like `192.168.178.1`). If not set it will determine IP from some online service."

	DefaultRecordName     = "@"
	DefaultCronExpression = "*/5 * * * *"
	DefaultTimeToLive     = 86400

	IPv4           = "IPv4"
	IPv6           = "IPv6"
	IPv6RecordType = "AAAA"
)

type DynDnsConf struct {
	DnsConf      DnsConf
	RecordConf   RecordConf
	ProviderConf ProviderConf
	CronConf     CronConf
}

type DnsConf struct {
	ApiToken string
	ZoneName string
}

type RecordConf struct {
	RecordType string
	RecordName string
	TTL        int
}

type CronConf struct {
	CronExpression string
}

type ProviderConf struct {
	IpVersion       string
	FritzBoxAddress string
}

type ArgumentMissingError struct {
	argumentName string
}

func (e *ArgumentMissingError) Error() string {
	return "The mandatory argument " + e.argumentName + " is missing"
}

func Read() DynDnsConf {
	// Mandatory flags
	var zoneName, apiToken, recordType string
	flag.StringVar(&zoneName, EnvZoneName, zoneName, DescZoneName)
	flag.StringVar(&apiToken, EnvApiToken, apiToken, DescApiToken)
	flag.StringVar(&recordType, EnvRecordType, recordType, DescRecordType)

	// Optional flags
	var recordName = DefaultRecordName
	flag.StringVar(&recordName, EnvRecordName, recordName, DescRecordName)
	var cronExpression = DefaultCronExpression
	flag.StringVar(&cronExpression, EnvCronExpression, cronExpression, DescCronExpression)
	var ttl = DefaultTimeToLive
	flag.IntVar(&ttl, EnvTimeToLive, ttl, DescTimeToLive)
	var fritzboxAddress string
	flag.StringVar(&fritzboxAddress, EnvFritzboxAddress, fritzboxAddress, DescFritzboxAddress)

	// Parse flags
	flag.Parse()

	// Computed confs
	var ipVersion = IPv4
	if recordType == IPv6RecordType {
		ipVersion = IPv6
	}

	dynDnsConf := DynDnsConf{
		DnsConf: DnsConf{ApiToken: apiToken, ZoneName: zoneName},
		RecordConf: RecordConf{
			RecordType: recordType,
			RecordName: recordName,
			TTL:        ttl,
		},
		ProviderConf: ProviderConf{
			IpVersion:       ipVersion,
			FritzBoxAddress: fritzboxAddress,
		},
		CronConf: CronConf{CronExpression: cronExpression},
	}

	validatedConf, err := validate(dynDnsConf)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return validatedConf
}

func PrintUsage() {
	flag.Usage()
}

func validate(dynDnsConf DynDnsConf) (DynDnsConf, error) {
	// Check api token
	if dynDnsConf.DnsConf.ApiToken == "" {
		return dynDnsConf, &ArgumentMissingError{
			argumentName: EnvApiToken,
		}
	}

	// Check zone name
	if dynDnsConf.DnsConf.ZoneName == "" {
		return dynDnsConf, &ArgumentMissingError{
			argumentName: EnvZoneName,
		}
	}

	// Check record type
	if dynDnsConf.RecordConf.RecordType == "" {
		return dynDnsConf, &ArgumentMissingError{
			argumentName: EnvRecordType,
		}
	}

	return dynDnsConf, nil
}
