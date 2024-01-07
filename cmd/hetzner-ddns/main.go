package main

import (
	"matthias-kutz.com/hetzner-ddns/pkg/conf"
	"matthias-kutz.com/hetzner-ddns/pkg/ddns"
	"matthias-kutz.com/hetzner-ddns/pkg/dns"
	"matthias-kutz.com/hetzner-ddns/pkg/ip"
)

func main() {
	dynDnsConf := conf.Read()

	dnsProvider := dns.Hetzner{
		ApiToken: dynDnsConf.DnsConf.ApiToken,
	}

	var ipProvider ip.Provider

	if dynDnsConf.ProviderConf.FritzBoxAddress != "" {
		ipProvider = ip.FritzBox{
			IpVersion:       dynDnsConf.ProviderConf.IpVersion,
			FritzBoxAddress: dynDnsConf.ProviderConf.FritzBoxAddress,
		}
	} else {
		ipProvider = ip.Ipify{
			IpVersion: dynDnsConf.ProviderConf.IpVersion,
		}
	}

	ddnsParameter := ddns.Parameter{
		ZoneName:   dynDnsConf.DnsConf.ZoneName,
		RecordName: dynDnsConf.RecordConf.RecordName,
		RecordType: dynDnsConf.RecordConf.RecordType,
		TTL:        dynDnsConf.RecordConf.TTL,
	}

	ddnsService := ddns.Service{
		DnsProvider: dnsProvider,
		IpProvider:  ipProvider,
		Parameter:   ddnsParameter,
	}

	ddnsScheduler := ddns.Scheduler{
		CronExpression: dynDnsConf.CronConf.CronExpression,
		Service:        ddnsService,
	}

	ddnsScheduler.Start()
}
