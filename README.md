# Web Crawler System Design

### MVP 
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
