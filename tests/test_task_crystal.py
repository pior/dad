import pytest


def test_env(cmd, project):
    project.write_devyml("""
        up:
        - crystal: '0.35.1'
    """)

    cmd.run("bud up")
    
    output = cmd.run("crystal version")
    assert "Crystal 0.35.1" in output
