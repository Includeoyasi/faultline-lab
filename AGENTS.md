# Instructions for this project

- The human user is the team lead and code reviewer. Codex implements the project and prepares reviewable changes.
- Read PROJECT_CONTEXT.md first for the portable project handoff, current status and approvals. Read docs/development/WORKFLOW.md for review expectations and docs/ROADMAP.md for the project scope when the repository is available.
- G0 was explicitly approved by the human on 2026-10-13 (Europe/Moscow): "G0 полностью принят." The architecture and business rules in docs/architecture/DESIGN.md and docs/architecture/PROCESSES.md are approved. Do not ask for G0 approval again. G1-G4 remain future review gates; prepare the work for a gate before requesting its approval.
- Work autonomously within an approved stage. Keep changes small enough for focused code review and include actual validation results and known limitations.
- Stop before progressing through an unapproved review gate. Present a concrete, complete review packet and identify the decision needed from the team lead.
- Significant changes to approved service boundaries, invariants, incompatible contracts or infrastructure scope require a concrete proposal and human approval.
- Keep exactly three business services: Orders, Inventory and Payments. Do not add another business service to host workflow orchestration.
- Never report measurements, tests or approvals that have not actually occurred. Keep proposed decisions, implemented behavior and observed results distinct.
- Keep the portable context current after every substantive task, review resolution, approval, experiment, environment change or newly discovered blocker. Update CONTEXT_STATE.json and the affected source documents, increment the context revision, then run tools/update_context.ps1 and tools/update_context.ps1 -CheckOnly. Include the refreshed PROJECT_CONTEXT.md in the review packet.
- PROJECT_CONTEXT.md is generated from CONTEXT_STATE.json, tools/context_intro.md and the embedded canonical documents. Update those inputs instead of editing only the generated local file. Do not put secrets in any context artifact.
- An agent without the repository works from the dated snapshot and returns an updated context/change proposal with its result. It must distinguish proposed edits from locally applied changes and list checks it could not run. Separate snapshots are reconciled by revision and evidence, not by silently overwriting a newer context.
