{
  // escapeStringRegex(string s) string
  // escapeStringRegex escapes all regular expression metacharacters and returns a regular expression that matches the literal text.
  escapeStringRegex: std.native('escapeStringRegex'),

  // regexMatch(string regex, string s) boolean
  // regexMatch returns whether the given string is matched by the given RE2 regular expression.
  regexMatch: std.native('regexMatch'),

  // regexSubst(string regex, string src, string repl) string
  // regexSubst replaces all matches of the re2 regular expression with the replacement string.
  regexSubst: std.native('regexSubst'),
}
