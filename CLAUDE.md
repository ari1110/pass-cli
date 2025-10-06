# Claude Operational Guide for pass-cli

## Project Context

You are working on **pass-cli**, a secure, cross-platform command-line password and API key manager designed for developers. This is a Go application with dual interfaces: CLI commands and an interactive TUI dashboard.

**Full Project Details**: @.spec-workflow/steering/product.md

---

## Operational Responsibilities

### 1. Spec-Driven Development (MANDATORY)

**ALWAYS use the spec-workflow MCP server tools** for all feature work:

```bash
# Load workflow instructions
mcp__spec-workflow__spec-workflow-guide

# Load steering guide (for project-level docs)
mcp__spec-workflow__steering-guide

# Check spec progress
mcp__spec-workflow__spec-status projectPath:. specName:{spec-name}

# Manage approvals
mcp__spec-workflow__approvals action:request|status|delete
```

**Critical Rules**:
- ⚠️ **NEVER proceed without dashboard approval** - Verbal approval is NOT accepted
- ⚠️ **MUST follow workflow**: Requirements → Design → Tasks → Implementation
- ⚠️ **ALWAYS read steering docs first** when creating a new spec
- ⚠️ **ONE spec at a time** - Complete before starting new work

**Steering Documents** (@.spec-workflow/steering/):
- @.spec-workflow/steering/product.md - Product vision, users, features, goals
- @.spec-workflow/steering/tech.md - Technology stack, architecture, dependencies
- @.spec-workflow/steering/structure.md - Directory organization, file patterns, conventions

These documents define the project's foundation. **Read them first** before creating any spec.

### 2. Technology Standards

**Reference**: @.spec-workflow/steering/tech.md

This document contains:
- Programming language and version requirements
- Key dependencies and libraries
- Application architecture and layers
- Data storage patterns
- External integrations
- Security standards

**CRITICAL**: When you add, remove, or change any technology (dependencies, frameworks, libraries), you **MUST** update @.spec-workflow/steering/tech.md immediately.

### 3. File Organization

**Reference**: @.spec-workflow/steering/structure.md

This document contains:
- Complete directory structure
- File naming conventions
- Module responsibilities
- Testing patterns
- Documentation standards

**CRITICAL**: When you create new directories, move files, or change the project structure, you **MUST** update @.spec-workflow/steering/structure.md immediately.

**When creating new files**:
- Place in appropriate layer directory (per structure.md)
- Follow existing naming patterns (`*_test.go` for tests)
- Use package names matching directory names
- Add unit tests in same directory (`filename_test.go`)

### 4. Code Quality Standards

**Testing Requirements**:
- Unit tests for all new code (`*_test.go`)
- Table-driven tests preferred for Go
- Integration tests in `test/` directory
- Target: 90%+ code coverage

**Before Committing**:
```bash
go fmt ./...           # Format code
go vet ./...           # Static analysis
golangci-lint run      # Comprehensive linting
gosec ./...            # Security scanning
go test ./...          # Run all tests
```

**Commit Messages**:
- Use conventional commits: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
- Include context and "why" not just "what"
- End with co-authorship:
  ```
  🤖 Generated with Claude Code

  Co-Authored-By: Claude <noreply@anthropic.com>
  ```

### 5. Security Standards (CRITICAL)

This is a **security-focused credential management tool**. Security is paramount.

**NEVER**:
- Store passwords in plaintext
- Log sensitive data
- Use weak encryption (always AES-256-GCM)
- Skip secure memory clearing
- Expose credentials in error messages

**ALWAYS**:
- Use `golang.org/x/crypto` for cryptographic functions
- Set file permissions to 600 for vault files
- Validate all inputs before processing
- Use secure random generation (`crypto/rand`)
- Follow OWASP secure coding practices

**Encryption Standards**:
- AES-256-GCM (authenticated encryption)
- PBKDF2 key derivation (100k iterations minimum)
- Cryptographically secure random IV generation

### 6. Workflow for Starting Work

**When you begin a session**:

1. **Check current spec status**:
   ```bash
   mcp__spec-workflow__spec-status projectPath:. specName:{active-spec}
   ```

2. **Read active spec tasks**:
   ```bash
   Read .spec-workflow/specs/{active-spec}/tasks.md
   ```

