# ipcdn
Check which CDN providers an IP list belongs to

# Description

This tool is based on the CIDR ranges collected by ProjectDiscovery in their [cdncheck](https://github.com/projectdiscovery/cdncheck) project. This is a cli implementation for easy use and it also contains a [GitHub action](https://github.com/six2dez/ipcdn/actions/workflows/download.yml) to download everyday the json containing CIDRs.

> **Its main use case is to avoid port scanning and/or other intensive tasks on IPs hosted outside our target and that can also result in blocking or banning.**

The tool reads from stdin the list of IPs and (by default) prints the IPs belonging to any CDN provider. Using the flags you can specify to print the ones that do not belong to CDNs or to print the detail of which provider it uses.

It detects the following CDNs: Azure, Cloudflare, Cloudfront, Fastly and Incapsula.

# Install

`go install -v github.com/six2dez/ipcdn@latest`

# Modes

- m (method):
  - cdn - prints only IPs on CDN ranges
  - not - prints only IPs NOT on CDN ranges
  - all - prints all IPs with description (verbose)

- verbose - prints description of which CDN provider is using

# Usage

```
> ./ipcdn -h
Usage of ./ipcdn:
  -m string
        Output method:
            cdn - Only prints IPs on CDN
            not - only prints NOT on CDN
            all - prints all with verbose
         (default "cdn")
  -v    verbose mode
```