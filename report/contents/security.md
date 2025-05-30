# Security Assessment (May 24, 2025)

**Note:** This assessment was conducted before Docker Swarm was implemented and is therefore outdated. An updated security review is a high priority for future development.

## Risk Identification

Our public-facing asset is a single web-app, while the database instance is protected by a firewall.
To identify the attack surface, we performed a TCP SYN scan of the most common ports against web's IP.
The scan revealed open ports:

- SSH (port 22)
- HTTP (port 80)
- Grafana (port 3000)
- Prometheus (port 9090)

Ports 3000 and 9090 are default monitoring services, and should be privated. Exposing SSH on port 22 is expected for maintenance access, and HTTP on port 80 is the web-app's interface and should remain open.

To uncover vulnerabilities, we used Nmap’s vulnerability scripts against ports 22 and 80, which identified exposure to cross-site request forgery (CSRF) and Slowloris denial-of-service attacks. Given prior incidents of idle-connection exhaustion, the Slowloris finding was expected. A subsequent Nikto scan revealed missing security headers to prevent clickjacking and content sniffing.

## Risk Scenarios

- A successful CSRF attack could trick authenticated users into unknowingly executing malicious actions.
- A Slowloris attack could exhaust server memory and trigger an OOM kill; because Docker’s restart policy does not recover from OOM kills[@docker_docs_resource_constraints], the service would require manual intervention.
- Clickjacking could be achieved by embedding malicious code in transparent frames, deceiving users into performing unintended actions.
- Content sniffing attacks could exploit the browser’s tendency to reclassify responses based on payload content, potentially executing embedded malicious scrips within user-submitted posts.

## Risk Analysis

|                        | **Impact: Low** | **Impact: Medium** | **Impact: High**      |
|------------------------|:-------------:|:--------------:|:-----------------:|
| **Likelihood: Low**    | Clickjacking  |                |                   |
| **Likelihood: Medium** |               |                | Content Sniffing  |
| **Likelihood: High**   | CSRF          |                | Slowloris         |

Table: Overview of likelihood and impact-level of identified scenarios.

Based on this analysis, we prioritized patches in the following order: Slowloris protection, CSRF mitigation, content-sniffing prevention, and clickjacking hardening.

## Mitigation and Remediation

- Slowloris attacks: Configure Read, Write, and Idle connection timeouts on the web-server and impose limits on header size (see PR [#160](https://github.com/DuwuOps/minitwit/pull/160)).
- Slowloris attacks: Enforce maximun database connections, with reduced lifetimes, to prevent resource exhaustion (see PR [#160](https://github.com/DuwuOps/minitwit/pull/160)).
- CSRF: Integrate middleware that issues and validates per-request tokens for all form submissions (see PR [#152](https://github.com/DuwuOps/minitwit/pull/158)).
- Content sniffing: Add response headers instructing browsers not to infer MIME types (see PR [#157](https://github.com/DuwuOps/minitwit/pull/167)).
- Clickjacking: Set the `X-Frame-Options: DENY` header on all responses (see PR [#157](https://github.com/DuwuOps/minitwit/pull/167)).
