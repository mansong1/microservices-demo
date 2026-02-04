# Working with Forks: Contributing to Online Boutique from a Fork

This guide explains how to contribute to Online Boutique when you've forked the repository, including how to keep your fork synchronized with the upstream repository and create pull requests.

## Table of Contents

- [Initial Fork Setup](#initial-fork-setup)
- [Syncing Your Fork with Upstream](#syncing-your-fork-with-upstream)
- [Creating a Pull Request from Your Fork](#creating-a-pull-request-from-your-fork)
- [Handling Merge Conflicts](#handling-merge-conflicts)
- [Best Practices](#best-practices)

## Initial Fork Setup

### 1. Fork the Repository

1. Visit the [Online Boutique repository](https://github.com/GoogleCloudPlatform/microservices-demo) on GitHub
2. Click the "Fork" button in the top-right corner
3. Select your account or organization as the destination

### 2. Clone Your Fork

Clone your forked repository to your local machine:

```bash
git clone https://github.com/YOUR_USERNAME/microservices-demo.git
cd microservices-demo/
```

Replace `YOUR_USERNAME` with your GitHub username.

### 3. Add the Upstream Remote

Add the original repository as an upstream remote to track changes:

```bash
git remote add upstream https://github.com/GoogleCloudPlatform/microservices-demo.git
```

Verify your remotes:

```bash
git remote -v
```

You should see:
```
origin    https://github.com/YOUR_USERNAME/microservices-demo.git (fetch)
origin    https://github.com/YOUR_USERNAME/microservices-demo.git (push)
upstream  https://github.com/GoogleCloudPlatform/microservices-demo.git (fetch)
upstream  https://github.com/GoogleCloudPlatform/microservices-demo.git (push)
```

## Syncing Your Fork with Upstream

It's important to keep your fork synchronized with the upstream repository to avoid conflicts and ensure you're working with the latest code.

### Quick Sync (Recommended Method)

Sync your fork's default branch with upstream:

```bash
# Fetch the latest changes from upstream
git fetch upstream

# Switch to your main branch
git checkout main

# Merge upstream changes into your main branch
git merge upstream/main

# Push the updates to your fork
git push origin main
```

### Alternative: Sync via GitHub UI

GitHub provides a built-in sync feature:

1. Go to your fork on GitHub
2. Click "Sync fork" button (if available)
3. Click "Update branch" to sync with upstream

### Syncing Feature Branches

If you have an existing feature branch that needs to be updated:

```bash
# Switch to your feature branch
git checkout your-feature-branch

# Fetch and merge upstream changes
git fetch upstream
git merge upstream/main

# Resolve any conflicts if they occur (see section below)

# Push the updated branch to your fork
git push origin your-feature-branch
```

## Creating a Pull Request from Your Fork

### 1. Create a Feature Branch

Always create a new branch for your changes:

```bash
# Make sure you're starting from an updated main branch
git checkout main
git pull upstream main

# Create and switch to a new feature branch
git checkout -b feature/your-feature-name
```

Use descriptive branch names like:
- `feature/add-new-service`
- `fix/cart-service-bug`
- `docs/update-readme`

### 2. Make Your Changes

Make your code changes following the [contribution guidelines](../.github/CONTRIBUTING.md):

```bash
# Make your changes to the code

# Stage your changes
git add .

# Commit with a descriptive message
git commit -m "Add: descriptive commit message"
```

### 3. Push to Your Fork

Push your feature branch to your fork:

```bash
git push origin feature/your-feature-name
```

If this is your first push of this branch, you might need:

```bash
git push -u origin feature/your-feature-name
```

### 4. Create the Pull Request

1. Go to the [upstream repository](https://github.com/GoogleCloudPlatform/microservices-demo)
2. Click "Pull requests" → "New pull request"
3. Click "compare across forks"
4. Set the base repository to `GoogleCloudPlatform/microservices-demo` and base branch to `main`
5. Set the head repository to `YOUR_USERNAME/microservices-demo` and compare branch to your feature branch
6. Click "Create pull request"
7. Fill in the pull request template with:
   - A clear title
   - Description of changes
   - Any related issue numbers
   - Testing performed
8. Submit the pull request

## Handling Merge Conflicts

Merge conflicts occur when your changes conflict with changes made in the upstream repository.

### Resolving Conflicts During Sync

If you encounter conflicts while merging upstream changes:

```bash
# After git merge upstream/main shows conflicts

# 1. View conflicted files
git status

# 2. Open and manually resolve conflicts in each file
#    Look for conflict markers: <<<<<<<, =======, >>>>>>>

# 3. After resolving, stage the resolved files
git add path/to/resolved/file

# 4. Complete the merge
git commit -m "Merge upstream changes and resolve conflicts"

# 5. Push to your fork
git push origin your-branch-name
```

### Avoiding Conflicts

To minimize merge conflicts:

1. **Sync frequently**: Regularly update your fork with upstream changes
2. **Small, focused changes**: Make smaller, targeted pull requests
3. **Communicate**: Check existing issues and PRs to avoid duplicate work
4. **Branch per feature**: Use separate branches for different features

## Best Practices

### Before Starting Work

- [ ] Review [contribution guidelines](../.github/CONTRIBUTING.md)
- [ ] Check [existing issues](https://github.com/GoogleCloudPlatform/microservices-demo/issues)
- [ ] Sync your fork with upstream
- [ ] Create a feature branch from updated main

### While Working

- [ ] Make focused, atomic commits
- [ ] Write clear commit messages
- [ ] Test your changes locally
- [ ] Follow the project's code style
- [ ] Update documentation if needed

### Before Creating PR

- [ ] Sync your branch with latest upstream
- [ ] Run tests and ensure they pass
- [ ] Review your own changes
- [ ] Write a clear PR description
- [ ] Link related issues

### After Creating PR

- [ ] Respond to review feedback promptly
- [ ] Make requested changes in new commits
- [ ] Keep your PR updated with upstream changes
- [ ] Be patient and respectful

## Important Notes for CI/CD

**Note**: The Online Boutique CI/CD pipelines have limitations when working with forks:

> In order for the current CI/CD setup to work on your pull request, you must branch directly off the repo (no forks). This is because the Github secrets necessary for these tests aren't copied over when you fork.

This means:
- Some automated tests may not run on PRs from forks
- Deploy tests require access to the main repository
- Maintainers may need to run tests after merging

If you need to test your changes before submitting a PR:
1. Follow the [development guide](development-guide.md) to test locally
2. Use your own GCP project and Kubernetes cluster
3. Provide test results in your PR description

## Getting Help

If you encounter issues:

1. Check the [development guide](development-guide.md) for local testing instructions
2. Review [existing issues](https://github.com/GoogleCloudPlatform/microservices-demo/issues)
3. Ask questions in your pull request
4. [Create a new issue](https://github.com/GoogleCloudPlatform/microservices-demo/issues/new/choose) for bugs or feature requests

## Additional Resources

- [GitHub Fork Documentation](https://docs.github.com/en/get-started/quickstart/fork-a-repo)
- [GitHub Pull Request Documentation](https://docs.github.com/en/pull-requests)
- [Git Documentation](https://git-scm.com/doc)
- [Online Boutique Contributing Guidelines](../.github/CONTRIBUTING.md)
- [Online Boutique Development Guide](development-guide.md)