3. **For current task**:
   - Mark as in-progress: Change `[ ]` to `[-]` in tasks.md
   - Read the `_Prompt` field for guidance
   - Follow `_Leverage` to reuse existing code
   - Implement according to task description
   - Test implementation
   - Mark as complete: Change `[-]` to `[x]` in tasks.md

4. **Continue systematically** through all tasks

### 6.1. Using the Task Tool for Spec Execution

**The Task tool is an agent launcher** that executes complex, multi-step tasks autonomously. When working through a spec, you can delegate tasks to specialized agents by injecting the task's `_Prompt` field.

#### Anatomy of a Spec Task

Each task in `tasks.md` follows a structured format. Every section has a specific purpose:

```markdown
- [ ] 1. Update Model struct view field types to tview
  - File: cmd/tui/model.go
  - Update Model struct to use tview view types instead of Bubble Tea view types
  - Purpose: Enable integration of existing tview view implementations
  - _Leverage: Existing tview view implementations (ListViewTview, DetailViewTview)_
  - _Requirements: 1.1_
  - _Prompt: Role: Go developer completing a framework migration | Task: Update the Model struct in cmd/tui/model.go to use tview view field types... | Restrictions: Only change field TYPE declarations, do not change field names... | Success: Model struct view fields are typed as tview variants..._
```

**Task Anatomy Breakdown**:

1. **Checkbox & Title**: `- [ ] 1. Task title`
   - **Purpose**: Track progress ([ ] pending, [-] in-progress, [x] completed)
   - **Importance**: Provides visual progress tracking and clear task identification

2. **File**: `File: path/to/file.go`
   - **Purpose**: Identifies exactly which files will be modified
   - **Importance**: Scopes the work to specific locations, prevents scope creep

3. **Description**: Brief explanation of what the task does
   - **Purpose**: Summarizes the task in one sentence
   - **Importance**: Quick reference for understanding the task's action

4. **Purpose**: Why this task is needed
   - **Purpose**: Explains the task's role in the larger feature
   - **Importance**: Provides context for why the work matters

5. **_Leverage**: `_Leverage: Existing code/utilities to use_`
   - **Purpose**: Points to existing code, patterns, or utilities to reuse
   - **Importance**: Promotes code reuse, prevents reinventing the wheel, shows you what already exists
   - **Critical**: ALWAYS check this before implementing - don't write what already exists

6. **_Requirements**: `_Requirements: 1.1, 1.2_`
   - **Purpose**: Maps task back to requirements.md acceptance criteria
   - **Importance**: Ensures traceability from requirements → design → tasks → implementation
   - **Critical**: These are the acceptance criteria you MUST satisfy

7. **_Prompt**: `_Prompt: Role: ... | Task: ... | Restrictions: ... | Success: ..._`
   - **Purpose**: Complete implementation guide for autonomous execution
   - **Importance**: THE MOST CRITICAL FIELD - contains everything needed to execute the task
   - **Critical**: This is your execution blueprint - follow it exactly

#### The _Prompt Field Structure

The `_Prompt` field has four parts separated by `|`:

**Format**: `Role: [specialist role] | Task: [detailed instructions] | Restrictions: [constraints] | Success: [completion criteria]`

**1. Role**: Specialist perspective to assume
```
Role: Go developer completing a framework migration
```
- **Purpose**: Sets context and expertise level for the task
- **Importance**: Helps you approach the task with the right mindset and knowledge base

**2. Task**: Detailed step-by-step instructions
```
Task: Update the Model struct in cmd/tui/model.go to use tview view field types by changing listView *views.ListView to listView *views.ListViewTview, changing detailView *views.DetailView to detailView *views.DetailViewTview...
```
- **Purpose**: Provides exact actions to take, with specific code changes
- **Importance**: Removes all ambiguity - tells you EXACTLY what to do

**3. Restrictions**: What NOT to do and constraints to follow
```
Restrictions: Only change field TYPE declarations, do not change field names, do not modify any methods yet, the file must compile after changes
```
- **Purpose**: Defines boundaries and prevents scope creep
- **Importance**: Keeps task focused and prevents breaking changes

**4. Success**: Completion criteria
```
Success: Model struct view fields are typed as tview variants, field names remain unchanged, no other code modified
```
- **Purpose**: Defines "done" - testable completion criteria
- **Importance**: Tells you when to stop and how to verify success

