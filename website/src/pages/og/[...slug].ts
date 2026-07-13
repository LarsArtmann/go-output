import { getCollection } from "astro:content";
import { OGImageRoute } from "astro-og-canvas";

const docs = await getCollection("docs");
const docPages = Object.fromEntries(docs.map(({ data, id }) => [id, { data }]));

const pages = {
  ...docPages,
  home: {
    data: {
      title: "go-output",
      description:
        "Write your data once. Render it anywhere. 16 formats, 3 data shapes, zero lock-in. Plus NOM-style real-time progress for Go.",
    },
  },
};

export const { getStaticPaths, GET } = await OGImageRoute({
  pages,
  param: "slug",
  getImageOptions: (_path, page) => ({
    title: page.data.title,
    description: page.data.description,
    bgGradient: [[12, 10, 9]],
    border: { color: [34, 211, 238], width: 4 },
    padding: 80,
  }),
});
