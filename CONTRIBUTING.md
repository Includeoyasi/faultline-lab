# Contributing to Faultline Lab

This repository documents a project built with AI implementation and human technical leadership. Changes should be understandable to a reviewer and supported by reproducible evidence.

## Before changing anything

Read the [decisions and evidence snapshot](PROJECT_CONTEXT.md), [approved design](docs/architecture/DESIGN.md) and [development workflow](docs/development/WORKFLOW.md). The executor, AI-Your-Bro, follows [AGENTS.md](AGENTS.md).

The [code review standard](docs/development/CODE_REVIEW.md) covers idiomatic Go, errors, context, concurrency, data integrity, tests and evidence. It also defines blocking comments, optional suggestions and the baseline Go checks to activate with the first module. The detailed standard is currently written in Russian.

G0 is approved. G1–G4 remain future review gates. Significant changes to approved service boundaries, invariants or incompatible contracts require a concrete proposal and the project owner's decision.

## A reviewable change

1. Name the problem, intended behavior and roadmap packet.
2. Keep the diff focused; explain decisions that affect correctness or operability.
3. Run checks appropriate to the change. Record exact commands, observed results and checks that could not run.
4. Update affected specifications and dated decisions/evidence in `CONTEXT_STATE.json`; regenerate `PROJECT_CONTEXT.md`. Live task progress is recorded only in GitHub Issues and PRs.
5. Use the pull request template to present the result. Address review comments before moving through the next gate.

AI assistance should be disclosed in the review packet. The human reviewer remains responsible for accepting decisions. Never claim tests, approvals or measurements that did not occur.

## Documentation checks

Requires Git with the referenced source commit available locally, Node.js 20+ and PowerShell 7+ for the commands below. A shallow checkout may need an explicit fetch of that SHA or full history. The generator also supports Windows PowerShell 5.1 (`powershell` instead of `pwsh`).

```powershell
node --test tools/doc_safety.test.mjs
node tools/check_docs.mjs
pwsh -NoProfile -File tools/update_context.ps1
pwsh -NoProfile -File tools/update_context.ps1 -CheckOnly
```

The checker validates local links, Markdown fence balance, text encoding and accidental machine-specific paths. The context check compares the generated snapshot with all registered sources. These checks do not validate service behavior or fully parse Mermaid syntax.

## Branches and commits

Use `main` as the baseline, short descriptive branches such as `feat/orders-bootstrap`, and focused commits with clear intent. After bootstrap, submit changes through a pull request; the AI does not approve or merge its own work. The human reviews the current change and authorizes merging; passing CI alone is insufficient. Prefer squash merge after acceptance. Suggested prefixes: `feat:`, `fix:`, `docs:`, `test:`, `chore:`. Include both canonical document changes and their generated context in the same review packet.

See the [runtime guide](docs/development/CONTAINER_RUNTIME.md) and [Go smoke](docs/development/SMOKE.md) for executable checks.
