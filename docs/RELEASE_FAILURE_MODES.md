# Release Failure Modes

Companion to `RELEASE_CHECKLIST.md`. Every entry here is a failure this repo has
actually hit, with the evidence that proved it and the command that recovers it.
Read this before a tag push and after any red Release run.

Last verified: 2026-10-20 (v0.38.7 release).

## The one-line summary

**A green `Release` run is not the same as a published release.** For v0.36.0,
v0.37.0, v0.38.0, v0.38.5, v0.38.6, and v0.38.7 every build/test/lint/benchmark
step passed, and multiple of them still produced no GitHub Release and no website
redeploy. Always assert the artifacts, not the run colour.

## Failure modes

| # | Mode                                                    | Symptom                                                                                      | Evidence (2026-10-19/20)                              |
| - | ------------------------------------------------------- | -------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| 1 | **>3 tags in one push**                                 | No `Release` run is created at all; nothing fails either — the workflow simply never starts. | v0.38.6: pushing `refs/tags/*v0.38.6` produced zero runs. GitHub: "Events will not be created for tags when more than three tags are pushed at once". |
| 2 | **No git identity on the runner**                       | `Create and push submodule tags` fails with `fatal: empty ident name (…not allowed)`, exit 128. | v0.38.7 run `39328156656`. Same step failed for every earlier tag (15 s–4 min runs). Fixed in `#59cd988`. |
| 3 | **Duplicate annotated tag object**                      | `git push` of a submodule tag is rejected (exit 128) when origin already has that tag with a different object hash. | v0.36.0/v0.37.0 era. The step now skips tags that exist on origin. |
| 4 | **Frozen workflow file at the tag**                     | A workflow fix cannot repair an existing tag: the tagged commit's workflow is what runs. Re-running a failed tag job re-runs the same broken file. | v0.38.7 rerun failed identically.                      |
| 5 | **CI-only gate red at the tag**                         | A gate that exists only in CI (e.g. the 80% coverage threshold) can fail a tag's check permanently. | v0.38.6 tagged at `daghtml` 77.2%; the fix landed after (`#c3a9094`), so only `v0.38.7` shows green. |
| 6 | **Website blocked independently of the Go release**     | `Website`/`Build Website` fails while the Go release is fine (e.g. `astro check` rejects `typescript@7`, or `pnpm install --frozen-lockfile` fails because the manifest was bumped without re-locking). | v0.38.6/7: TS 7 rejection; BuildFlow `pnpm-update` manifest bump. |
| 7 | **Transient network failure mid-run**                   | A module zip download fails (`stream error: … INTERNAL_ERROR`), failing the build step. | v0.38.7 first attempt, `charmbracelet/x/exp/teatest/v2`. |

## Publish order (avoids mode 1 by construction)

```bash
git push origin master        # ordinary commits first
git push origin vX.Y.Z        # ROOT TAG ALONE — this is the event that fires release.yml
```

`release.yml` then creates and pushes the 16 submodule tags (mode 2 is fixed) and
creates the GitHub Release. Push the submodule tags by hand only if that step did
not run:

```bash
git push origin 'refs/tags/*vX.Y.Z'   # many tags at once is fine here — no workflow needs those events
```

## Recovery procedures

| Missing artifact        | Command                                                                                        |
| ----------------------- | ---------------------------------------------------------------------------------------------- |
| GitHub Release          | `gh release create vX.Y.Z --generate-notes --latest`                                            |
| Submodule tags          | `git push origin 'refs/tags/*vX.Y.Z'` from a clone that has the family (`scripts/tag-release.sh` creates it) |
| Website redeploy        | `gh workflow run website.yml --ref vX.Y.Z` (use `--ref master` if the tag's tree is itself broken) |
| Transient failure       | `gh run rerun <run-id> --failed`                                                                |
| A red check on a tag    | Nothing — cut the next patch release (mode 4/5). Do **not** plan to re-tag.                     |

## Post-release assertions (run these, do not assume)

```bash
VERSION=v0.38.7
git ls-remote --tags origin | grep -c "$VERSION"   # expect 34 (17 tags + 17 peeled)
gh release view "$VERSION" --json tagName --jq .tagName
gh run list --workflow=release.yml --limit 3        # inspect the CONCLUSION, not the step list
gh workflow run website.yml --ref "$VERSION"        # then confirm Build + Deploy both succeed
go list -m -versions github.com/larsartmann/go-output/daghtml | tr ' ' '\n' | grep -qx "$VERSION"  # proxy serves it
```

Any missing item means the release is not published, regardless of how the
`Release` run looks.

## History

- v0.36.0, v0.37.0, v0.38.0, v0.38.5, v0.38.6, v0.38.7 — Releases created
  manually or not at all because the tag step failed (modes 2/3). This is why the
  checklist's step 8 exists and why the workflow now sets a git identity and runs
  the Release before the tags.
- v0.38.6 — the newest tag showed a red `Build & Test` check (mode 5) until
  v0.38.7 superseded it.
