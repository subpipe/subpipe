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

You can install SubPipe using Go (recommended) or by downloading a pre-compiled binary.

### Option 1: Install via Go (Recommended)
If you have Go installed on your system, this is the fastest way to get started. It will automatically compile the binary for your specific OS and architecture:
```bash
go install [github.com/subpipe/subpipe@latest](https://github.com/subpipe/subpipe@latest)
```

## 🛠️ Usage

SubPipe takes a list of subdomains via standard input (`stdin`). To keep your workflow clean and your credentials out of your `.bash_history`, export your API key as an environment variable before running your scans.

```bash
# 1. Set your API Key securely in your environment
export SUBPIPE_API_KEY="your-api-key-here"

# 2. Pipe your recon data into the engine
cat subdomains.txt | subpipe
```

## Example Output
SubPipe color-codes findings by severity and provides the exact context you need to immediately pivot to exploitation:

```
🚀 SubPipe Analysis Started: 26 targets sent to [https://api.subpipe.run/v1/scan](https://api.subpipe.run/v1/scan)
----------------------------------------------------------------------  LOW      Internal/RFC1918 IP Exposure (SSRF Enabler): internal-test.subpipe.run [10.0.0.5]  CRITICAL Nameserver Domain Expired: demo-ns.subpipe.run [Takeover: definitely-not-registered-subpipe-123.com]  HIGH     Mail Exchange (MX) Domain Expired: demo-mx.subpipe.run [Takeover: totally-fake-mail-server-subpipe-999.com]  MEDIUM   Potential GCP Elastic IP Takeover: demo-gcp.subpipe.run [GCP - 34.152.86.1]  HIGH     Microsoft Azure Takeover Detection: demo-azure.subpipe.run (demo-azure.subpipe.run)  HIGH     ElasticBeanstalk Subdomain Takeover Detection: demo-eb.subpipe.run (demo-eb.subpipe.run)
----------------------------------------------------------------------
✅ Scan Finished in 22.94s
```

## ⚙️ CI/CD Integration & Pipeline Chaining
SubPipe is designed to play perfectly with other command-line tools in your existing recon pipelines. You can pass the API key directly via a flag if needed for automation scripts:
```bash
# Chain directly from subfinder into subpipe
subfinder -d target.com -silent | subpipe --SUBPIPE_API_KEY="your-api-key-here"
```

## ⚠️ Disclaimer
SubPipe is designed exclusively for authorized security research, bug bounty hunting on platforms with explicit scope, and defensive posture assessment. Users are strictly responsible for ensuring they have permission to scan target infrastructure. The developers assume no liability and are not responsible for any misuse or damage caused by this program.
