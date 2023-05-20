import { DocsThemeConfig } from 'nextra-theme-docs'

export default {
  logo: <span>Prisma Client Go</span>,
  project: {
    link: 'https://github.com/prisma/prisma-client-go',
  },
  docsRepositoryBase: 'https://github.com/prisma/prisma-client-go/tree/main/docs',
  footer: {
    text: (
      <span>
        All source code and content licensed under
        <a href="https://github.com/prisma/prisma-client-go/blob/main/LICENSE" target="_blank">
          Apache 2.0
        </a>
        <a href="https://goprisma.org" target="_blank">
          Prisma Client Go
        </a>
      </span>
    ),
  },
  useNextSeoProps() {
    return {
      titleTemplate: '%s – Prisma Client Go',
    }
  },
  sidebar: {
    defaultMenuCollapseLevel: 1,
  },
} as DocsThemeConfig
