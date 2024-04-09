{
  // semverParse(string v) { major, minor, patch, prerelease, metadata }
  // semverParse parses a semver string returning an object containing the components
  semverParse: std.native('semverParse'),

  // semverMatchesConstraint(string v, constraint s) bool
  // semverMatchesConstraint returns true if semver string v matches the constraint.
  // See https://github.com/Masterminds/semver?tab=readme-ov-file#basic-comparisons
  // for full description of constaint syntax
  semverMatchesConstraint: std.native('semverMatchesConstraint'),
}
