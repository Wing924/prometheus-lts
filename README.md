# Prometheus LTS

[![CircleCI](https://circleci.com/gh/Wing924/prometheus-lts/tree/master.svg?style=svg)](https://circleci.com/gh/Wing924/prometheus-lts/tree/master)

A distributed long term remote storage for Prometheus

## Overview

## System Design

### prom-lts writer

* implements the prometheus remote storage write interface
* convert prometheus WriteRequest into samples (a list of metrics)
* store kafka cluster with samples

### prom-lts partition

* consume kafka cluster and get samples
* partition samples into many segments
* store kafka cluster with segments

### prom-lts appender

* consume kafka cluster and get segments
* write segments into InfluxDB

### prom-lts reader

* implements the prometheus remote storage read interface 
* read samples from InfluxDB