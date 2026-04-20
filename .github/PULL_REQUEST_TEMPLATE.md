<!--
Copyright 2026 The ARCORIS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Before you submit

Use this template to explain the change clearly to reviewers and maintainers.

- Replace the guidance text with actual content.
- If a section does not apply, write `None`.
- Be concrete: include commands, file paths, config keys, metrics, versions, and rollout details where relevant.
- Do not skip validation, compatibility, or rollout information when the change affects runtime behavior.

# Summary

What to write:
- Describe what changed in this PR.
- Keep it short and concrete.
- Preferred size: 3-7 bullets or 1-3 short paragraphs.

Good examples:
- Added queue-pressure input to the admission controller.
- Updated retry lease handling for worker failures.
- Added new shed-rate and backlog metrics.

Avoid:
- Improved the system.
- Fixed several things.

- 
- 
- 

## Linked work

What to write:
- Link the issue, proposal, requirement, or tracking item related to this PR.
- If there is no linked design or requirement, write `None`.
- Use `Closes #123` or `Fixes #123` if this PR should automatically close an issue.

- Issue(s):
- Design / proposal:
- Requirement / tracking reference:
- Auto-close directive:

## Change classification

Select all items that apply to this PR.

- [ ] Bug fix
- [ ] Feature
- [ ] Refactoring with no intended behavior change
- [ ] Performance or efficiency improvement
- [ ] Observability or debugging improvement
- [ ] Documentation change
- [ ] Build / CI / tooling / deployment change
- [ ] Security / isolation / privacy change
- [ ] Breaking change
- [ ] Mechanical change only

## Affected subsystems

Select all ARCORIS subsystems materially affected by this PR.

- [ ] Admission control
- [ ] Scheduling / dispatch
- [ ] Rate control / concurrency regulation
- [ ] Workload isolation / fairness
- [ ] Overload protection / shedding
- [ ] Workers / execution / leases / retries
- [ ] Queue / broker adapters
- [ ] Control plane / policy distribution
- [ ] Coordination / cluster state
- [ ] Observability / metrics / tracing / logging
- [ ] API / CLI / configuration schema
- [ ] Deployment / Kubernetes / Helm / CI
- [ ] Documentation

## Why this change

What to write:
- What problem existed before this change?
- Why is this change needed now?
- What user, operator, or maintainer pain does it remove?

- Problem statement:
- Why now:

## Reviewer guidance

What to write:
- Tell reviewers where they should focus.
- Call out important invariants, contracts, risky paths, or non-obvious trade-offs.
- Say what is intentionally out of scope for this PR.

Examples:
- Review focus: lease expiry handling under worker restarts.
- Key invariants or contracts: no change to default fairness semantics.
- Files worth focused review: `scheduler/policy/*` and `adapter/retry/*`.

- Review focus:
- Key invariants or contracts:
- Files, commits, or flows worth focused review:
- Explicitly out of scope:

## Behavioral and operational impact

What to write:
- Explain what changes at runtime.
- If the PR affects scheduling, admission, fairness, overload handling, coordination, observability, or operator workflows, say so explicitly.
- If a category does not apply, write `None`.

Examples:
- Runtime behavior impact: admission decisions now incorporate queue backlog.
- Operator-visible impact: one new Helm value and two new Prometheus metrics.
- Failure mode sensitivity: bad policy config can over-throttle background work.

- Runtime behavior impact:
- Scheduling, admission, or control-loop impact:
- State, persistence, or coordination impact:
- Operator-visible impact:
- Metrics, logs, or traces added or changed:
- Failure mode, safety, or rollback sensitivity:

## Validation

What to write:
- Give exact validation steps whenever possible.
- Prefer copy-pasteable commands and reproducible scenarios.
- If validation is partial, say what was proven and what was not.

### Validation steps

1.
2.
3.

### Validation evidence

What to write:
- Exact commands or suites that ran.
- Environment used for validation.
- Short evidence summary.
- Remaining gaps, if any.

Example:
- Commands / suites: `go test ./...`, `make integration-test`
- Environment: local Kind cluster with RabbitMQ adapter
- Evidence summary: integration suite passed; new shed-rate metric emitted as expected
- Remaining validation gaps: no multi-node control-plane validation yet

- Commands / suites:
- Environment:
- Evidence summary:
- Remaining validation gaps:

### Tests executed

Select what actually ran for this PR.

- [ ] Unit tests
- [ ] Integration tests
- [ ] End-to-end tests
- [ ] Benchmarks / load tests
- [ ] Lint / static analysis
- [ ] Manual validation
- [ ] Not run, explained above

## Compatibility and rollout

What to write:
- Be explicit even when the answer is no impact.
- If any checkbox below is `Yes`, explain the consequence in the rollout or migration notes.

- Breaking changes: [ ] No  [ ] Yes
- Migration required: [ ] No  [ ] Yes
- API / CLI / schema contract changed: [ ] No  [ ] Yes
- Configuration or policy changes required: [ ] No  [ ] Yes
- Default behavior changed: [ ] No  [ ] Yes
- Feature gate or opt-in path involved: [ ] No  [ ] Yes

### Rollout plan

What to write:
- How should this change be introduced safely?
- Mention feature gates, canary rollout, phased enablement, or monitoring expectations if relevant.

Example:
- Deploy behind a feature gate, enable for one workload class, and watch shed-rate and admission latency for 30 minutes.

- 

### Rollback plan

What to write:
- How should operators disable or revert the change if it causes issues?
- Include config rollback, feature gate disablement, or cleanup steps if needed.

Example:
- Disable the feature gate, roll back Helm values, revert the policy field, and confirm control-plane convergence.

- 

### Upgrade or migration notes

What to write:
- Describe upgrade, migration, schema, config, or operator action required by this PR.
- If nothing is required, write `None`.

Example:
- Operators must add `policy.rateControl.mode`; old configs continue to work unchanged.

- 

## Security, privacy, and provenance

What to write:
- Say whether this PR changes trust boundaries, isolation, credential handling, privacy behavior, or third-party provenance.
- If nothing applies, write `None`.

Examples:
- Security, isolation, or trust-boundary impact: None
- Third-party material copied or adapted: adapted retry backoff approach from <project>, no code copied

- Security, isolation, or trust-boundary impact:
- Privacy or data-handling impact:
- New secrets, credentials, or trust assumptions:
- Third-party material copied or adapted:
- License or attribution follow-up:

## Known limitations and follow-up work

What to write:
- State what is still risky, intentionally deferred, or not yet validated.

- Known limitations:
- Follow-up work:

## Author checklist

Select all items that are true for this PR.

- [ ] I reviewed the diff myself before requesting review.
- [ ] I removed accidental secrets, tokens, keys, personal data, and debug-only artifacts.
- [ ] I added or updated tests where needed, or explained why they were not run.
- [ ] I updated documentation, examples, or comments where needed.
- [ ] I documented behavioral, compatibility, migration, rollout, and rollback impact where relevant.
- [ ] I documented operator-facing and observability impact where relevant.
- [ ] I reviewed licensing and provenance for any third-party material.
- [ ] CI is green, or failing and non-required jobs are explained.
