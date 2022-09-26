#  Integration service for k6 
[K6](https://k6.io/) is an open-source tool widely used for load testing. With Keptn’s Quality gates, we can automate the process of evaluating the results of the test and also monitor them. Keptn already has support for JMeter, Locust, Litmus Chaos, and others. This project will enable users to use K6 for performance testing. 

## Tutorials

### Performance Testing in Keptn using K6 [Part 1]
In this tutorial we run standard K6 tests using [Job Executor Service](https://github.com/keptn-contrib/job-executor-service) of Keptn. Please find the tutorial [here](./docs/k6-jes-example/README.md).

### Using K6 Extension in JES - Prometheus [Part 2]
After running standard K6 in JES, we make use of K6 extensions for writing test metrics to extenal source. We'll use Prometheus remote write K6 extension. Please find the tutorial [here](./docs/k6-prometheus-example/README.md).

### Quality Gates Evaluation on Exported Metrics [Part 3]
Once we have K6 tests metrics written in Prometheus, we can use Keptn's Quality Gates evaluation for SLO compliance. Please find the tutorial [here](./docs/k6-prometheus-quality-gate-example/README.md).

## Custom Service State
There were two ways of implementing this integration
1. Using Job Executor Service
2. Having a Custom Keptn Service

We decided the Job Executor Service will be sufficient to deliver all the features. More information regarding this discussion can be found in this GitHub [issue](https://github.com/keptn-sandbox/k6-service/issues/20).

Altough we decided to go ahead with JES, the template for listening to `sh.keptn.event.test.triggered` CloudEvent and sending  `sh.keptn.event.test.started` and `sh.keptn.event.test.finished` CloudEvents. The service could be run using Helm Charts.

## About 
This repository contains details about my work in Google Summer of Code 2022. 

- [Project Page at GSOC 2022 website](https://summerofcode.withgoogle.com/programs/2022/projects/0xICJhw8)

## License

[MIT](https://github.com/jainammm/keptn-k6-service/blob/main/LICENSE)
