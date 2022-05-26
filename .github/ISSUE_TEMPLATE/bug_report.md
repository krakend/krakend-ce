---
name: Bug report
about: Create a report to help us improve
title: ''
labels: ''
assignees: ''

---
**Version of KrakenD you are using**
Get it with: `krakend help | grep Version`

**Describe the bug**
A clear and concise description of what the bug is.

**Your configuration file**
The content of your `krakend.json`. When using the flexible configuration option, the computed file can be generated using `FC_OUT=out.json`
```
{
  "version": 2,
  ...
}
```
**Commands used**
How did you start the software?
```
#Example:
docker run --rm -it -v $PWD:/etc/krakend \
        -e FC_ENABLE=1 \
        -e FC_SETTINGS="/etc/krakend/config/settings" \
        -e FC_PARTIALS="/etc/krakend/config/partials" \
        -e FC_OUT=out.json \
        devopsfaith/krakend \
        run -c /etc/krakend/config/krakend.json -d
```
**Expected behavior**
A clear and concise description of what you expected to happen.

**Logs**
If applicable, any logs you saw in the console and debugging information

**Additional context**
Add any other context about the problem here.