#### How to Use the Task Tool

**Option 1: Execute task directly yourself**
```markdown
1. Read the task's _Prompt field
2. Follow Role → Task → Restrictions → Success
3. Implement according to instructions
4. Verify success criteria met
5. Mark task as complete
```

**Option 2: Delegate to Task agent (for complex tasks)**
```bash
Task tool:
  description: "Execute spec task 5"
  subagent_type: "general-purpose"
  prompt: "Implement the task for spec {spec-name}, task 5:

  Role: Go developer completing a framework migration

  Task: Update the Model struct in cmd/tui/model.go to use tview view field types by changing listView *views.ListView to listView *views.ListViewTview, changing detailView *views.DetailView to detailView *views.DetailViewTview...

  Restrictions: Only change field TYPE declarations, do not change field names, do not modify any methods yet, the file must compile after changes

  Success: Model struct view fields are typed as tview variants, field names remain unchanged, no other code modified

  After completing the task, update .spec-workflow/specs/{spec-name}/tasks.md by changing the task checkbox from [-] to [x]."
```

**When to use the Task tool**:
- ✅ Complex tasks spanning multiple files
- ✅ Tasks requiring extensive searching/analysis
- ✅ Tasks with many steps that can be executed independently
- ✅ When you want an agent to execute the task autonomously

**When NOT to use the Task tool**:
- ❌ Simple 1-2 file edits (do it yourself)
- ❌ When you need tight control over each step
- ❌ Tasks requiring user approval mid-execution

#### Critical Rules for Task Execution

**Whether executing yourself or using Task tool**:

1. ✅ **Read ALL task sections** before starting
   - Understand File, Description, Purpose
   - Check _Leverage for existing code to reuse
   - Note _Requirements to ensure you satisfy acceptance criteria
   - Study _Prompt for exact instructions

2. ✅ **Follow the _Prompt EXACTLY**:
   - Assume the Role
   - Execute the Task step-by-step
   - Respect all Restrictions
   - Verify all Success criteria before marking complete

3. ✅ **Update tasks.md accurately**:
   - `[ ]` → `[-]` when you START the task
   - `[-]` → `[x]` ONLY when all Success criteria are met
   - NEVER mark complete if tests fail or criteria aren't met

4. ✅ **Respect _Leverage**:
   - Read the existing code/utilities mentioned
   - Reuse patterns and structures
   - Don't reinvent what already exists

5. ✅ **Verify _Requirements**:
   - Cross-reference with requirements.md
   - Ensure acceptance criteria are satisfied
   - Don't mark complete if requirements aren't met

**The _Prompt field is your contract** - it contains everything you need to execute the task correctly. Follow it exactly, no interpretation needed.

### 7. Workflow for Creating New Specs

**ALWAYS start by loading the workflow guide**:
```bash
mcp__spec-workflow__spec-workflow-guide
```

**Then follow the sequence**:

1. **Requirements Phase**:
   - Read steering docs (product.md, tech.md, structure.md)
   - Create requirements.md with user stories and EARS criteria
   - Request approval (filePath only)
   - Poll status until approved (NEVER accept verbal)
   - Delete approval before proceeding

2. **Design Phase**:
   - Analyze codebase for patterns to reuse
   - Create design.md with architecture and components
   - Request approval → Poll → Delete

3. **Tasks Phase**:
   - Convert design to atomic tasks
   - Create tasks.md with checkboxes and _Prompt fields
   - Request approval → Poll → Delete

4. **Implementation Phase**:
   - Execute tasks systematically
   - Update checkboxes as you progress

### 8. Handling Errors and Blockers

**When compilation fails**:
1. Read the error message carefully
2. Identify which layer is affected
3. Check if it's a type mismatch (common during migration)
4. Fix at the source, not with workarounds

**When tests fail**:
1. Run individually: `go test -v ./path/to/package -run TestName`
2. Check if test needs updating for current framework/patterns
3. Fix implementation OR update test (whichever is wrong)

**When stuck on a task**:
1. Re-read the task's `_Prompt` field
2. Check `_Leverage` for existing code to reference
3. Read the `_Requirements` to understand acceptance criteria
4. Search codebase for similar patterns (use Grep, Glob)

