# StatsD HTTP Proxy

StatsD HTTP proxy with REST interface for using in browsers.


Sample code to send metric in browser with JWT token in header:

```javascript
$.ajax({
    url: 'http://127.0.0.1:8080/count/some.key.name',
    method: 'POST',
    headers: {
        'X-JWT-Token': 'some-jwt-token'
    },
    data: {
        value: 100500
    }
});
```
## Supported metrics

For the general reference see https://www.librato.com/docs/kb/collect/collection_agents/stastd/#

All metrics accept `tags` as comma-separated key=value pairs (InfluxDB tag format):

```javascript
data: {
    value: 100500,
    tags: 'env=prod,locale=en-us'
}
```

### `count`

Adds count to the bucket. Expected `n` as integer. By default `n` is 0.

### `incr`

Increments the given bucket. It is equivalent to count with `n` default to 1.

### `gauge`

Sets the gauge metric. Expected `value` as integer. Before setting negative gauge, it needs to be set to `0`.

### `timing`

Adds timing to the bucket. Expected `dur` as milliseconds integer. Default is `0`.

### `uniq`

Adds unique value in a set bucket. Expected `value` as string. Sets are a relatively new concept in recent versions of StatsD. Sets track the number of unique elements belonging to a group. At each flush interval, the statsd backend will push the number of unique elements in the set as a single gauge value.
