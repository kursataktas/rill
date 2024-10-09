/**
 * Creating a sidebar enables you to:
 * - create an ordered group of docs
 * - render a sidebar for each doc of that group
 * - provide next/previous navigation
 *
 * The sidebars can be generated from the filesystem, or explicitly defined here.
 *
 * Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  docsSidebar: [
    { type: "autogenerated", dirName: ".", },
  ],

  tutorialsSidebar: [
    {
      type: 'category',
      label: 'Tutorials',
      collapsible: false,
      collapsed: false,
      link: {
        type: 'doc',
        id: 'tutorials/index',

      },
      items: [
        {
          type: 'category',
          label: 'Rill Developer to Rill Cloud in 5 steps!',
          description: 'Rill Developer to  to Rill Cloud',

          items: [
            'tutorials/rill_basics/launch',
            'tutorials/rill_basics/import',
            'tutorials/rill_basics/model',
            'tutorials/rill_basics/dashboard',
            'tutorials/rill_basics/deploy',
            'tutorials/rill_basics/success',
          ]
        },



        {
          type: 'category',
          label: "Rill's Advanced Features",
          description: 'Advanced Features and beyond',
          link: {
            type: 'doc',
            id: 'tutorials/rill_advanced_features/overview',

          },

          items: [

            {
              type: 'category',
              label: 'Back to Rill Developer',
              description: 'Make some changes to our SQL model and Dashboard',
              items: [
                'tutorials/rill_advanced_features/advanced_developer/advanced-modeling',
                'tutorials/rill_advanced_features/advanced_developer/advanced-dashboard',
                'tutorials/rill_advanced_features/advanced_developer/update-rill-cloud',

              ]
            },
            {
              type: 'category',
              label: 'Rill Canvas Dashboards',
              items: [
                'tutorials/rill_advanced_features/canvas_dashboards/getting-started',
                'tutorials/rill_advanced_features/canvas_dashboards/template-charts',
                'tutorials/rill_advanced_features/canvas_dashboards/vega-lite',
                'tutorials/rill_advanced_features/canvas_dashboards/vega-lite2',
                'tutorials/rill_advanced_features/canvas_dashboards/canvas-dashboards',
                'tutorials/rill_advanced_features/canvas_dashboards/filters'

              ]
            },

            {
              type: 'category',
              label: 'Incremental Models',
              
              items: [
                'tutorials/rill_advanced_features/incremental_models/introduction',
                {
                  type: 'category',
                  label: 'Basic Incremental and Split Model Examples',
                  items: [
                    'tutorials/rill_advanced_features/incremental_models/simple-examples/incremental_now',
                    'tutorials/rill_advanced_features/incremental_models/simple-examples/split_now',
                  ]
                },
                'tutorials/rill_advanced_features/incremental_models/cloud-storage-splits',
                'tutorials/rill_advanced_features/incremental_models/data-warehouse-splits',
                'tutorials/rill_advanced_features/incremental_models/staging-connectors'

              ]
            },
            {
              type: 'category',
              label: 'Custom APIs',
              items: [
                'tutorials/rill_advanced_features/custom_api/getting-started',
                'tutorials/rill_advanced_features/custom_api/create-api',
                'tutorials/rill_advanced_features/custom_api/test-api',
              ]
            },
          ]
        },
        {
          type: 'category',
          label: 'Administration Topics',
          description: 'Rill Administration Topics',
          link: {
            type: 'doc',
            id: 'tutorials/administration/index',

          },
          items: [
            {
              type: 'category',
              label: 'User Management',
              items: [
                'tutorials/administration/user/user-management',
                'tutorials/administration/user/user-group-management',
              ]
            },
            {
              type: 'category',
              label: 'Project Management',
              items: [
                'tutorials/administration/project/project-maintanence',
                'tutorials/administration/project/alerts',
                'tutorials/administration/project/credential-envvariable-mangement',
                'tutorials/administration/project/github',
              ]
            },

          ]
        },

        {
          type: 'category',
          label: 'Rill and ClickHouse to Dashboarding in 4 steps!',
          description: 'For our friends from ClickHouse, a revamped guide.',
          items: [
            'tutorials/rill_clickhouse/r_ch_launch',
            'tutorials/rill_clickhouse/r_ch_connect',
            'tutorials/rill_clickhouse/r_ch_dashboard',
            'tutorials/rill_clickhouse/r_ch_deploy',
            {
              type: 'category',
              label: 'Extra Topics:',
              items: [
                'tutorials/rill_clickhouse/r_ch_ingest',
              ],

            },


          ]
        },
        {
          type: 'category',
          label: "Other Concepts and How-to's",
          description: 'For guides that are not quite Rill related but needs consideration',
          items: [
            'tutorials/other/add-column-dimension',
            'tutorials/other/dashboard-row-policies',
            'tutorials/other/custom-charts',
            'tutorials/other/create-map-component',
            'tutorials/other/embed-dashboard',
      //      'tutorials/other/deep-dive-incremental-modeling',
            'tutorials/other/Rill Cloud/share-dashboard-publicly',
            'tutorials/other/Rill Cloud/views',
      //      'tutorials/other/Rill Cloud/visual-metric-editor-rc',
       //     'tutorials/other/yaml-vs-ui',
       //     'tutorials/other/trial-check',
            'tutorials/other/avg_avg',

          ]
        },


      ],
    },
  ],

  refSidebar: [
    { type: "autogenerated", dirName: "reference" },
  ],
};



module.exports = sidebars;
