local x = import 'lib.libsonnet';

{
  testcase1: function() {
    actual: std.trace("this is a trace", x()),
    expectJSON: './fixtures/test1.testcase1.json',
  },
  testcase2: function() {
    actual: std.manifestYamlDoc({ hello: 'world' }),
    expectYAML: './fixtures/test1.testcase2.yaml',
  },
  testcase3: function() {
    actual: |||
      Hello World

      This is good old plain text.
    |||,
    expectPlainText: './fixtures/test1.testcase3.txt',
  },
} + {
  ['testgen' + i]: function() {
    actual: x(),
    expectJSON: './fixtures/test1.testcase1.json',
  }
  for i in std.range(1, 6)
}
