# ipcdn
Check which cdn provider your IPs belong to

# Description

This tool is based on CIDR ranges collected by ProjectDiscovery on cdncheck project. This is a cli implementation for easy use. Also it contains a GitHub action to download everyday the json containing CIDRs.

The tool reads from stdin the IP and (by default) it prints the IPs belonging to any CDN provider.

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