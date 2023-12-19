local x = import 'lib.libsonnet';

{
  testcase1: {
    actual: x(),
    expectJSON: './fixtures/test1.testcase1.json',
  },
  testcase2: {
    actual: std.manifestYamlDoc({ hello: 'world' }),
    expectYAML: './fixtures/test1.testcase2.yaml',
  },
}
