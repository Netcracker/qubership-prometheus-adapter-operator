import json
import re
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PATHS_FILTER = "dorny/paths-filter@ceb8a2b8f2d89434be7ff52d3de7ec3738c5cc9d # v4.0.3"


class RenovateAutomergeContractTest(unittest.TestCase):
    def workflow(self, name: str) -> str:
        return (ROOT / ".github" / "workflows" / name).read_text()

    def test_renovate_uses_only_the_shared_automerge_policy(self) -> None:
        config = json.loads((ROOT / "renovate.json").read_text())

        self.assertIn("github>Netcracker/renovate-config:automerge", config["extends"])
        self.assertNotIn("platformAutomerge", config)
        self.assertNotIn("automergeType", config)
        self.assertNotIn("autoApprove", config)
        for rule in config.get("packageRules", []):
            self.assertNotIn("automerge", rule)

    def test_required_workflows_have_pr_only_fail_closed_gates(self) -> None:
        expected = {
            "super-linter.yaml": "Lint Gate",
            "test-sonar-go-coverage.yaml": "Go CI Gate",
            "build.yaml": "Build Gate",
            "helm-apiservice-smoke-test.yaml": "Helm CI Gate",
        }

        for filename, gate in expected.items():
            with self.subTest(filename=filename):
                workflow = self.workflow(filename)
                self.assertIn(f"name: {gate}", workflow)
                self.assertIn("if: always() && github.event_name == 'pull_request'", workflow)
                self.assertRegex(workflow, r"needs\.[a-z0-9-]+\.result != 'success'")

        all_workflows = "\n".join(self.workflow(name) for name in expected)
        self.assertNotIn("Renovate Config Gate", all_workflows)

    def test_path_filtered_validation_starts_on_every_pull_request(self) -> None:
        for filename in (
            "test-sonar-go-coverage.yaml",
            "build.yaml",
            "helm-apiservice-smoke-test.yaml",
        ):
            with self.subTest(filename=filename):
                workflow = self.workflow(filename)
                triggers = workflow.split("jobs:", 1)[0]
                self.assertNotRegex(triggers, re.compile(r"^\s+paths(?:-ignore)?:", re.MULTILINE))
                self.assertIn(PATHS_FILTER, workflow)
                self.assertIn("pull-requests: read", workflow)
                self.assertIn(f".github/workflows/{filename}", workflow)

    def test_go_workflow_discovers_all_modules(self) -> None:
        workflow = self.workflow("test-sonar-go-coverage.yaml")

        self.assertNotIn("go-module-dir:", workflow)
        self.assertNotIn("go-module-configs:", workflow)
        self.assertNotIn("test-api:", workflow)

    def test_pull_request_container_build_is_read_only(self) -> None:
        workflow = self.workflow("build.yaml")

        self.assertIn("multiplatform_build:", workflow)
        self.assertIn("multiplatform_publish:", workflow)
        self.assertIn("packages: read", workflow)
        self.assertIn("push: false", workflow)
        self.assertIn("packages: write", workflow)

    def test_sonar_classifies_python_contracts_as_tests(self) -> None:
        config = (ROOT / "sonar-project.properties").read_text()

        self.assertIn("sonar.exclusions=**/*_test.go,tests/**/*.py", config)
        self.assertIn("sonar.test.inclusions=**/*_test.go,tests/**/*.py", config)


if __name__ == "__main__":
    unittest.main()
