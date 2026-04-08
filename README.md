# SubPipe 🚀

> **Detect DNS Vulnerabilities With Absolute Accuracy.**

A high-octane sniper engine designed for bug bounty hunters. Stop chasing theoretical noise; get verified takeovers, dangling cloud IPs, and SSRF enablers with zero false positives. 

SubPipe acts as a lightning-fast CLI client that pipes your reconnaissance data directly into our concurrent cloud engines, streaming actionable vulnerabilities back to your terminal in real-time.

## ✨ Core Features

* **Zero False Positives:** Actively suppresses noise by verifying live hosts before flagging.
* **Real-Time Streaming:** Built on Server-Sent Events (SSE) to deliver findings the millisecond they are discovered—no waiting for batch jobs to finish.
* **Cloud IP Engine:** Identifies dangling Elastic IPs (AWS, GCP), Azure takeovers, and internal RFC1918 IP exposures (SSRF enablers).
* **Registrar Engine:** Bypasses standard resolvers using raw API checks to catch dangling CNAMEs, expired MX records, and critical NS takeovers.
* **PaaS Takeovers:** Natively detects unprotected Surge, ElasticBeanstalk, and other managed service endpoints.
* **Smart Deduplication:** Automatically strips messy URLs, ports, and duplicate targets to save your quota.

## 📥 Installation

Download the latest pre-compiled binary for your OS (Linux, macOS, Windows) from the [Releases page](https://github.com/subpipe/subpipe/releases).

**Linux/macOS Setup:**
```bash
chmod +x subpipe
mv subpipe /usr/local/bin/
