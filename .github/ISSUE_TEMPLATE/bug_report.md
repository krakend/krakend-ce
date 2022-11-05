---
name: Bug report
about: Create a report to help us improve
title: ''
labels: ''
assignees: ''

---
<!--
Thank you for reporting a bug of KrakenD. Please spend some time to fill all the requested information in this template.

Having the proper context and detailed information will help us KrakenD maintainers to investigate this issue faster. Unfortunately, we have to leave issues that we don't wholly understand or require more information for a much later processing.
-->

**Environment info:**

* KrakenD version: Run `krakend help | grep Version` and copy the output here
* System info: Run `uname -srm` or write `docker` when using containers
* Hardware specs: Number of CPUs, RAM, etc
* Backend technology: Node, PHP, Java, Go, etc.
* Additional environment information:

**Describe the bug**
A clear and concise description of what the bug is.


**Your configuration file**:
<!-- The content of your `krakend.json`. When using the flexible configuration option, the computed file can be generated specifying the env var FC_OUT=out.json -->

```json
{
  "version": 3,
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
