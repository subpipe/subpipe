# SubPipe 🚀

> **Detect DNS Vulnerabilities With Absolute Accuracy.**

A high-octane sniper engine designed for bug bounty hunters and red teamers. Stop chasing theoretical noise; get verified takeovers, dangling cloud IPs, and SSRF enablers with zero false positives. 

SubPipe acts as a lightning-fast CLI client that pipes your reconnaissance data directly into our concurrent cloud engines, streaming actionable vulnerabilities back to your terminal in real-time.

## ✨ Core Features

* **Zero False Positives:** Actively verifies live hosts to eliminate noise.
* **Real-Time Streaming:** Instant vulnerability delivery via Server-Sent Events (SSE).
* **Cloud IP Engine:** Detects dangling IPs (AWS/GCP) and internal SSRF enablers (RFC1918).
* **Registrar Engine:** Catches critical NS, MX, and CNAME takeovers using raw API checks.
* **PaaS Takeovers:** Flags unprotected Surge, ElasticBeanstalk, and Azure assets.
* **Smart Deduplication:** Cleans messy URLs and drops duplicates to protect your scan quota.

## 📥 Installation

You can install SubPipe using Go (recommended) or by downloading a pre-compiled binary.

### Option 1: Install via Go (Recommended)
If you have Go installed on your system, this is the fastest way to get started. It will automatically compile the binary for your specific OS and architecture:
```bash
go install https://github.com/subpipe/subpipe@latest
```

## 🔑 Getting an API Key

SubPipe's CLI acts as a lightweight client that routes targets to our concurrent cloud engines. To run scans, you need to grab a free API key from the web dashboard.

1. Go to [subpipe.run](https://subpipe.run)
2. Sign in to access your dashboard
3. Copy your `SUBPIPE_API_KEY`
4. You will automatically receive a quota of free scans (credits are only deducted for valid, actionable targets, not messy input or duplicates).

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
