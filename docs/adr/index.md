# Architecture Decision Records

This section will contain Architecture Decision Records (ADRs) documenting key architectural decisions made during CMDR development.

## What are ADRs?

Architecture Decision Records document important architectural decisions made in a project, including:

- The context and problem being addressed
- The decision that was made
- The rationale behind the decision
- The consequences of the decision

## ADR Format

Each ADR should follow this template:

```markdown
# ADR-XXXX: Title

**Status:** Proposed | Accepted | Deprecated | Superseded

**Date:** YYYY-MM-DD

## Context

What is the issue that we're seeing that is motivating this decision or change?

## Decision

What is the change that we're proposing and/or doing?

## Rationale

Why did we choose this option over alternatives?

## Consequences

What becomes easier or more difficult to do because of this change?

### Positive

- Benefit 1
- Benefit 2

### Negative

- Drawback 1
- Drawback 2

## Alternatives Considered

What other options were considered and why were they rejected?
```

## Existing ADRs

Currently, there are no formal ADRs documented. However, some implicit architectural decisions can be inferred from the codebase:

### Implicit Decision: Factory Pattern for Extensibility

**Context:** Need to support multiple implementations of interfaces (CommandManager, Initializer, etc.)

**Decision:** Use factory pattern with registration

**Rationale:**
- Decouples interface consumers from implementations
- Allows runtime selection of implementations
- Enables easy testing with mock implementations

**Evidence:** [`core/command.go`](https://github.com/mrlyc/cmdr/blob/master/core/command.go), [`core/initializer.go`](https://github.com/mrlyc/cmdr/blob/master/core/initializer.go)

### Implicit Decision: Layered Manager Architecture

**Context:** Need to separate concerns (download, file management, persistence)

**Decision:** Chain managers where each adds functionality

**Rationale:**
- Single Responsibility Principle
- Easy to test each layer independently
- Flexibility to recombine layers

**Evidence:** DownloadManager → BinaryManager → DatabaseManager

### Implicit Decision: Storm/BoltDB for Persistence

**Context:** Need lightweight, embedded database for command metadata

**Decision:** Use Storm (wrapper around BoltDB)

**Rationale:**
- No external database required
- Single file database
- Strong typing with Storm ORM
- Cross-platform compatibility

**Evidence:** [`core/manager/database.go`](https://github.com/mrlyc/cmdr/blob/master/core/manager/database.go)

## Future ADRs

Future architectural decisions should be documented as formal ADRs in this directory, such as:

- ADR-001: Choice of CLI framework (Cobra vs alternatives)
- ADR-002: Configuration management approach (Viper)
- ADR-003: Download strategy pattern
- ADR-004: Shell initialization approach
- ADR-005: Version normalization strategy
