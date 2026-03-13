const { getDefaultConfig } = require('expo/metro-config');
const path = require('path');

const projectRoot = __dirname;
const monorepoRoot = path.join(projectRoot, '../..');

const config = getDefaultConfig(projectRoot);

config.resolver.nodeModulesPaths = [
  path.resolve(projectRoot, 'node_modules'),
  path.resolve(monorepoRoot, 'node_modules'),
];

config.resolver.extraNodeModules = {
  '@rcm/shared': path.resolve(monorepoRoot, 'package/shared/src'),
};

config.watchFolders = [
  projectRoot,
  path.resolve(monorepoRoot, 'package/shared'),
];

module.exports = config;