**When discovering incomplete work**:
1. **STOP immediately** - Don't continue building on broken foundation
2. Document the gap (what was claimed vs. what exists)
3. Create remediation plan
4. Get user approval before proceeding

### 9. Communication Standards

**Be concise and direct**:
- Avoid preamble like "Great!", "Sure!", "Let me help"
- State facts and actions clearly
- Only explain when complexity requires it

**When reporting progress**:
- Use file paths with line numbers: `cmd/tui/model.go:54`
- Show before/after for changes
- Confirm completion, don't elaborate unless asked

**When asking for approval**:
- **NEVER proceed on verbal "approved"**
- Always use approvals tool and check dashboard status
- MUST delete approval successfully before proceeding

### 10. Committing Work During Specs

**Commit frequently and often** when working through spec tasks:

**When to commit**:
- ✅ After completing each task (mark as `[x]`)
- ✅ After completing each phase of a spec
- ✅ After any significant milestone or working state
- ✅ Before switching to a different task
- ✅ When you update steering docs (tech.md, structure.md, product.md)

**Commit message format**:
```
<type>: <description>

<body explaining changes>

<phase reference if applicable>

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Examples**:
```
feat: Integrate tview view implementations into Model struct

- Update Model view field types to tview variants
- Update NewModel() to use tview view constructors
- Fix view method calls for tview compatibility

Phase 1 of tview-migration-remediation spec.

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Why commit frequently**:
- Enables easy rollback to working states
- Provides clear audit trail of implementation
- Allows atomic changes that can be reviewed independently
- Demonstrates systematic progress through spec tasks

### 11. Accuracy and Transparency (CRITICAL)

**Accurate assessments and transparency are the #1 priority in this repository.**

**NEVER**:
- ❌ Claim a task is complete when it's only partially done
- ❌ Mark a task as `[x]` if tests are failing
- ❌ Skip steps in a task to save time
- ❌ Take shortcuts that deviate from the spec
- ❌ Implement differently than the spec describes
- ❌ Ignore acceptance criteria in requirements
- ❌ Hide errors or issues you encounter

**ALWAYS**:
- ✅ Report the actual state of work, not aspirational state
- ✅ If you discover incomplete work, STOP and document the gap
- ✅ If you cannot complete a task, explain why clearly
- ✅ If a spec has errors, surface them immediately
- ✅ Follow the spec exactly as written - no interpretation
- ✅ Execute all steps in a task, even if they seem redundant
- ✅ Test thoroughly before marking tasks complete

**If a spec exists, you MUST follow it with NO QUESTIONS ASKED, ONLY EXECUTION:**

The spec represents deliberate planning and design. If time was taken to create:
- **requirements.md** - Follow every acceptance criterion exactly
- **design.md** - Implement the architecture as specified
- **tasks.md** - Execute every step in every task, in order

**No shortcuts. No deviations. No assumptions.**

If you think the spec is wrong, unclear, or could be improved:
1. **STOP implementation**
2. Document the specific issue
3. Ask the user for clarification or correction
4. Wait for spec update and approval
5. THEN continue implementation

**Do not reinterpret, optimize, or "improve" the spec on your own.** Execute it exactly as written.

---

## Summary: Your Responsibilities

✅ **ALWAYS**:
- Use spec-workflow MCP tools for all feature work
- Read steering docs (@.spec-workflow/steering/) before creating specs
- Follow approval workflow exactly (request → poll → delete)
- Follow specs exactly as written - NO shortcuts, NO deviations
- Report accurate state of work - transparency is #1 priority
- Commit frequently (after each task, phase, milestone)
- Update tech.md when changing dependencies/frameworks
- Update structure.md when changing file organization
- Respect architectural layers (never mix)
- Write tests for all new code
- Follow security standards strictly
- Be concise and direct in communication
- Update task checkboxes as you progress

❌ **NEVER**:
- Proceed without dashboard approval (verbal NOT accepted)
- Take shortcuts or skip steps in spec tasks
- Mark tasks complete if tests are failing
- Reinterpret or "improve" specs on your own
- Mix architectural layers
- Store sensitive data insecurely
- Skip testing or security checks
- Work on multiple specs simultaneously
- Claim work is done when it's only partial

**Critical Rule**: If a spec exists, follow it exactly. No questions asked, only execution.

**When in doubt**: Read steering docs, check existing patterns, follow the spec-workflow exactly.
