load("@com_github_datadog_rules_oci//oci:defs.bzl", "oci_image_layer", "oci_image_index", "oci_image")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

def go_multiarch_image(name, base, archs, binary_name = "", binary_dir = "/app", **kwargs):
    """
    Create a multiarch image from a go library.

    NOTE: This probably should be called something like go_image, but doing this
    to prevent confusion with rules_docker.

    """

    os = "linux"

    visibility = kwargs.get("visibility", None)
    tags = kwargs.get("tags", None)

    if binary_name == "":
        binary_name = name

    manifests = []
    for arch in archs:
        go_binary_name = "{name}.{os}-{arch}-go-binary".format(name = name, os = os, arch = arch)
        go_binary_out = "{binary_name}-{os}-{arch}".format(binary_name = binary_name, os = os, arch = arch)
        go_binary(
            name = go_binary_name,
            goos = os,
            goarch = arch,
            out = go_binary_out,
            **kwargs
        )

        layer_name = "{name}.{os}-{arch}-go-layer".format(name = name, os = os, arch = arch)
        oci_image_layer(
            name = layer_name,
            files = [
                go_binary_name,
            ],
            symlinks = {
                "{}/{}".format(binary_dir, binary_name): "{}/{}".format(binary_dir, go_binary_out),
            },
            directory = binary_dir,
        )

        image_name = "{name}.{os}-{arch}-image".format(name = name, os = os, arch = arch)
        oci_image(
            name = image_name,
            base = base,
            layers = [
                layer_name,
            ],
            os = os,
            arch = arch,
        )
        manifests.append(image_name)

    oci_image_index(
        name = name,
        manifests = manifests,
        tags = tags,
        visibility = visibility,
    )
