""" oci_image_tar rule """

load(
    "@com_github_datadog_rules_oci//oci:providers.bzl",
    "OCIDescriptor",
    "OCIImageTar",
    "OCILayout",
)

def _impl(ctx):
    toolchain = ctx.toolchains["@com_github_datadog_rules_oci//oci:toolchain"]

    args = ctx.actions.args()
    inputs = []

    args.add("tar")

    if ctx.attr.gzip:
        out = ctx.actions.declare_file("{}_/oci_layout.tar.gz".format(ctx.label.name))
    else:
        out = ctx.actions.declare_file("{}_/oci_layout.tar".format(ctx.label.name))
    args.add("--out", out.path)

    blob_index = ctx.attr.image[OCILayout].blob_index
    args.add("--blob-index", blob_index.path)
    inputs.append(blob_index)

    descriptor_file = ctx.attr.image[OCIDescriptor].descriptor_file
    args.add("--descriptor-file", descriptor_file.path)
    inputs.append(descriptor_file)

    files = ctx.attr.image[OCILayout].files
    for f in files.to_list():
        args.add("--file", f.path)
        inputs.append(f)

    args.add("--gzip", ctx.attr.gzip)

    ctx.actions.run(
        outputs = [out],
        executable = toolchain.sdk.ocitool,
        arguments = [args],
        inputs = inputs,
    )

    return [
        DefaultInfo(
            files = depset([out]),
        ),
        OCIImageTar(
            file = out,
            gzip = ctx.attr.gzip,
        ),
    ]

_oci_image_tar = rule(
    implementation = _impl,
    attrs = {
        "gzip": attr.bool(
            mandatory = True,
        ),
        "image": attr.label(
            mandatory = True,
            providers = [OCIDescriptor, OCILayout],
        ),
    },
    toolchains = ["@com_github_datadog_rules_oci//oci:toolchain"],
)

def oci_image_tar(
        *,
        name,
        gzip,
        image,
        **kwargs):
    tags = kwargs.pop("tags", [])
    tags = {x: True for x in tags}
    tags["manual"] = True
    tags = [x for x in tags.keys()]

    _oci_image_tar(
        name = name,
        gzip = gzip,
        image = image,
        tags = tags,
        **kwargs
    )
