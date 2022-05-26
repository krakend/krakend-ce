---
name: Help and questions
about: You are stuck trying to do something, you get unexpected results or you simply have
  a question or suggestion
title: ''
labels: 'question'
assignees: ''

---
**Version of KrakenD you are using**
Get it with: `krakend help | grep Version`

**Describe what are you trying to do**
A clear and concise description of what you want to do and how is the setup with KrakenD.

**Your configuration file**
The content of your `krakend.json`. When using the flexible configuration option, the computed file can be generated using `FC_OUT=out.json`

```
{
  "version": 3,
  ...
}
```

**Configuration check output**
Ater running `krakend check -dtc krakend.json --lint` I got this output:

```
Output of the linter here.
```

**Commands used**
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

**Logs**
Logs you saw in the console and debugging information
