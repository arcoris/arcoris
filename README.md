<div align="center">
  <h1>ARCORIS</h1>
  <p><strong>Open Source Platform for Distributed Scheduling and Execution Control</strong></p>
  <p>Admission control · Adaptive resource management · Workload isolation · Feedback-driven overload protection</p>
  <p>
    <a href="https://github.com/ARCORIS/arcoris">
      <img alt="Repository" src="https://img.shields.io/badge/Repository-GitHub-111827?style=for-the-badge&logo=github" />
    </a>
    <a href="./LICENSE">
      <img alt="License" src="https://img.shields.io/badge/License-Apache%202.0-D22128?style=for-the-badge&logo=apache" />
    </a>
    <img alt="Domain" src="https://img.shields.io/badge/Domain-Distributed%20Scheduling-0F766E?style=for-the-badge" />
    <img alt="Scope" src="https://img.shields.io/badge/Scope-Queues%20%7C%20Workers%20%7C%20Backends-1D4ED8?style=for-the-badge" />
  </p>
  <p>
    <a href="https://github.com/ARCORIS/arcoris/actions">
      <img alt="Build" src="https://img.shields.io/badge/Build-passing-2EA043?style=for-the-badge" />
    </a>
    <a href="https://github.com/ARCORIS/arcoris/releases">
      <img alt="Release" src="https://img.shields.io/badge/Release-v0.1.0--dev-2563EB?style=for-the-badge" />
    </a>
    <a href="https://pkg.go.dev/github.com/ARCORIS/arcoris">
      <img alt="Go Reference" src="https://img.shields.io/badge/Go%20Reference-pkg.go.dev-00ADD8?style=for-the-badge&logo=go" />
    </a>
    <img alt="Coverage" src="https://img.shields.io/badge/Coverage-pending-6B7280?style=for-the-badge" />
  </p>
  <p>
    <img alt="OpenSSF Scorecard" src="https://img.shields.io/badge/OpenSSF%20Scorecard-pending-6B7280?style=for-the-badge" />
    <img alt="OpenSSF Best Practices" src="https://img.shields.io/badge/OpenSSF%20Best%20Practices-in%20progress-9A6700?style=for-the-badge" />
    <img alt="REUSE" src="https://img.shields.io/badge/REUSE-pending-6B7280?style=for-the-badge" />
    <img alt="Go Report Card" src="https://img.shields.io/badge/Go%20Report%20Card-pending-6B7280?style=for-the-badge" />
  </p>
</div>

---

ARCORIS is an open source platform for building distributed schedulers and execution control systems across queues, workers, and heterogeneous execution backends. It provides core mechanisms for admission control, adaptive resource management, workload isolation, and feedback-driven overload protection, helping systems remain stable, fair, and efficient under changing load.

ARCORIS is designed for environments where scheduling is more than simple placement. It helps determine when work should be admitted, how much concurrency the system can safely sustain, how capacity should be shared across competing workloads, and how execution should be shaped when demand exceeds safe operating limits. In practice, this means giving platform and infrastructure teams the building blocks to regulate throughput, protect critical workloads, isolate noisy neighbors, and avoid cascading overload.

ARCORIS builds on established ideas from distributed scheduling, queue-based execution, resource governance, and runtime control. It is intended for modern systems that operate across multiple queues, worker pools, and backend types, where resilience, fairness, and controlled execution matter as much as raw throughput.
