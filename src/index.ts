import { OpenAPIRouter } from '@cloudflare/itty-router-openapi';
import { version } from '../package.json';

export const router = OpenAPIRouter({
  schema: {
    info: {
      title: 'Edge Image Classification',
      description:
        'A proof of concept to perform image classification at the edge using Cloudflare Workers',
      version: `v${version}`,
    },
  },
  docs_url: '/',
});

// 404 for everything else
router.all('*', () => new Response('Not Found.', { status: 404 }));

export default {
  fetch: router.handle,
};
