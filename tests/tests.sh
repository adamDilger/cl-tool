#!/usr/bin/env bats

function setup {
	mkdir .changelog
}

function teardown {
	rm -rf .changelog
}

get_index() {
	echo "${lines[$((${#lines[@]} + $1))]}"
}

@test "basic changelog entry" {
	root=".changelog/2.0.0_2022-07-31"
	mkdir "$root"

	echo -e "added:\n - 'there was a change'" >"$root/2022-07-31-new-feature.yml"

	run "../cl-tool"

	[ "$(get_index -3)" == "## [2.0.0] - 2022-07-31" ]
	[ "$(get_index -2)" == "### Added" ]
	[ "$(get_index -1)" == "- there was a change" ]
}

@test "basic changelog entry with multiple files" {
	root=".changelog/2.0.0_2022-07-31"
	mkdir "$root"

	echo -e "added:\n - 'first added'" >"$root/2022-07-31-new-feature.yml"
	echo -e "added:\n - 'second added'\n - 'third added in same file'" >"$root/2022-08-31-second-feature.yml"
	echo -e "changed:\n - 'there was a change'" >"$root/2022-09-31-second-feature.yml"
	echo -e "added:\n - 'forth added'" >"$root/2022-10-31-second-feature.yml"

	run "../cl-tool"

	[ "$(get_index -8)" == "## [2.0.0] - 2022-07-31" ]
	[ "$(get_index -7)" == "### Added" ]
	[ "$(get_index -6)" == "- forth added" ]
	[ "$(get_index -5)" == "- second added" ]
	[ "$(get_index -4)" == "- third added in same file" ]
	[ "$(get_index -3)" == "- first added" ]
	[ "$(get_index -2)" == "### Changed" ]
	[ "$(get_index -1)" == "- there was a change" ]
}

@test "basic changelog entry with multiple versions" {
	root1=".changelog/2.0.0_2022-07-31"
	root2=".changelog/2.0.1_2022-07-31"
	mkdir "$root1"
	mkdir "$root2"

	echo -e "added:\n - 'first version'" >"$root1/2022-07-31-new-feature.yml"
	echo -e "added:\n - 'second version'" >"$root2/2022-08-31-second-feature.yml"

	run "../cl-tool"

	[ "$(get_index -3)" == "## [2.0.0] - 2022-07-31" ] # lower version number is last
	[ "$(get_index -6)" == "## [2.0.1] - 2022-07-31" ]
}

@test "Unreleased folder creates an unreleased changelog group" {
	root=".changelog/Unreleased"
	mkdir "$root"

	echo -e "changed:\n  - 'an unreleased entry'" >.changelog/Unreleased/2022-01-01-helloworld.yml

	run "../cl-tool"

	[ "$(get_index -3)" == "## Unreleased" ]
	[ "$(get_index -2)" == "### Changed" ]
	[ "$(get_index -1)" == "- an unreleased entry" ]
}

@test "head.md works" {
	echo -e "testing\nhead.md" >.changelog/head.md

	run "../cl-tool"

	[ "${lines[0]}" == "testing" ]
	[ "${lines[1]}" == "head.md" ]
}

@test "tail.md works" {
	echo -e "testing\ntail.md" >.changelog/tail.md

	run "../cl-tool"

	[ "$(get_index -2)" == "testing" ]
	[ "$(get_index -1)" == "tail.md" ]
}

@test "version command works" {
	run "../cl-tool" "-v"
	echo $output
	[[ "$output" == "Version: "* ]]
}
