Results:

```
Sequential signing took: 72.3618ms
Concurrent signing took: 64.4503ms
Redis cache retrieval/signing took: 121.6283ms
```

The issue is that concurrent signing is not any faster than sequential signing, despite being more complex and taking up more resources.

Using redis as a cache is slower if there are more than a few cache-misses, because it has to sign the missing URLs and then store them in cache. The initial fetch is very fast, so if there are no cache-misses it could be as long as 1ms. But redis adds an external service + complexity as well.

