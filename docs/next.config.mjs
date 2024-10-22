import nextra from "nextra";

const withNextra = nextra({
  latex: true,
  search: {
    codeblocks: true,
  },
  useContentDir: true,
});

export default withNextra({
  reactStrictMode: true,
});
