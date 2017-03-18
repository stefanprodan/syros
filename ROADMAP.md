# Syros Roadmap

The roadmap represents our estimate of when each feature will enter full production release. 

### Near-term

***Health checks***

* Collect the health check results using the Docker collector and/or develop a new collector to extract heath checks from Consul
* Display the health check result on the container page, add a health column to containers tables
* Track health status changes into db container health history table
* Create a container health history page with a line chart and health check results table view

### Medium-term

***PagerDuty integration***

* Collect incidents details from PagerDuty
* Create On-Call dashboard 
* Display incidents stats (opened, resolved last 24h, resolved last 30 days)
* Display incidents history (table view)
* On-Call handover submit form (save to db and send by mail to the On-Call team)
* Display On-Call handover history (table view)
* Generate incidents monthly reports

***ELK integration***

* Collect from Elasticseach logs statistics per container (INFO, WARN, ERROR) 
* Display the logs stats in the container table view (compose Kibana hyperlinks)
* Add the logs stats to the environment dashboard
* Add the logs stats and chart to the container health history page or create a dedicated page

### Long-Term

***Prometheus integration***

* Render host graphs (CPU, Memory, IO, Disk, Network) 
* Render container graphs (CPU, Memory, IO, Network)

***Docker Registry integration***

* Collect registry images
* Display images (table view) and render images deploy graph
* Track used/unused images by linking to the running container 