[<img src="https://monitoringartist.github.io/managed-by-monitoringartist.png" alt="Managed by Monitoring Artist: DevOps / Docker / Kubernetes / AWS ECS / Zabbix / Zenoss / Terraform / Monitoring" align="right"/>](http://www.monitoringartist.com 'DevOps / Docker / Kubernetes / AWS ECS / Zabbix / Zenoss / Terraform / Monitoring')

# Zabbix template XML/JSON/YAML converter

[![GPL license](https://img.shields.io/badge/license-GPL-brightgreen.svg)](https://opensource.org/licenses/gpl-license)

```
$ go get -u github.com/monitoringartist/zabbix-template-converter
$ zabbix-template-converter --in in-template.[xml|json|yaml] --out out-template.[xml|json|yaml]
```

App depends heavily on the [Zabbix template model](https://github.com/monitoringartist/go-zabbix), which supports:
- Zabbix 4.4 - the majority is supported = all officials 4.4 templates are supported
- Zabbix 4.0 - only partial support at the moment

You may verify support for your template with commands, where you compare input with output xml templates:
```
$ zabbix-template-converter --in <in-template.xml> --out <out-template.xml>
$ diff --text --ignore-blank-lines --ignore-space-change <in-template.xml> <out-template.xml>
```
No difference is OK.

Feel free to open PR, which improves used template model.

# Author

[Devops Monitoring Expert](http://www.jangaraj.com 'DevOps / Docker / Kubernetes / AWS ECS / Google GCP / Zabbix / Zenoss / Terraform / Monitoring'),
who loves monitoring systems and cutting/bleeding edge technologies: Docker,
Kubernetes, ECS, AWS, Google GCP, Terraform, Lambda, Zabbix, Grafana, Elasticsearch,
Kibana, Prometheus, Sysdig,...

Summary:
* 3000+ [GitHub](https://github.com/monitoringartist/) stars
* 100 000+ [Grafana dashboard](https://grafana.net/monitoringartist) downloads
* 1 000 000+ [Docker image](https://hub.docker.com/u/monitoringartist/) pulls

Professional devops / monitoring / consulting services:

[![Monitoring Artist](http://monitoringartist.com/img/github-monitoring-artist-logo.jpg)](http://www.monitoringartist.com 'DevOps / Docker / Kubernetes / AWS ECS / Google GCP / Zabbix / Zenoss / Terraform / Monitoring')
