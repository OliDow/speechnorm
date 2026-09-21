export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    // Dependabot's grouped-update descriptions contain unwrapped tables and URLs.
    'body-max-line-length': [0],
  },
};
