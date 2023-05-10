# **This repository is archived since 10/05/2023 in favor of https://github.com/projectdiscovery/cdncheck**

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

# Screenshots

## Default (only CDN)
![image](https://user-images.githubusercontent.com/24670991/162330965-bc984bd8-388f-4f6f-aead-d90a8eb0e63e.png)

## Default with verbose
![image](https://user-images.githubusercontent.com/24670991/162331123-3489f885-497e-4b52-9e98-766c3ba8c580.png)

## Not (not CDN)
![image](https://user-images.githubusercontent.com/24670991/162331204-8f13bb33-6dc9-471b-9a82-10cfd8437585.png)

## All (verbose)
![image](https://user-images.githubusercontent.com/24670991/162330779-0be6163c-a661-4c9c-9b43-93cb9227d342.png)
