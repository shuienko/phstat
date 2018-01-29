# phstat
Get [Pi-Hole](https://github.com/pi-hole/pi-hole) metrics remotely using command line.

## download

Get binary from [Releases](https://github.com/shuienko/phstat/releases) and save it somewhere in your `$PATH`

## or build yourself
```bash
go get github.com/shuienko/phstat/gohole
go build -o phstat
```

## set environment variables
* get `PIHOLE_TOKEN` here http://your-pi-hole-ip/admin/settings.php?tab=api
* add to your `.bashrc` or `.zshrc`:

```bash
export PIHOLE_HOST=your-pi-hole-ip
export PIHOLE_TOKEN=longtokenstring
```


## use

```bash
Usage: phstat [-n NUMBER] summary|blocked|queries|clients
  -n number
    	number of returned entries (default 10)
```

# example

```bash
$ phstat -n 5 clients
=== Clients over last 24h:
- 192.168.1.72 : 7231
- 192.168.1.46 : 6359
- 192.168.1.60 : 1667
- 192.168.1.67 : 721
- 192.168.1.61 : 685
```
