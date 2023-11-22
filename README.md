# Web Crawler System Design

### Design Plan
The system will involve 7 parts overall:
- Seed Url (Starting url for crawling)
  - For this project I will only be allowing a single url input
  - client will be a long running web server
- URL Frontier
  - data structure to store URLs for future downloads
  - This will ensure priority/politeness to not DDOS a website
  - Queue router to put data in queues - queue selector to select data from given queues 
  - Managed redis queue for FIFO data 
    - Each redis key with the primary host will have a queue associated with it
  - Workers will spin up to ingest data from the FIFO queue per key 
- HTML Downloader (including DNS resolution)
  - Gets IP addresses from the DNS resolver and starts downloading html content
- Content Parser
  - Parses HTML to ensure raw text is not malformed
- Content Seen?
  - Data store of MD5 hashes of html content - if this data store has the md5 hash from the parser it throws away the data 
  and continues work - if it doesn't have the hash it stores it. 
- Link extractor
  - Extracts links from HTML page 
- URL filter
  - Gets passed the links and stores URLs
  - URLs will then be stored in the URL Frontier and the whole process will continue

### Diagram
```
                If either the DNS resolver fails
                or parser log error and restart
                  ┌─────────────────────────┐
                  │                         │
                  │         ┌─────────┐     │
                  │         │DNS      │     │
                  │   ┌─────┤Resolver │     │
                  │   │     └───▲─────┘     │
                  │   │         │           │
                  │   │         │           │
┌─────────┐   ┌───▼───▼─┐   ┌───┴─────┐   ┌─┴───────┐
│         │   │         │   │         │   │         │
│Client   ├───►Frontier ├───►Html     ├───►Html     │
│         │   │         │   │Download │   │Parser   │
└─────────┘   └──▲───▲──┘   └─────────┘   └────┬────┘
                 │   │                         │
                 │   │       ┌───────┐    ┌────▼────┐  ┌─────────┐
                 │   │       ├───────┘    │         │  │         │
                 │   │       │ Data  ◄────┤Content  ├──►Link     │
                 │   │       │ Store │    │seen?    │  │extract  │
                 │   │       └───────┘    └┬────────┘  └────┬────┘
                 │   │                     │                │
                 │   └─────────────────────┘           ┌────▼────┐
                 │     If MD5 hash exists              │         │
                 │     restart to beginning            │URL      │
                 │                                     │Filter   │
                 │           ┌───────┐                 └───┬─────┘
                 │           ├───────┘                     │
                 └───────────┤Redis  ◄─────────────────────┘
                             │MQ     │   Url's are pushed to
                             └───────┘  redis MQ for processing
```

### Models
The data models for this will be incredibly simple. The queue data model will take form as a 
redis queue per host. 
```
{ "wikipedia": ["https://wikipedia.com", "https://wikipedia.com/test"] }
{ "go": ["https://pkg.go.com/net/http", "go.com"] }
```

Using this model will enable the use of grouping crawler workers only within the desired host.

For the `Seen content` data store it will simple be a SQLite DB containing MD5 hashes of all seen sites.

```
interface SeenContentModel {
  id PK int unique
  hash string
}
```

