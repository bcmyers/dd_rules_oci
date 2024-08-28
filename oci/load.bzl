""" oci_image_load rule """

load(
    "@com_github_datadog_rules_oci//oci:providers.bzl",
    "OCIImageTar",
)

def _impl(ctx):
    file = ctx.attr.tar[OCIImageTar].file

    out = ctx.actions.declare_file("{}_/run.sh".format(ctx.label.name))
    ctx.actions.write(
        output = out,
        content = """
#!/usr/bin/env bash
set -euo pipefail
{ocitool} load --file={file}
""".format(
            ocitool = ctx.executable._bin.short_path,
            file = file.short_path,
        ),
        is_executable = True,
    )

    runfiles = ctx.runfiles(files = [
        ctx.executable._bin,
        file,
    ])

    return [
        DefaultInfo(
            executable = out,
            runfiles = runfiles,
        ),
    ]

_oci_image_load = rule(
    implementation = _impl,
    attrs = {
        "tar": attr.label(
            mandatory = True,
            providers = [OCIImageTar],
        ),
        "_bin": attr.label(
            default = Label("//go/cmd/ocitool"),
            executable = True,
            cfg = "exec",
        ),
    },
    executable = True,
)

# TODO: support more than docker
def oci_image_load(
        *,
        name,
        tar,
        **kwargs):
    tags = kwargs.pop("tags", [])
    tags = {x: True for x in tags}
    tags["manual"] = True
    tags = [x for x in tags.keys()]

    _oci_image_load(
        name = name,
        tar = tar,
        tags = tags,
        **kwargs
    )
