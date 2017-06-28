# Summary
Auto downsample version of prometheus. It just for some one who want to save long term data in prometheus. 

# Introduction
* Fork from prometheusv1.5.2.
* Add auto downsample feature by modify promql a little bit and storage a little bit more.
* Just run without problems for now, no union test (also may break exist test).
* Origin retention is useless now, won't drop any data, just downsample to the next retention level.
* A small problem while use with grafana(should set "null connected" if your query step set smaller than retention level's interval), some step get a zerosample because of after downsaple gap is larger and step is decided by the defined retention, should be dev a new plugin for this version prometheus.

# Usage
* Change retention in storage/local/retentions.go, see func init().
* Build and run it the same way as prometheus.

# Other
* Any problem please feel free to contact me here or Email: qinguoan2007@hotmail.com.
* Last update: 2017-06-28 with a lot of bugfix.
