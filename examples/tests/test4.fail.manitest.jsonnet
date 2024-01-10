local x = import 'lib.libsonnet';

{
  testcase1: function() {
    actual: { hello: std.trace("hello world", 'world') },
    expectJSON: './fixtures/test4.fail.MISSING.json',
  },
  testcase2: function() {
    actual: std.manifestYamlDoc({ hello: std.trace("hello world", 'world') }),
    expectYAML: './fixtures/test4.fail.MISSING.yaml',
  },
  testcase3: function() {
    actual: std.manifestYamlDoc({ hello: std.trace("hello world", 'world') }),
    expectYAML: './fixtures/test4.fail.MISSING.yaml',
  },
  testcase4: function() {
    actual: {
      deep: {
        json: 'hello world'
      }
    },
    expect: {
      deep: {
        json: 'hello earth',
        other: 'field'
      }
    }
  },
}
