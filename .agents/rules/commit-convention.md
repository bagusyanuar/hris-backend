# Git Commit & Versioning (Conventional Commits)

Semua commit harus menggunakan format **Conventional Commits**:
- `feat:` (fitur baru)
- `fix:` (perbaikan bug)
- `docs:` (dokumentasi, termasuk PRD/RFC)
- `refactor:` (restrukturisasi kode)
- `chore:` (maintenance, update dependensi, generate wire)

**Aturan Penting:**
1. **Atomik:** Dilarang menggabung perubahan fitur A dengan bugfix B dalam satu commit. Pecah menjadi beberapa commit jika perlu.
2. **Workflow:** Gunakan *slash command* `/git-commit` agar AI mengelompokkan file dan menulis pesan commit secara otomatis.
3. **Changelog:** Jika ada perubahan besar, AI akan memperbarui file `CHANGELOG.md`.
