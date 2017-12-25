# CloudflareCTL
A tool to help manage cloudflare resources

```
NAME:
   cloudflarectl - golang tool to clear cloudflare cache

USAGE:
   cloudflarectl [global options] command [command options] [arguments...]

VERSION:
   0.9.0

AUTHOR:
   Tuan Nguyen <tuan.nguyen@honestbee.com>

COMMANDS:
     clear, c      Clear list of files's cache
     clearAll, ca  Clear everything
     status, s     Show account status
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --apiKey value                          Cloudflare's API key (REQUIRED) [$CF_API_KEY]
   --email value                           Cloudflare's email account (REQUIRED) [$CF_API_EMAIL]
   --file <files slice>, -f <files slice>  <files slice> which need to be cleared  (default: "./files_list.txt") [$CF_FILES]
   --url value                             URL to clear cache (non-scheme) [$URL_BASE]
   --domain value                          cloudflare domain [$CF_DOMAIN]
   --help, -h                              show help
   --version, -v                           print the version
```
