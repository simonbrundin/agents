// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxt/eslint',
    '@nuxt/ui',
    '@vueuse/nuxt',
    '@nuxt/icon',
    'nuxt-auth-utils'
  ],

  runtimeConfig: {
    hasuraUrl: process.env.HASURA_URL || 'http://localhost:8080',
    hasuraAdminSecret: process.env.HASURA_ADMIN_SECRET || 'hasura-dev-secret',
    oauth: {
      github: {}
    },
    session: {
      maxAge: 60 * 60 * 24 * 7,
      name: 'nuxt-session',
      cookie: {
        sameSite: 'lax',
        secure: false,
        httpOnly: true,
        path: '/'
      }
    }
  },

  auth: {
    origin: process.env.NUXT_AUTH_ORIGIN || 'http://localhost:3000'
  },

  devtools: {
    enabled: true
  },

  css: ['~/assets/css/main.css'],

  routeRules: {
    '/api/**': {
      cors: true
    },
    '/v1/graphql': {
      cors: true
    }
  },

  compatibilityDate: '2024-07-11',

  eslint: {
    config: {
      stylistic: {
        commaDangle: 'never',
        braceStyle: '1tbs'
      }
    }
  }
})
