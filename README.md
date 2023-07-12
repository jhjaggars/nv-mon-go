This program fetches NVidia Graphics Card information and prints it out in influx line format for use with telegraf.

You can use it like this:

```
[[inputs.execd]]
  command = ["/path/to/nv-mon-go"]
  signal = "SIGHUP"
```
