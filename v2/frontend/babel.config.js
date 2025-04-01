module.exports = {
    presets: [
      ['@babel/preset-env', { targets: { node: 'current' } }],
      ['@babel/preset-react', { runtime: 'automatic' }]
    ],
    plugins: [
      function () {
        return {
          visitor: {
            MetaProperty(path) {
              path.replaceWithSourceString('process')
            },
          },
        }
      },
    ]
  };