#!/usr/bin/env python3
# This file is used to update the version number in all relevant places
# The SemVer (https://semver.org) versioning system is used.
import re

makefile_path = "./Makefile"
changelog_path = "./CHANGELOG.md"

# Extract old version from the Makefile
with open(makefile_path, "r") as makefile:
    content = makefile.read()
    old_version = content.split('version := ')[1].split('\n')[0]
    print(f"Found old version in {makefile_path}: {old_version}")

# Attempt to read new version from user input
try:
    VERSION = input(
        f"Current version: {old_version}\nNew version (without 'v' prefix): ")
except KeyboardInterrupt:
    print("\nCanceled by user")
    quit()

if VERSION == "":
    VERSION = old_version

if not re.match(r"^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$", VERSION):
    print(
        f"\x1b[31mThe version: '{VERSION}' is not a valid SemVer version.\x1b[0m")
    quit()

# Update version in the Makefile
with open(makefile_path, 'w') as makefile:
    makefile.write(content.replace(old_version, VERSION))

# Update version in Changelog
with open(changelog_path, 'r') as changelog:
    content = changelog.read()
    old_version = content.split("## Changelog for v")[1].split('\n')[0]
    print(f"Found old version in {changelog_path}: {old_version}")

with open(changelog_path, "w") as changelog:
    changelog.write(content.replace(old_version, VERSION))

print(f"Version has been changed from '{old_version}' -> '{VERSION}'")
