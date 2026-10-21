## Problem and result

<!-- Concrete trigger, previous behavior and resulting behavior. -->

## Scope and decisions

<!-- Roadmap packet / review gate; decisions that affect correctness or operation. -->

<!-- Если схема помогает понять изменение, добавьте раздел «Схема взаимодействия» с краткой Mermaid-схемой прямо в описание. Отделяйте реализованное поведение от предложений и будущих частей. Если схема не нужна, пропустите раздел без отметки N/A и объяснения. -->

## Validation

<!-- Exact commands, observed results and material limitations. Do not claim unrun checks. -->

## Review checklist

<!-- Follow docs/development/CODE_REVIEW.md. Mark irrelevant items N/A with a reason; an empty checkbox is not a passing check. -->

- [ ] Scope and resulting behavior match the approved design.
- [ ] Applicable Go format, build, lint, unit, race and vulnerability checks pass; commands/results or N/A reasons are listed above.
- [ ] Errors, cancellation, goroutine/resource lifetime and shutdown are accounted for.
- [ ] Affected transactions, operation keys, duplicate/late events and crash windows preserve business invariants.
- [ ] Applicable regression/integration checks demonstrate the changed behavior; compatibility and recovery are explained.
- [ ] Logs, metrics, input limits and dependency changes meet the review standard.

## Context and review

- [ ] Canonical documentation and portable context reflect this change.
- [ ] Approvals are recorded only when explicitly given by the human reviewer.
- [ ] AI assistance and remaining decisions are described below.

<!-- AI assistance used, human review needed, and any unresolved questions. -->

<!-- The AI prepares the change; the human reviews this HEAD and authorizes merging. CI success does not replace that decision. -->
