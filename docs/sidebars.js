/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  // docsSidebar: [{ type: "autogenerated", dirName: "." }],

  docsSidebar: [
    {
      type: 'doc',
      id: 'README',
      label: 'Get Started',
    },
    {
      type: 'category',
      label: 'Install',
      items: ['install/binary', 'install/docker', 'install/npm'],
    },
    {
      type: 'doc',
      id: 'cli',
      label: 'CLI Documentation',
    },
    {
      type: 'doc',
      id: 'sql-dialect',
      label: 'SQL Dialect',
    },
    {
      type: 'doc',
      id: 'defining-metrics',
      label: 'Defining Metrics',
    },
    {
      type: 'category',
      label: 'Contributors',
      items: ['contributors/development', 'contributors/guidelines', 'contributors/local-testing'],
    },
  ],
};

module.exports = sidebars;
