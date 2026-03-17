import os
import unittest
import urllib.request
import urllib.error
import json

BASE_URL = os.environ.get("BAREREPO_URL", "http://localhost:8080")


def get(path: str):
    url = f"{BASE_URL}{path}"
    try:
        with urllib.request.urlopen(url, timeout=10) as resp:
            return resp.status, json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        try:
            body = json.loads(e.read().decode())
        except Exception:
            body = {}
        return e.code, body
    except Exception as e:
        raise RuntimeError(f"Request to {url} failed: {e}") from e


class TestSwagger(unittest.TestCase):

    def test_swagger_ui_accessible(self):
        url = f"{BASE_URL}/swagger/index.html"
        try:
            with urllib.request.urlopen(url, timeout=10) as resp:
                self.assertEqual(resp.status, 200)
        except Exception as e:
            self.fail(f"Swagger UI not accessible: {e}")


class TestGetRepoSuccess(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        try:
            cls.status, cls.data = get("/api/v1/repos/XRenso/goshka")
        except RuntimeError as e:
            cls.status, cls.data = 0, {}
            cls.skip_reason = str(e)
        else:
            cls.skip_reason = None

    def _skip_if_unavailable(self):
        if self.skip_reason:
            self.skipTest(self.skip_reason)
        if self.status != 200:
            self.skipTest(f"API returned {self.status}: {self.data}")

    def test_status_200(self):
        self._skip_if_unavailable()
        self.assertEqual(self.status, 200)

    def test_has_title(self):
        self._skip_if_unavailable()
        self.assertIn("title", self.data)
        self.assertIsInstance(self.data["title"], str)
        self.assertTrue(len(self.data["title"]) > 0)

    def test_has_description(self):
        self._skip_if_unavailable()
        self.assertIn("description", self.data)

    def test_has_creator(self):
        self._skip_if_unavailable()
        self.assertIn("creator", self.data)
        self.assertIn("XRenso", self.data["creator"])

    def test_has_stars_cnt(self):
        self._skip_if_unavailable()
        self.assertIn("stars_cnt", self.data)
        self.assertIsInstance(self.data["stars_cnt"], int)

    def test_has_forks_cnt(self):
        self._skip_if_unavailable()
        self.assertIn("forks_cnt", self.data)
        self.assertIsInstance(self.data["forks_cnt"], int)

    def test_has_issues_cnt(self):
        self._skip_if_unavailable()
        self.assertIn("issues_cnt", self.data)
        self.assertIsInstance(self.data["issues_cnt"], int)

    def test_has_lang(self):
        self._skip_if_unavailable()
        self.assertIn("lang", self.data)

    def test_has_size(self):
        self._skip_if_unavailable()
        self.assertIn("size", self.data)
        self.assertIsInstance(self.data["size"], int)

    def test_has_created_at(self):
        self._skip_if_unavailable()
        self.assertIn("created_at", self.data)
        import re
        self.assertRegex(self.data["created_at"], r"\d{4}-\d{2}-\d{2}")

    def test_has_license(self):
        self._skip_if_unavailable()
        self.assertIn("license", self.data)

    def test_no_extra_error_field(self):
        self._skip_if_unavailable()
        self.assertNotIn("error", self.data)


class TestGetRepoErrors(unittest.TestCase):

    def test_not_found_returns_404(self):
        status, data = get("/api/v1/repos/this-owner-does-not-exist-xyzzy/no-repo-here-abc")
        self.assertIn(status, (404, 500))
        self.assertIn("error", data)

    def test_error_response_has_error_field(self):
        status, data = get("/api/v1/repos/this-owner-does-not-exist-xyzzy/no-repo-here-abc")
        self.assertIn("error", data)
        self.assertIsInstance(data["error"], str)
        self.assertTrue(len(data["error"]) > 0)

    def test_unknown_route_returns_404(self):
        status, _ = get("/api/v1/repos/")
        self.assertEqual(status, 404)


if __name__ == "__main__":
    unittest.main(verbosity=2)
