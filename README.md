# kipf

Token-slim CLI proxy for AI coding agents. Swallows verbose command output, spits out only what matters.

Named after Vienna's Kipferl - the crumbly little croissant - with a nod to Kirby, the video game character who inhales everything and keeps only what's useful. That's what kipf does to your command output.

## Why

AI coding agents like Claude Code burn tokens on noisy command output. A simple `git push` dumps 15 lines of progress nobody reads. `docker ps` returns a wide table when you need names and status. Every wasted token costs money and shrinks your context window.

kipf sits between the agent and the shell, compressing output before it hits the LLM context.

## Example

```
# without kipf: ~200 tokens
$ git push
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Delta compression using up to 8 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 312 bytes | 312.00 KiB/s, done.
Total 3 (delta 2), reused 0 (delta 0), pack-reused 0
remote: Resolving deltas: 100% (2/2), completed with 2 local objects.
To github.com:user/repo.git
   abc1234..def5678  main -> main

# with kipf: ~10 tokens
$ kipf git push
ok ✓ main
```

## Install


## Usage

```bash
kipf <command>
```

### Git

| Command | What kipf does |
|---|---|
| `kipf git status` | Short format with branch |
| `kipf git log` | Oneline, last 10 commits |
| `kipf git diff` | Stat summary only |
| `kipf git push` | Just the result + branch |
| `kipf git pull` | Just the result + file count |
| `kipf git fetch` | Quiet, confirms ok |
| `kipf git add .` | Confirms ok |
| `kipf git commit -m "msg"` | Confirms ok + short hash |

## Claude Code Integration

kipf works standalone, but shines when paired with a Claude Code PreToolUse hook that transparently rewrites commands:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "~/.claude/hooks/kipf-rewrite.sh"
          }
        ]
      }
    ]
  }
}
```

With the hook installed, when Claude Code runs `git status`, it transparently becomes `kipf git status`. Zero config, zero token waste.

## Design Principles

- **One binary, zero dependencies** — single Go binary, no runtime needed
- **Transparent** — same exit codes, stderr on failure, machine-parseable output
- **Conservative** — only rewrites what it knows. Unknown commands pass through unchanged
- **Composable** — plays nice with pipes, redirects, and scripts

## Roadmap

- [ ] Git commands
- [ ] GitHub CLI
- [ ] Bash commands
- [ ] Docker commands

## License

MIT
