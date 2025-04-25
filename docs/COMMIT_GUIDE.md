# Commit Message Guidelines

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification for creating clear and standardized commit messages.

## Commit Format

Each commit message consists of:
```
<type>(<optional scope>): <description>

<optional body>

<optional footer>
```

Examples:
```
feat: add input validation template
```
```
fix(auth): resolve JWT token refresh issue
```
```
docs: update installation instructions

The installation section was missing the step for environment setup.
```

## Commit Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code changes that neither fix bugs nor add features
- **perf**: Performance improvements
- **test**: Test additions or corrections
- **chore**: Changes to the build process or auxiliary tools
- **ci**: Changes to CI configuration files
- **revert**: Reverts a previous commit

## Using the Commit Helper

We provide an interactive commit script that guides you through creating properly formatted commits:

```bash
# Stage your changes first
git add .

# Then run the commit helper
npm run commit
```

The tool will prompt you for:
1. The type of change
2. Optional scope
3. Short description
4. Longer description if needed
5. Information about breaking changes
6. Issues affected by this commit

## Why We Use Conventional Commits

1. **Automated versioning**: Enables semantic versioning based on commit types
2. **Automated changelogs**: Generates detailed release notes
3. **Clarity**: Makes the project history more readable and structured
4. **Searchability**: Easier to find specific types of changes

## Breaking Changes

For breaking changes, add `BREAKING CHANGE:` in the commit message body:

```
feat: change API authentication flow

BREAKING CHANGE: This changes the authentication flow and requires clients to update their integration.
```

Breaking changes will automatically trigger a major version bump. 