---
name: Help and questions
about: You are stuck trying to do something, you get unexpected results, or you have a question or suggestion
title: ''
labels: 'question'
assignees: ''

---
<!--
Thank you for using KrakenD! Please spend some time to fill all the requested information in this template.

There are hundreds of questions across all our repositories. Having the proper context and detailed information will help us KrakenD maintainers to provide a faster answer. Unfortunately, we have to leave issues that we don't wholly understand or require more information for a much later processing.
-->

**Environment info:**

* KrakenD version: Run `krakend help | grep Version` and copy the output here
* System info: Run `uname -srm` or write `docker` when using containers
* Hardware specs: Number of CPUs, RAM, etc
* Backend technology: Node, PHP, Java, Go, etc.
* Additional environment information:

**Describe what are you trying to do**:
A clear and concise description of what you want to do and what is the expected result.

**Your configuration file**:
<!-- The content of your `krakend.json`. When using the flexible configuration option, the computed file can be generated specifying the env var FC_OUT=out.json -->

```json
{
  "version": 3,
  ...
}
```

**Configuration check output**:
Result of `krakend check -dtc krakend.json --lint` command

```
Output of the linter here.
```

**Commands used:**
How did you start the software?
```
#Example:
krakend run -d -c krakend.json

# Or maybe...
docker run --rm -it -v $PWD:/etc/krakend \
        -e FC_ENABLE=1 \
        -e FC_SETTINGS="/etc/krakend/config/settings" \
        -e FC_PARTIALS="/etc/krakend/config/partials" \
        -e FC_OUT=out.json \
        devopsfaith/krakend \
        run -c /etc/krakend/config/krakend.json -d
```

**Logs:**
Logs you saw in the console and debugging information

**Additional comments:**
