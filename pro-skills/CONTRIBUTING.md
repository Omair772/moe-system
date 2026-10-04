# Contributing to pro-skills-senior-full-stack-software-engineer

> **Version:** 1.0.0  
> **Author:** Eng. Salah Allsayani  
> **Email:** eng.salahalssayani@gmail.com  
> **Project:** pro-skills-senior-full-stack-software-engineer  
> **Date:** 2026-09-08  

We welcome contributions from the community! This document provides guidelines for contributing to the project. Please read carefully before submitting any changes.

---

## 📜 Code of Conduct

All contributors are expected to uphold the following standards:

- **Be respectful and inclusive** in all interactions
- **Be constructive and supportive** in feedback and discussions
- **Respect diverse perspectives and experiences**
- **Use professional and welcoming language**
- **Accept criticism gracefully** and focus on improvement
- **Focus on what is best not just for ourselves, but for the overall community**
- **Show empathy towards other community members**

Harassment, discrimination, or abusive behavior of any kind will not be tolerated.

---

## 🚀 Getting Started

### Prerequisites
- Git installed on your system
- A GitHub account
- Basic familiarity with Markdown and YAML
- Access to the repository at [https://github.com/salahAlssayani/pro-skills-senior-full-stack-software-engineer](https://github.com/salahAlssayani/pro-skills-senior-full-stack-software-engineer)

### Setting Up Your Development Environment
1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/pro-skills-senior-full-stack-software-engineer.git
   cd pro-skills-senior-full-stack-software-engineer
   ```
3. **Add the upstream remote:**
   ```bash
   git remote add upstream https://github.com/salahAlssayani/pro-skills-senior-full-stack-software-engineer.git
   ```
4. **Create a branch** for your contribution:
   ```bash
   git checkout -b feature/your-feature-name
   ```

---

## 📝 Types of Contributions

We accept the following types of contributions:

| Type | Description |
|------|-------------|
| **Bug Fixes** | Fix errors, typos, or incorrect information in skill files |
| **New Skills** | Add new expert skill domains to the package |
| **Skill Improvements** | Enhance existing skills with better frameworks, examples, or patterns |
| **Documentation** | Improve README, CONTRIBUTING, or any other documentation |
| **Cross-Skill Patterns** | Add new orchestration patterns or workflows |
| **Testing** | Add test cases, validation checks, or quality gates |
| **Tooling** | Improve scripts, automation, or build processes |
| **Translation** | Translate skill content into other languages |
| **Examples** | Add practical examples, templates, or case studies |

---

## 📄 Skill File Guidelines

All skill files must follow strict formatting and content requirements to maintain consistency across the package.

### YAML Frontmatter

Every skill file must begin with a YAML frontmatter block containing the following fields:

```yaml
---
name: skill-name
version: 1.0.0
author: Eng. Salah Allsayani
email: eng.salahalssayani@gmail.com
skill_number: 01
domain: [Category Name]
tags: [tag1, tag2, tag3]
status: Complete
description: Brief description of the skill domain.
created: 2026-09-08
updated: 2026-09-08
---
```

### Required Sections

Each skill file must contain the following sections:

1. **Core Principles** — Fundamental rules and philosophies of the domain
2. **Frameworks & Methodologies** — Structured approaches and methodologies
3. **Checklists & Quality Gates** — Verification checklists and acceptance criteria
4. **Templates & Patterns** — Reusable templates, code patterns, and examples
5. **Cross-Skill Integration** — How this skill integrates with other skill domains
6. **Reference Materials** — Links, citations, and additional reading

### Version Requirements

- All skill files must use **Semantic Versioning** (SemVer)
- Version format: `MAJOR.MINOR.PATCH` (e.g., `1.0.0`)
- Increment MAJOR for incompatible changes
- Increment MINOR for backward-compatible additions
- Increment PATCH for backward-compatible bug fixes
- Every version change must be documented in the changelog

### Classification Levels

Each skill is classified into one of the following levels:

| Level | Description |
|-------|-------------|
| **Beginner** | Foundational knowledge for those starting in the domain |
| **Intermediate** | Practical application for experienced practitioners |
| **Advanced** | Expert-level strategies for senior professionals |
| **Master** | Cutting-edge techniques for industry leaders |

All 14 skills in this package are classified at the **Advanced/Master** level, reflecting the senior full-stack engineering target audience.

---

## 🔄 Development Workflow

Contributions follow a structured five-phase workflow:

### Phase 1: RFC (Request for Comments)
1. Open an issue describing the proposed change or new skill
2. Provide a detailed description including:
   - Motivation and problem statement
   - Proposed solution or approach
   - Impact on existing skills
   - Expected outcomes
3. Discuss with the community and maintainers
4. Obtain consensus before proceeding

### Phase 2: Draft
1. Create a dedicated branch from `main`
2. Develop the skill file or change according to guidelines
3. Ensure all required sections are present
4. Validate YAML frontmatter syntax
5. Follow the project naming conventions
6. Write comprehensive documentation

### Phase 3: Review
1. Submit a Pull Request against `main`
2. Assign at least one reviewer
3. Address all feedback and requested changes
4. Ensure the contribution meets all quality standards
5. Verify cross-skill integration compatibility

### Phase 4: Test
1. Run all validation checks
2. Verify the skill file loads correctly with the ENTRY.md orchestration
3. Test any new cross-skill patterns
4. Confirm checklists and quality gates function as expected
5. Validate formatting consistency across all files

### Phase 5: Approve & Deploy
1. Maintainer approval and sign-off
2. Merge the Pull Request to `main`
3. Update version numbers if applicable
4. Tag the release with the new version
5. Update the README.md skills inventory table
6. Publish any related documentation updates

---

## 💬 Commit Conventions

All commits must follow the Conventional Commits specification:

### Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
| Type | Description |
|------|-------------|
| `feat` | A new feature or skill |
| `fix` | A bug fix |
| `docs` | Documentation changes |
| `style` | Formatting, missing semicolons, etc. |
| `refactor` | Code or content restructuring |
| `test` | Adding or updating tests |
| `chore` | Maintenance, tooling, or build changes |
| `revert` | Reverting a previous commit |

### Examples
```bash
git commit -m "feat(skills-15): add new cloud-native architecture skill"
git commit -m "fix(06-security): correct threat model examples"
git commit -m "docs(readme): update skills inventory table"
git commit -m "refactor(entry): improve orchestration logic"
```

### Commit Message Body
When the commit affects multiple files or has significant impact, include a body:
```
feat(skills): add requirements engineering skill domain

- Adds SKILL-14 requirements engineering module
- Includes templates and checklists
- Integrates with SKILL-07 architecture and SKILL-13 HCI/UI-UX
```

---

## 📥 Pull Request Process

1. **Ensure your branch is up to date** with the latest `main` branch
2. **Run all checks** and validations before submitting
3. **Update relevant documentation** (README, skill inventory, version history)
4. **Link the related issue** in your PR description
5. **Provide a clear description** of what the PR changes and why
6. **Request a review** from at least one maintainer
7. **Address review feedback** promptly and professionally
8. **Do not merge your own PR** — wait for maintainer approval
9. **Squash commits** if multiple small commits address one change

### PR Checklist
- [ ] YAML frontmatter is valid and complete
- [ ] All required sections are present
- [ ] Version follows semantic versioning
- [ ] Naming conventions are consistent
- [ ] Cross-skill integration is documented
- [ ] README.md is updated if needed
- [ ] No placeholder text remains
- [ ] All tests pass
- [ ] Code formatting matches project standards

---

## ✅ Quality Standards

All contributions must meet the following quality benchmarks:

### Content Quality
- All content must be **accurate, actionable, and production-ready**
- No placeholder text, TODO comments, or incomplete sections
- All examples must be realistic and based on professional experience
- Content must be written at the **Advanced/Master** level

### Formatting Standards
- Consistent Markdown formatting throughout
- Proper heading hierarchy (H1 → H2 → H3 → H4)
- Code blocks with appropriate language identifiers
- Tables properly formatted with aligned columns
- Consistent use of emojis for visual hierarchy

### Technical Standards
- All skills must be **tested against the ENTRY.md orchestration**
- Cross-skill patterns must be validated
- Checklists must be complete and actionable
- Templates must be ready for immediate use
- No broken links or references

### Documentation Standards
- Every skill file must have complete YAML frontmatter
- Version history must be maintained
- Author attribution must be present in every file
- All external references must be verified and current

---

## 🐛 Reporting Issues

If you find a bug, inconsistency, or have a suggestion, please report it:

1. **Check existing issues** to avoid duplicates
2. **Create a new issue** with the following information:
   - Clear title describing the issue
   - Description of the problem
   - Expected behavior vs. actual behavior
   - Steps to reproduce (if applicable)
   - Relevant skill file and section
   - Proposed solution (if any)

### Issue Labels
| Label | Description |
|-------|-------------|
| `bug` | Something is broken |
| `enhancement` | New feature or improvement |
| `documentation` | Documentation issue |
| `question` | Question or discussion |
| `priority-high` | Needs immediate attention |
| `skill-file` | Issue within a skill file |

---

## 📬 Contact

For questions, suggestions, or collaboration opportunities:

- **Author:** Eng. Salah Allsayani
- **Email:** eng.salahalssayani@gmail.com
- **Career:** Senior Full-Stack Engineer (10+ Years)
- **University:** Taiz University, Alsaeed Faculty of Engineering & IT
- **City:** Taiz, Yemen
- **GitHub:** [salahAlssayani](https://github.com/salahAlssayani)
- **Repository:** [pro-skills-senior-full-stack-software-engineer](https://github.com/salahAlssayani/pro-skills-senior-full-stack-software-engineer)

---

> **Note:** By contributing to this project, you agree to the Code of Conduct and the MIT License terms. All contributions become part of the project and are licensed under the MIT License.

---

**[END OF CONTRIBUTING]**
