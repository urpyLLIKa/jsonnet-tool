local x = import 'lib.libsonnet';

{
  testcase1: function() {
    actual: std.trace("calling x()", x()),
    expectJSON: './fixtures/test3.fail.testcase1.json',
  },
  testcase2: function() {
    actual: std.trace("manifest", std.manifestYamlDoc({ hello: 'world' })),
    expectYAML: './fixtures/test3.fail.testcase2.yaml',
  },
  testcase3: function() {
    actual: |||
      Hello World

      This is good old plain text.
    |||,
    expectYAML: './fixtures/test3.fail.testcase3.txt',
  },
}
