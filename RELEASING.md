# Releasing

This project uses semantic version tags for releases.

## Release Checklist

1. Ensure the working tree is clean.
2. Run tests locally:

   ```bash
   go test ./...
   ```

3. Update [CHANGELOG.md](./CHANGELOG.md):
   - Move relevant entries from `Unreleased` to the new version section.
   - Set the release date.
4. Commit changelog or release note adjustments if needed.
5. Create an annotated tag:

   ```bash
   git tag -a v0.1.0 -m "v0.1.0"
   ```

6. Push branch and tag:

   ```bash
   git push origin main
   git push origin v0.1.0
   ```

7. Create the GitHub release from the tag and use the matching changelog section as release notes.

## Versioning Notes

- `v0.x.y` is used while the API is still evolving.
- Backward-incompatible changes may still happen before `v1.0.0`.
