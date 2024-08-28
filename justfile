build:
    #!/usr/bin/env bash
    set -euo pipefail
    bzl build //tests/load:image2.tar
    rm -rf foo
    mkdir -p foo
    cp bazel-bin/tests/load/image2.tar_/oci_layout.tar foo/
    (
        cd foo
        tar -xzf oci_layout.tar
    )
    find "foo" -type f | while read -r file; do
        if jq empty "$file" > /dev/null 2>&1; then
            chmod 660 "$file"
            jq . "$file" > "$file.tmp" && mv "$file.tmp" "$file"
        fi
    done
