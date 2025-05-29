# Security Assessment (May 24, 2025)

**Note:** This assessment was conducted before Docker Swarm was implemented and is therefore outdated. An updated security review is a high priority for future development.

## Risk Identification

Our public-facing asset is a single web application droplet, while the database instance is protected by a firewall. To enumerate the attack surface, we performed a TCP SYN scan of the most common ports against our staging droplet (configured identically to production) at `167.71.64.240`. This scan revealed SSH on port 22, HTTP on port 80, Grafana on port 3000, and the Zeus admin interface on port 9090. Although ports 3000 and 9090 correspond to default monitoring services, we have already changed Grafana’s default credentials; closing these ports remains an open action item. Exposing SSH on port 22 is expected for maintenance access, and HTTP on port 80 is required for normal operation.

A version scan against our secondary droplet at `134.209.137.191` confirmed that OpenSSH 9.7p1 and a Go net/http server are running on ports 22 and 80, respectively.

To uncover vulnerabilities, we used Nmap’s vulnerability scripts against ports 22 and 80 on `167.71.64.240`, which identified exposure to cross-site request forgery (CSRF) and Slowloris denial-of-service attacks. Given prior incidents of idle-connection exhaustion, the Slowloris finding was expected. A subsequent Nikto scan of `http://134.209.137.191` revealed missing security headers to prevent clickjacking and content sniffing.

## Risk Scenarios

A successful CSRF attack could trick authenticated users into unknowingly executing state-changing actions—while our application has no privileged admin roles, the impact could include unauthorized posts or exploitation of influencer accounts to promote fraud. Slowloris attacks could originate from a single compromised host running a half-open connection script, exhausting server memory and triggering an OOM kill; because Docker’s restart policy does not recover from OOM kills, the service would require manual intervention. Clickjacking could be achieved by embedding our interface within a transparent iframe on a malicious site, deceiving users into performing unintended actions. Content sniffing attacks could exploit the browser’s tendency to reclassify responses based on payload content, potentially executing embedded JavaScript within user-submitted posts if our headers do not explicitly forbid sniffing.

## Risk Analysis

|                        | Impact: Low   | Impact: Medium | Impact: High      |
|------------------------|:-------------:|:--------------:|:-----------------:|
| **Likelihood: Low**    | Clickjacking  |                |                   |
| **Likelihood: Medium** |               |                | Content Sniffing  |
| **Likelihood: High**   | CSRF          |                | Slowloris         |

Based on this analysis, we have prioritized remediations in the following order: Slowloris protection, CSRF mitigation, content-sniffing prevention, and clickjacking hardening.

## Mitigation and Remediation

All identified vulnerabilities have been addressed. To guard against Slowloris attacks, we configured Read, Write, and Idle connection timeouts on the web server and imposed limits on header size. In addition, database connection pooling now enforces maximum open and idle connections with reduced lifetimes to prevent resource exhaustion. CSRF protection was implemented by integrating middleware that issues and validates per-request tokens for all form submissions. To prevent content sniffing, we added response headers instructing browsers not to infer MIME types. Finally, clickjacking is blocked by setting the `X-Frame-Options: DENY` header on all responses.
