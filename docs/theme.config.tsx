import Image from 'next/image'
import { DocsThemeConfig } from 'nextra-theme-docs'

export default {
  logo: <span>Go Prisma</span>,
  project: {
    link: 'https://github.com/steebchen/prisma-client-go',
  },
  head: (
    <>
      <link rel="icon" href="/icon.png" sizes="any"/>
    </>
  ),
  docsRepositoryBase: 'https://github.com/steebchen/prisma-client-go/tree/main/docs',
  footer: {
    text: (
      <>
        <div style={{ width: '100%' }}>
          All source code and content licensed under&nbsp;
          <a href="https://github.com/steebchen/prisma-client-go/blob/main/LICENSE" target="_blank">
            Apache 2.0
          </a>
        </div>

        <div>
          <a href="https://vercel.com/?utm_source=prisma-client-go&utm_campaign=oss" target="_blank">
            <Image
              alt="Powered by Vercel"
              src="https://images.ctfassets.net/e5382hct74si/78Olo8EZRdUlcDUFQvnzG7/fa4cdb6dc04c40fceac194134788a0e2/1618983297-powered-by-vercel.svg"
              width="212"
              height="44"
            />
          </a>
        </div>
      </>
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
