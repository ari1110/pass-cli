# Development Documentation

This directory contains implementation tracking documents and development notes. These are **NOT user-facing documentation**.

## Purpose

Development documentation serves to track implementation progress, record design decisions, and provide context for ongoing development work. Unlike user-facing docs (USAGE.md, INSTALLATION.md) or contributor docs (DEVELOPMENT.md), these documents capture the detailed mechanics of building specific features.

## Contents

### Implementation Summaries
- **DASHBOARD_IMPLEMENTATION_SUMMARY.md** - Dashboard feature implementation details and architecture decisions
- Development tracking for the approval dashboard system

### Testing Documentation
- **DASHBOARD_TESTING_CHECKLIST.md** - Comprehensive testing checklist for dashboard functionality
- Validation steps and test coverage documentation

### Audits and Analysis
- **KEYBINDINGS_AUDIT.md** - Keybinding system analysis and documentation
- Records of UI/UX audits and accessibility reviews

## What Belongs Here?

**Include:**
- Implementation tracking documents for specific features
- Testing checklists and validation procedures
- Code audits and technical analysis reports
- Development session notes and decision logs

**Don't Include:**
- User-facing documentation (goes in docs/)
- API documentation (goes in code comments/godoc)
- Contribution guidelines (goes in DEVELOPMENT.md)
- Release notes (goes in CHANGELOG.md or GitHub releases)

## Why Separate from Other Docs?

These documents are developer-oriented and implementation-specific, serving a different purpose than:
- **docs/** - User and maintainer documentation
- **docs/archive/** - Historical artifacts from past milestones
- **test/** - Test code and test utilities

For current user documentation, see [docs/README.md](../README.md).
For archived historical documents, see [docs/archive/README.md](../archive/README.md).
