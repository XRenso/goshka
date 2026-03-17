import os
import re
import subprocess
import unittest
from pathlib import Path

TOOL_DIR = Path(__file__).resolve().parent.parent   # GitCliTool/
BINARY   = TOOL_DIR / "gitpeek"


def build_binary():
    result = subprocess.run(
        ["go", "build", "-o", str(BINARY), "."],
        cwd=str(TOOL_DIR),
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        raise RuntimeError(f"go build failed:\n{result.stderr}")


def run_gitpeek(*args, env=None):
    cmd_env = os.environ.copy()
    if env:
        cmd_env.update(env)
    result = subprocess.run(
        [str(BINARY), *args],
        capture_output=True,
        text=True,
        env=cmd_env,
    )
    return result.returncode, result.stdout, result.stderr


def setUpModule():
    if not BINARY.exists():
        build_binary()


class TestParseRepoInput(unittest.TestCase):

    def test_empty_arg_prints_usage(self):
        rc, _, stderr = run_gitpeek()
        self.assertNotEqual(rc, 0)
        self.assertIn("Usage", stderr)

    def test_garbage_input_error(self):
        rc, _, stderr = run_gitpeek("not_a_repo!!!")
        self.assertNotEqual(rc, 0)
        self.assertIn("invalid repository format", stderr)

    def test_single_word_error(self):
        rc, _, stderr = run_gitpeek("justoneword")
        self.assertNotEqual(rc, 0)
        self.assertIn("invalid repository format", stderr)

    def _is_parse_error(self, stderr: str) -> bool:
        return "invalid repository format" in stderr

    def test_owner_slash_repo_format(self):
        _, _, stderr = run_gitpeek("torvalds/linux")
        self.assertFalse(self._is_parse_error(stderr), stderr)

    def test_full_https_url(self):
        _, _, stderr = run_gitpeek("https://github.com/torvalds/linux")
        self.assertFalse(self._is_parse_error(stderr), stderr)

    def test_github_com_prefix(self):
        _, _, stderr = run_gitpeek("github.com/torvalds/linux")
        self.assertFalse(self._is_parse_error(stderr), stderr)

    def test_dot_git_suffix_stripped(self):
        _, _, stderr = run_gitpeek("github.com/torvalds/linux.git")
        self.assertFalse(self._is_parse_error(stderr), stderr)

    def test_trailing_slash_path_ignored(self):
        _, _, stderr = run_gitpeek("github.com/torvalds/linux/tree/master")
        self.assertFalse(self._is_parse_error(stderr), stderr)

    def test_hyphens_and_dots_in_names(self):
        _, _, stderr = run_gitpeek("owner-name/repo.name")
        self.assertFalse(self._is_parse_error(stderr), stderr)

    def test_underscore_in_names(self):
        _, _, stderr = run_gitpeek("my_org/my_repo")
        self.assertFalse(self._is_parse_error(stderr), stderr)

class TestFetchRepo(unittest.TestCase):

    def test_404_prints_not_found(self):
        rc, _, stderr = run_gitpeek("this-owner-does-not-exist-xyzzy/no-repo-here-abc")
        if rc == 0:
            self.skipTest("Unexpected success — possibly cached/proxied")
        self.assertTrue(
            "404" in stderr or "403" in stderr or "rate" in stderr.lower(),
            f"Expected error message, got: {stderr!r}",
        )

    def test_missing_token_env_var_does_not_crash(self):
        env = {k: v for k, v in os.environ.items() if k != "GITHUB_TOKEN"}
        rc, stdout, stderr = run_gitpeek("owner/repo", env=env)
        self.assertNotIn("panic", stderr.lower())
        self.assertNotIn("panic", stdout.lower())

class TestOutputFormat(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        rc, stdout, stderr = run_gitpeek("XRenso/goshka")
        cls.rc = rc
        cls.stdout = stdout
        cls.stderr = stderr

    def _skip_if_network_error(self):
        if self.rc != 0:
            self.skipTest(f"Network/API unavailable: {self.stderr.strip()}")

    def test_output_contains_owner_and_repo(self):
        self._skip_if_network_error()
        self.assertIn("XRenso/goshka", self.stdout)

    def test_output_contains_separator_line(self):
        self._skip_if_network_error()
        self.assertIn("─", self.stdout)

    def test_output_contains_stars_symbol(self):
        self._skip_if_network_error()
        self.assertIn("★", self.stdout)

    def test_output_contains_forks_symbol(self):
        self._skip_if_network_error()
        self.assertIn("⑂", self.stdout)

    def test_output_contains_created_date(self):
        self._skip_if_network_error()
        self.assertRegex(self.stdout, r"\d{4}-\d{2}-\d{2}")

    def test_output_not_empty_on_success(self):
        self._skip_if_network_error()
        self.assertTrue(len(self.stdout.strip()) > 0)


if __name__ == "__main__":
    unittest.main(verbosity=2)
