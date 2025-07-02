go run . -mode=server -addr=:8080

go run . -mode=bench -server=http://localhost:8080 -runs=5

---

*Results*

Run 1   cdn: 435.6474ms raw: 87.3505ms (Δ=-348.2969ms)  private: 382.042ms (Δ=-53.6054ms)
Run 2   cdn: 168.7811ms raw: 50.1827ms (Δ=-118.5984ms)  private: 396.9195ms (Δ=228.1384ms)
Run 3   cdn: 165.5404ms raw: 77.6181ms (Δ=-87.9223ms)   private: 582.9899ms (Δ=417.4495ms)
Run 4   cdn: 149.688ms  raw: 38.1562ms (Δ=-111.5318ms)  private: 365.3579ms (Δ=215.6699ms)
Run 5   cdn: 395.4185ms raw: 44.2521ms (Δ=-351.1664ms)  private: 323.957ms (Δ=-71.4615ms)

The raw fetch (no CDN) is much faster. My content is not being cached; I need to add some sort of cache header to the metadata of the object in s3 in order to get cloudflare to cache it (currently it's returning cache as 'DYNAMIC'). If that doesn't work, we should simply bypass cloudflare altogether

Note also that public buckets, even with cloudflare, are much faster than private buckets. On the order of about 200ms faster. I don't know if there's a way to speed this up.
