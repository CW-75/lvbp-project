import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'

// https://vitepress.dev/reference/site-config
export default withMermaid(
  defineConfig({
    title: 'LVBP GameCast Documentation',
    description: 'Live Baseball GameCast & Standings Engine - Technical Documentation',
    base: '/lvbp-project/',
    lastUpdated: true,
    cleanUrls: true,
    
    head: [
      ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo-lvbp-light.svg' }],
      ['meta', { name: 'theme-color', content: '#0B1A32' }]
    ],

    markdown: {
      theme: {
        light: 'github-light',
        dark: 'github-dark'
      },
      lineNumbers: true,
      toc: { level: [2, 3] }
    },

    ignoreDeadLinks: true,

    themeConfig: {
      logo: '/logo-lvbp-light.svg',
      
      nav: [
        { text: 'Inicio', link: '/' },
        { text: 'Arquitectura', link: '/architecture/system-overview' },
        { text: 'Fases', link: '/phases/01-auth-system' },
        { text: 'Referencia DTOs', link: '/dto-reference' },
        { text: 'API Reference', link: '/api-reference' },
        { text: 'Runbooks', link: '/runbooks/local-dev' }
      ],

      sidebar: {
        '/architecture/': [
          {
            text: 'Arquitectura del Sistema',
            items: [
              { text: 'Visión General', link: '/architecture/system-overview' },
              { text: 'Flujo de Datos', link: '/architecture/data-flow' },
              { text: 'ADR - Decision Records', link: '/architecture/adr/index' }
            ]
          }
        ],
        '/phases/': [
          {
            text: 'Fases de Implementación',
            items: [
              { text: '01 - Sistema de Autenticación', link: '/phases/01-auth-system' },
              { text: '02 - Modelos de Dominio', link: '/phases/02-domain-models' },
              { text: '03 - Handler de Ingesta', link: '/phases/03-ingestion-handler' },
              { text: '04 - Motor Scorekeeper', link: '/phases/04-scorekeeper-engine' },
              { text: '05 - Streaming SSE', link: '/phases/05-sse-streaming' },
              { text: '06 - Handler Standings', link: '/phases/06-standings-handler' },
              { text: '07 - Handler Boxscore', link: '/phases/07-boxscore-handler' },
              { text: '08 - Capa Repositorio', link: '/phases/08-repository-layer' }
            ]
          }
        ],
        '/runbooks/': [
          {
            text: 'Runbooks Operacionales',
            items: [
              { text: 'Desarrollo Local', link: '/runbooks/local-dev' },
              { text: 'Despliegue', link: '/runbooks/deployment' }
            ]
          }
        ]
      },

      socialLinks: [
        { icon: 'github', link: 'https://github.com/your-org/lvbp-project' }
      ],

      footer: {
        message: 'Licencia MIT',
        copyright: 'Copyright © 2026 LVBP Project'
      },

      search: {
        provider: 'local'
      },

      editLink: {
        pattern: 'https://github.com/your-org/lvbp-project/edit/main/docs-site/:path',
        text: 'Editar en GitHub'
      }
    },

    vite: {
      plugins: [],
      ssr: {
        noExternal: ['redoc']
      }
    }
  })
)