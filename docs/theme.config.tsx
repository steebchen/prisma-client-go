import { DocsThemeConfig } from 'nextra-theme-docs'

export default {
  logo: <span>Go Prisma</span>,
  project: {
    link: 'https://github.com/prisma/prisma-client-go',
  },
  head: (
    <>
      <link rel="icon" href="/icon.png" sizes="any"/>
    </>
  ),
  docsRepositoryBase: 'https://github.com/prisma/prisma-client-go/tree/main/docs',
  footer: {
    text: (
      <span>
        All source code and content licensed under&nbsp;
        <a href="https://github.com/prisma/prisma-client-go/blob/main/LICENSE" target="_blank">
          Apache 2.0
        </a>
      </span>
    ),
  },
  useNextSeoProps() {
    return {
      titleTemplate: '%s – Prisma Client Go',
      description: 'Prisma Client Go is an auto-generated and fully type-safe database client',
      openGraph: {
        type: 'website',
        url: 'https://goprisma.org',
        description: 'Prisma Client Go is an auto-generated and fully type-safe database client',
        site_name: 'Go Prisma',
        title: 'Prisma Client Go',
      },
    }
  },
  sidebar: {
    defaultMenuCollapseLevel: 1,
  },
} as DocsThemeConfig
